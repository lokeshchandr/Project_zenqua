package database

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"authservice/config"
	"authservice/models"
	"authservice/utils"
)

// DB is the global database handle.
var DB *gorm.DB

// InitDatabase opens the SQLite database, runs migrations, and seeds the Super Admin.
func InitDatabase() error {
	c := config.Get()

	// Open SQLite database using GORM
	db, err := gorm.Open(sqlite.Open(c.DBPath), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("open sqlite: %w", err)
	}
	DB = db

	// Migrations
	if err := db.AutoMigrate(&models.User{}, &models.Notification{}); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	// Seed Super Admin if not exists
	if err := seedSuperAdmin(db); err != nil {
		return fmt.Errorf("seed superadmin: %w", err)
	}

	log.Println("Database Connected Sucessfully")
	return nil
}

func seedSuperAdmin(db *gorm.DB) error {
	c := config.Get()

	if c.SuperAdminEmail == "" || c.SuperAdminPassword == "" {
		return errors.New("super admin credentials missing")
	}

	var existing models.User
	tx := db.Where("email = ?", c.SuperAdminEmail).First(&existing)
	if tx.Error == nil {
		// Already exists
		return nil
	}
	if tx.Error != nil && !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return tx.Error
	}

	hash, err := utils.HashPassword(c.SuperAdminPassword)
	if err != nil {
		return err
	}

	super := models.User{
		Name:         valueOrDefault(c.SuperAdminName, "Super Admin"),
		Email:        c.SuperAdminEmail,
		PasswordHash: hash,
		Role:         models.RoleSuperAdmin,
		IsApproved:   true,
		CreatedAt:    time.Now(),
	}
	if err := db.Create(&super).Error; err != nil {
		return err
	}
	log.Printf("Seeded Super Admin: %s", c.SuperAdminEmail)
	return nil
}

func valueOrDefault(v, d string) string {
	if v == "" {
		return d
	}
	return v
}
