package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"authservice/utils"
)

// AuthMiddleware validates JWT from the Authorization header and sets user info in context.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(strings.ToLower(auth), "bearer ") {
			utils.JSONError(c, http.StatusUnauthorized, "missing or invalid authorization header")
			c.Abort()
			return
		}
		tokenStr := strings.TrimSpace(auth[len("Bearer "):])
		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			utils.JSONError(c, http.StatusUnauthorized, "invalid token: "+err.Error())
			c.Abort()
			return
		}
		// Store claims in context
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// RequireRoles ensures the authenticated user has one of the allowed roles.
func RequireRoles(roles ...string) gin.HandlerFunc {
	allowed := map[string]struct{}{}
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(c *gin.Context) {
		val, exists := c.Get("role")
		if !exists {
			utils.JSONError(c, http.StatusForbidden, "role not found in context")
			c.Abort()
			return
		}
		role, _ := val.(string)
		if _, ok := allowed[role]; !ok {
			utils.JSONError(c, http.StatusForbidden, "insufficient role")
			c.Abort()
			return
		}
		c.Next()
	}
}
