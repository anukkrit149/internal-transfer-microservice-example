package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"internal-transfer-microservice/internal/config"
)

// PostgresDB implements Database interface
type PostgresDB struct {
	db *gorm.DB
}

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB(cfg *config.Config) (Database, error) {
	dsn := cfg.GetDBConnectionString()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	log.Println("Connected to PostgreSQL database")
	return &PostgresDB{db: db}, nil
}

// GetConnection returns the database connection
func (p *PostgresDB) GetConnection() *gorm.DB {
	return p.db
}

// Close closes the database connection
func (p *PostgresDB) Close() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
