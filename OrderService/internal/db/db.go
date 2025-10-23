package db

import (
	"fmt"
	"log"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"orderservice/internal/models"
)

// InitDB initializes the SQLite database connection and runs migrations.
func InitDB() *gorm.DB {
	dbPath := "order.db"

	// Open SQLite database
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Run auto-migration
	if err := db.AutoMigrate(&models.Order{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	log.Println("Database initialized successfully")
	return db
}

// GetDB returns database connection (helper for testing).
func GetDB(db *gorm.DB) *gorm.DB {
	if db == nil {
		panic(fmt.Errorf("database connection is nil"))
	}
	return db
}
