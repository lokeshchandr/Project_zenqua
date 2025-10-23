package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"apigateway/internal/routes"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Validate JWT_SECRET is set
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	// Set Gin mode
	ginMode := os.Getenv("GIN_MODE")
	if ginMode != "" {
		gin.SetMode(ginMode)
	}

	// Create Gin router
	router := gin.Default()

	// Setup all routes and middlewares
	routes.SetupRoutes(router)

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// Start the gateway server
	log.Printf("API Gateway starting on port %s", port)
	log.Println("Proxying to backend services:")
	log.Printf("   - Auth Service: %s", getEnvOrDefault("AUTH_SERVICE_URL", "http://localhost:8001"))
	log.Printf("   - Product Service: %s", getEnvOrDefault("PRODUCT_SERVICE_URL", "http://localhost:8002"))
	log.Printf("   - Order Service: %s", getEnvOrDefault("ORDER_SERVICE_URL", "http://localhost:8003"))

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start gateway: %v", err)
	}
}

// getEnvOrDefault returns environment variable value or default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
