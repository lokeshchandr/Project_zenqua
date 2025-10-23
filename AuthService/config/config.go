package config

import (
	"errors"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds application configuration loaded from environment variables.
type Config struct {
	JWTSecret          string
	JWTExpiry          time.Duration
	DBPath             string
	SuperAdminName     string
	SuperAdminEmail    string
	SuperAdminPassword string
}

var cfg Config

// Load reads environment variables and populates a global config.
// It loads a .env file if present and ensures mandatory variables exist.
func Load() error {
	// Load .env if found; ignore error if file doesn't exist
	_ = godotenv.Load()

	expiryStr := getenvDefault("JWT_EXPIRY", "24h")
	dur, err := time.ParseDuration(expiryStr)
	if err != nil {
		return errors.New("invalid JWT_EXPIRY; use Go duration format like 24h, 30m")
	}

	cfg = Config{
		JWTSecret:          os.Getenv("JWT_SECRET"),
		JWTExpiry:          dur,
		DBPath:             getenvDefault("DB_PATH", "auth.db"),
		SuperAdminName:     getenvDefault("SUPERADMIN_NAME", "Super Admin"),
		SuperAdminEmail:    os.Getenv("SUPERADMIN_EMAIL"),
		SuperAdminPassword: os.Getenv("SUPERADMIN_PASSWORD"),
	}

	if cfg.JWTSecret == "" {
		return errors.New("JWT_SECRET is required")
	}
	if cfg.SuperAdminEmail == "" || cfg.SuperAdminPassword == "" {
		// Not strictly required to run, but needed to seed a super admin.
		// Return an error to ensure secure initial setup.
		return errors.New("SUPERADMIN_EMAIL and SUPERADMIN_PASSWORD are required")
	}
	return nil
}

// Get returns the loaded global configuration.
func Get() Config { return cfg }

func getenvDefault(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}
