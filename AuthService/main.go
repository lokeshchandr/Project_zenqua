package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"authservice/config"
	"authservice/database"
	"authservice/handlers"
	"authservice/middleware"
)

// bootstrap sets up configuration, database, routes, and starts the server.
func bootstrap() {
	// Load environment configuration (optionally from .env)
	if err := config.Load(); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Initialize database and run migrations/seeders
	if err := database.InitDatabase(); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	r := gin.Default()

	// Public routes
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	// Protected routes: require valid JWT
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		// Notifications for any logged in user
		auth.GET("/notifications", handlers.GetUnreadNotifications)

		// Super Admin only: approve Admins
		auth.PUT("/approve-admin/:id", middleware.RequireRoles("superadmin"), handlers.ApproveAdmin)
		// to get the list of pending admin approval requests
		auth.GET("/approve-request", middleware.RequireRoles("superadmin"), handlers.GetPendingAdminApprovals)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func main() {
	bootstrap()
}
