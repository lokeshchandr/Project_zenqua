package models

import "time"

// Product represents a product in the store.
type Product struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:text;not null" json:"name" binding:"required"`
	Description string    `gorm:"type:text" json:"description"`
	Price       float64   `gorm:"not null" json:"price" binding:"required,gt=0"`
	SellerID    uint      `gorm:"not null" json:"seller_id"`
	Quantity    int       `gorm:"not null;default:0" json:"quantity" binding:"gte=0"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// CreateProductRequest represents the request to create a product.
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required,min=1"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	SellerID    uint
	Quantity    int `json:"quantity" binding:"gte=0"`
}

// UpdateProductRequest represents the request to update a product.
type UpdateProductRequest struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty" binding:"omitempty,gt=0"`
	SellerID    *uint
	Quantity    *int `json:"quantity,omitempty" binding:"omitempty,gte=0"`
}

// UpdateStockRequest represents the request to update product stock.
type UpdateStockRequest struct {
	Quantity int `json:"quantity" binding:"required,gte=0"`
}
