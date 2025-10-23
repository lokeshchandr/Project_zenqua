package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"productservice/internal/db"
	"productservice/internal/handlers"
	"productservice/internal/middleware"
	"productservice/internal/repo"
	"productservice/internal/service"
)

func main() {
	// Load .env file if present
	// envPath := filepath.Join("..", ".env")
	err := godotenv.Load()
	if err != nil {
		log.Print("No .env file found")
	}
	// Load JWT secret from environment variable
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	// Initialize database and run migrations
	database, err := db.InitDB("product.db")
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	// Initialize layers (dependency injection)
	productRepo := repo.NewProductRepository(database)
	productService := service.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// Setup Gin router
	router := gin.Default()

	// Public routes - anyone can view products
	router.GET("/products", productHandler.ListProducts)
	router.GET("/products/:id", productHandler.GetProduct)

	// Protected routes - Admin/Super Admin only
	adminRoutes := router.Group("/")
	adminRoutes.Use(middleware.JWTAuth(jwtSecret), middleware.AdminOnly())
	{
		adminRoutes.POST("/products", productHandler.CreateProduct)
		adminRoutes.PATCH("/products/:id", productHandler.UpdateProduct)
		adminRoutes.PATCH("/products/:id/stock", productHandler.UpdateStock)
		adminRoutes.DELETE("/products/:id", productHandler.DeleteProduct)
		// adminRoutes.GET("/products/admin", productHandler.GetAdminProducts)
		adminRoutes.GET("/allProducts", productHandler.GetSalerProducts)
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8002"
	}

	log.Printf("Product service running on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
