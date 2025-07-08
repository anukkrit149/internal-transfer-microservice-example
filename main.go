package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	"internal-transfer-microservice/internal/config"
	"internal-transfer-microservice/internal/factory"
	"internal-transfer-microservice/internal/routes"
	"internal-transfer-microservice/pkg/logger"
)

var (
	configPath string
)

func main() {
	// Root command
	rootCmd := &cobra.Command{
		Use:   "account-transfer-microservice",
		Short: "A REST API for account transfers",
		Long:  `A REST API for account transfers built with Go, Gin, GORM, PostgreSQL, and Redis.`,
	}

	// API command
	apiCmd := &cobra.Command{
		Use:   "api",
		Short: "Start the API server",
		Long:  `Start the REST API server for account transfers.`,
		Run:   runAPI,
	}

	// Migrate command
	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		Long:  `Run database migrations to set up or update the database schema.`,
		Run:   runMigrate,
	}

	// Add flags to commands
	apiCmd.Flags().StringVar(&configPath, "config", "", "Path to configuration file")
	migrateCmd.Flags().StringVar(&configPath, "config", "", "Path to configuration file")

	// Add commands to root command
	rootCmd.AddCommand(apiCmd)
	rootCmd.AddCommand(migrateCmd)

	// Execute
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runMigrate(cmd *cobra.Command, args []string) {
	// Load configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	// Create factory
	appFactory, err := factory.NewFactory(cfg)
	if err != nil {
		logger.Fatalf("Failed to create factory: %v", err)
	}
	defer appFactory.Close()

	// Run database migrations
	logger.Info("Running database migrations...")
	if err := appFactory.MigrateDB(); err != nil {
		logger.Fatalf("Failed to migrate database: %v", err)
	}
	logger.Info("Database migrations completed successfully")
}

func runAPI(cmd *cobra.Command, args []string) {
	// Initialize logger
	logConfig := logger.DefaultConfig()
	logConfig.ReportCaller = false
	err := logger.Initialize(logConfig)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	// Set Gin mode
	gin.SetMode(cfg.GetGinMode())

	// Create factory
	appFactory, err := factory.NewFactory(cfg)
	if err != nil {
		logger.Fatalf("Failed to create factory: %v", err)
	}
	defer appFactory.Close()

	// Create router
	router := gin.Default()

	// Setup middleware
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	// Create controllers
	accountController := appFactory.CreateAccountController()

	// Setup routes
	routes.SetupAccountRoutes(router, accountController)

	// Health check route
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})
	})

	// Create server
	server := &http.Server{
		Addr:    ":" + cfg.GetServerPort(),
		Handler: router,
	}

	// Start the server in a goroutine
	go func() {
		logger.Infof("Server starting on port %s", cfg.GetServerPort())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for the interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Received shutdown signal. Starting graceful shutdown...")

	// Create a deadline for server shutdown based on configuration
	shutdownTimeout := time.Duration(cfg.GetShutdownTimeout()) * time.Second
	logger.Infof("Server will shutdown after %s or when all connections are closed", shutdownTimeout)

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server gracefully stopped")
}
