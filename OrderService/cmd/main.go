package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"orderservice/internal/db"
	"orderservice/internal/handlers"
	"orderservice/internal/middleware"
	"orderservice/internal/repo"
	"orderservice/internal/service"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Validate required environment variables
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	productServiceURL := os.Getenv("PRODUCT_SERVICE_URL")
	if productServiceURL == "" {
		productServiceURL = "http://localhost:8002" // default
		log.Printf("Using default PRODUCT_SERVICE_URL: %s", productServiceURL)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8003"
	}

	// Initialize database
	database := db.InitDB()

	// Initialize layers (dependency injection)
	orderRepo := repo.NewOrderRepository(database)
	orderService := service.NewOrderService(orderRepo, productServiceURL)
	orderHandler := handlers.NewOrderHandler(orderService)

	// Setup Gin router
	router := gin.Default()

	// Check
	router.GET("/Check", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "order-service"})
	})

	// Protected routes - require authentication
	authGroup := router.Group("/")
	authGroup.Use(middleware.AuthMiddleware())
	{
		// Users can create orders and view their own
		authGroup.POST("/orders", orderHandler.CreateOrder)
		authGroup.GET("/orders", orderHandler.GetOrders) // Role-based filtering inside handler
		authGroup.GET("/orders/:id", orderHandler.GetOrder)

		// Admin/Super Admin only - update order status
		adminGroup := authGroup.Group("/")
		adminGroup.Use(middleware.AdminOnlyMiddleware())
		{
			adminGroup.PATCH("/orders/:id/status", orderHandler.UpdateOrderStatus)
		}
	}

	log.Printf("Order service running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
