package models

import "time"

// Roles constants for RBAC
const (
	RoleSuperAdmin = "superadmin"
	RoleAdmin      = "saler"
	RoleUser       = "user"
)

// User represents an application user stored in the users table.
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"size:100;not null" json:"name"`
	Email        string    `gorm:"size:120;uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"size:255;not null" json:"-"`
	Role         string    `gorm:"size:20;not null;default:user" json:"role"`
	IsApproved   bool      `gorm:"not null;default:false" json:"is_approved"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	GSTNum       string    `gorm:"column:gstnumber" json:"gstnumber"`
}
