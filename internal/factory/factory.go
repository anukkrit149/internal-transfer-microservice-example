package factory

import (
	"internal-transfer-microservice/internal/domain/account"

	"internal-transfer-microservice/internal/config"
	"internal-transfer-microservice/internal/controller"
	"internal-transfer-microservice/internal/infrastructure/cache"
	"internal-transfer-microservice/internal/infrastructure/db"
	"internal-transfer-microservice/internal/repository"
	"internal-transfer-microservice/internal/service"
	"internal-transfer-microservice/pkg/logger"
)

// Factory is responsible for creating and wiring up all components
type Factory struct {
	database db.Database
	cache    cache.Cache
	config   *config.Config
}

// NewFactory creates a new factory
func NewFactory(cfg *config.Config) (*Factory, error) {
	// Initialize database
	database, err := db.NewPostgresDB(cfg)
	if err != nil {
		logger.Errorf("Failed to connect to database: %v", err)
		return nil, err
	}

	// Initialize cache
	redisCache, err := cache.NewRedisCache(cfg)
	if err != nil {
		logger.Errorf("Failed to connect to Redis: %v", err)
		// Close database connection if Redis fails
		database.Close()
		return nil, err
	}

	return &Factory{
		database: database,
		cache:    redisCache,
		config:   cfg,
	}, nil
}

// Close closes all connections
func (f *Factory) Close() {
	if f.database != nil {
		f.database.Close()
	}
	if f.cache != nil {
		f.cache.Close()
	}
}

func (f *Factory) CreateAccountController() *controller.AccountController {
	// Create repository
	accountRepo := repository.NewAccountRepo(f.database)

	// Create service
	accountService := service.NewAccountService(accountRepo, f.cache)

	// Create controller
	accountController := controller.NewAccountController(accountService)

	return accountController
}

// MigrateDB performs database migrations
func (f *Factory) MigrateDB() error {
	// Auto migrate models
	err := f.database.GetConnection().AutoMigrate(
		&account.Model{},
	)
	return err
}
