package db

import (
	"gorm.io/gorm"
)

// Database interface defines the operations for database access
type Database interface {
	// GetConnection returns the database connection
	GetConnection() *gorm.DB

	// Close closes the database connection
	Close() error
}
