package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// JSONError sends a standardized error JSON response with a message.
func JSONError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

// JSONOK sends a standard success JSON response with data.
func JSONOK(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

// BindJSONOrAbort binds JSON into obj and aborts with 400 on error.
func BindJSONOrAbort[T any](c *gin.Context, obj *T) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return false
	}
	return true
}
