package routes

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"apigateway/internal/middleware"
	"apigateway/internal/proxy"
)

// SetupRoutes configures all routes and middlewares for the gateway
func SetupRoutes(router *gin.Engine) {
	// Get backend service URLs from environment
	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	if authServiceURL == "" {
		authServiceURL = "http://localhost:8001"
	}

	productServiceURL := os.Getenv("PRODUCT_SERVICE_URL")
	if productServiceURL == "" {
		productServiceURL = "http://localhost:8002"
	}

	orderServiceURL := os.Getenv("ORDER_SERVICE_URL")
	if orderServiceURL == "" {
		orderServiceURL = "http://localhost:8003"
	}

	// Configure CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // In production, specify allowed origins
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-User-ID", "X-User-Role"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Add request logging middleware
	router.Use(middleware.LoggingMiddleware())

	// Check endpoint
	router.GET("/Check", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"services": gin.H{
				"auth":    authServiceURL,
				"product": productServiceURL,
				"order":   orderServiceURL,
			},
		})
	})

	// Auth Service routes - Public (no authentication needed for login/register)
	authGroup := router.Group("/auth")
	{
		authGroup.Any("/*path", proxy.StripPrefixProxy("/auth", authServiceURL))
	}

	// Admin routes (for approving users, notifications, etc.) - Requires authentication
	// Strip the /admin prefix and forward to AuthService root, since backend admin routes are not prefixed
	adminGroup := router.Group("/admin")
	adminGroup.Use(middleware.AuthMiddleware(), middleware.AdminOnlyMiddleware())
	{
		adminGroup.Any("/*path", proxy.StripPrefixProxy("/admin", authServiceURL))
	}

	// Notifications - Requires authentication; preserve original path and proxy to AuthService
	notifGroup := router.Group("/notifications")
	notifGroup.Use(middleware.AuthMiddleware())
	{
		notifGroup.Any("/*path", proxy.ProxyHandler(authServiceURL))
	}

	// Product Service routes
	productGroup := router.Group("/products")
	{
		// Public routes (list and get single product)
		productGroup.GET("", proxy.ProxyHandler(productServiceURL))
		productGroup.GET("/:id", proxy.ProxyHandler(productServiceURL))

		// Protected routes - Admin only (create, update, delete)
		protectedProducts := productGroup.Group("")
		protectedProducts.Use(middleware.AuthMiddleware(), middleware.AdminOnlyMiddleware())
		{
			protectedProducts.POST("", proxy.ProxyHandler(productServiceURL))
			protectedProducts.PATCH("/:id", proxy.ProxyHandler(productServiceURL))
			protectedProducts.PATCH("/:id/stock", proxy.ProxyHandler(productServiceURL))
			protectedProducts.DELETE("/:id", proxy.ProxyHandler(productServiceURL))
			// Seller's own products - match backend route /allProducts
			protectedProducts.GET("/allProducts", proxy.StripPrefixProxy("/products", productServiceURL))
		}
	}

	// Order Service routes - All require authentication
	orderGroup := router.Group("/orders")
	orderGroup.Use(middleware.AuthMiddleware())
	{
		// User can create and view orders
		orderGroup.POST("", proxy.ProxyHandler(orderServiceURL))
		orderGroup.GET("", proxy.ProxyHandler(orderServiceURL))
		orderGroup.GET("/:id", proxy.ProxyHandler(orderServiceURL))

		// Admin only - update order status
		adminOrders := orderGroup.Group("")
		adminOrders.Use(middleware.AdminOnlyMiddleware())
		{
			adminOrders.PATCH("/:id/status", proxy.ProxyHandler(orderServiceURL))
		}
	}
}
