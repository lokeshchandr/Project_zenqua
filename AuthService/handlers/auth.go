package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"authservice/database"
	"authservice/models"
	"authservice/utils"
)

// RegisterRequest represents the expected payload for user registration.
type RegisterRequest struct {
	Name      string `json:"name" binding:"required,min=2"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=6"`
	Role      string `json:"role" binding:"omitempty,oneof=superadmin saler user"`
	GSTNumber string `json:"gst_number" binding:"omitempty"` //only for the admin (Salers)
}

// Register registers a new user. Admins require approval; Super Admin creation is restricted by seeding.
func Register(c *gin.Context) {
	var req RegisterRequest
	if !utils.BindJSONOrAbort(c, &req) {
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	if req.Role == "" {
		req.Role = models.RoleUser
	}
	if req.Role == models.RoleSuperAdmin {
		utils.JSONError(c, http.StatusForbidden, "cannot self-register as superadmin")
		return
	}
	if req.Role == models.RoleAdmin {
		if strings.TrimSpace(req.GSTNumber) == "" {
			utils.JSONError(c, http.StatusBadRequest, "please provide gst number")
			return
		}
	}

	// Check if email exists
	var existing models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		utils.JSONError(c, http.StatusBadRequest, "email already in use")
		return
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "failed to hash password")
		return
	}

	user := models.User{
		Name:         strings.TrimSpace(req.Name),
		Email:        req.Email,
		PasswordHash: hash,
		Role:         req.Role,
		IsApproved:   req.Role != models.RoleAdmin, // Admins require approval
		CreatedAt:    time.Now(),
		GSTNum:       strings.TrimSpace(req.GSTNumber),
	}
	if err := database.DB.Create(&user).Error; err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "failed to create user")
		return
	}

	utils.JSONOK(c, http.StatusCreated, gin.H{
		"message":     "user registered successfully",
		"user_id":     user.ID,
		"is_approved": user.IsApproved,
		"role":        user.Role,
	})
}

// LoginRequest represents the expected payload for login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login authenticates the user and returns a JWT token if approved.
func Login(c *gin.Context) {
	var req LoginRequest
	if !utils.BindJSONOrAbort(c, &req) {
		return
	}
	email := strings.ToLower(strings.TrimSpace(req.Email))

	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		utils.JSONError(c, http.StatusUnauthorized, "invalid email or password")
		return
	}
	if !utils.CheckPassword(user.PasswordHash, req.Password) {
		utils.JSONError(c, http.StatusUnauthorized, "invalid email or password")
		return
	}
	// For admin role, ensure approved
	if user.Role == models.RoleAdmin && !user.IsApproved {
		utils.JSONError(c, http.StatusForbidden, "admin not approved yet")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, "failed to generate token")
		return
	}

	utils.JSONOK(c, http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":          user.ID,
			"name":        user.Name,
			"email":       user.Email,
			"role":        user.Role,
			"is_approved": user.IsApproved,
		},
	})
}
