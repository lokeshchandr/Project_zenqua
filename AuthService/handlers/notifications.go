package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"authservice/database"
	"authservice/models"
	"authservice/utils"
)

// GetUnreadNotifications returns unread notifications for the logged-in user.
func GetUnreadNotifications(c *gin.Context) {
	v, ok := c.Get("user_id")
	if !ok {
		utils.JSONError(c, http.StatusUnauthorized, "missing user context")
		return
	}
	userID, _ := v.(uint)

	var notifs []models.Notification
	if err := database.DB.Where("user_id = ? AND is_read = ?", userID, false).Order("created_at DESC").Find(&notifs).Error; err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "failed to fetch notifications")
		return
	}
	utils.JSONOK(c, http.StatusOK, gin.H{"notifications": notifs})
}
