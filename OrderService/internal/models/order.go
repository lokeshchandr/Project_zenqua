package models

import "time"

// Order statuses
const (
	StatusPending   = "Pending"
	StatusCompleted = "Completed"
	StatusCancelled = "Cancelled"
)

// Payment methods
const (
	PaymentCOD    = "COD"
	PaymentOnline = "Online"
)

// Order represents an order in the system.
type Order struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        uint      `gorm:"not null;index" json:"user_id"`
	ProductID     uint      `gorm:"not null;index" json:"product_id"`
	Quantity      int       `gorm:"not null" json:"quantity"`
	TotalAmount   float64   `gorm:"not null" json:"total_amount"`
	PaymentMethod string    `gorm:"type:text;not null" json:"payment_method"`
	Status        string    `gorm:"type:text;not null;default:Pending" json:"status"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// CreateOrderRequest represents the request to create an order.
type CreateOrderRequest struct {
	ProductID     uint   `json:"product_id" binding:"required"`
	Quantity      int    `json:"quantity" binding:"required,gt=0"`
	PaymentMethod string `json:"payment_method" binding:"required,oneof=COD Online"`
}

// UpdateOrderStatusRequest represents the request to update order status.
type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=Pending Completed Cancelled"`
}

// Product represents product info from Product Service (for validation).
type Product struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	SellerID uint    `json:"seller_id"`
}
