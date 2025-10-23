package middleware

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware logs all incoming requests with user context
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get user_id from context if available
		userID, exists := c.Get("user_id")
		userIDStr := "anonymous"
		if exists {
			userIDStr = fmt.Sprintf("%v", userID)
		}

		// Log request details
		log.Printf("[%s] %s %s | Status: %d | Latency: %v | User: %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.Request.Proto,
			c.Writer.Status(),
			latency,
			userIDStr,
		)
	}
}
