package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"authservice/database"
	"authservice/models"
	"authservice/utils"
)

// ApproveAdmin allows Super Admin to approve an Admin user by ID and creates a notification.
func ApproveAdmin(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		utils.JSONError(c, http.StatusBadRequest, "invalid user id")
		return
	}
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		utils.JSONError(c, http.StatusNotFound, "user not found")
		return
	}
	if user.Role != models.RoleAdmin {
		utils.JSONError(c, http.StatusBadRequest, "only admin users require approval")
		return
	}
	if user.IsApproved {
		utils.JSONOK(c, http.StatusOK, gin.H{"message": "admin already approved"})
		return
	}
	// Approve and create notification in a transaction
	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		user.IsApproved = true
		if err := tx.Save(&user).Error; err != nil {
			return err
		}
		notif := models.Notification{UserID: user.ID, Message: "Your admin account has been approved"}
		if err := tx.Create(&notif).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "failed to approve admin")
		return
	}

	utils.JSONOK(c, http.StatusOK, gin.H{"message": "admin approved"})
}

// GetPendingAdminApprovals returns a list of Admin users pending approval.
func GetPendingAdminApprovals(c *gin.Context) {
	var admins []models.User
	if err := database.DB.Where("role = ? AND is_approved = ?", models.RoleAdmin, false).Find(&admins).Error; err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "failed to retrieve pending admin approvals")
		return
	}
	utils.JSONOK(c, http.StatusOK, admins)
}
