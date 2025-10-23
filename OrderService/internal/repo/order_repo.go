package repo

import (
	"errors"

	"gorm.io/gorm"

	"orderservice/internal/models"
)

// OrderRepository defines the interface for order data operations.
type OrderRepository interface {
	Create(order *models.Order) error
	FindAll() ([]models.Order, error)
	FindByID(id uint) (*models.Order, error)
	FindByUserID(userID uint) ([]models.Order, error)
	UpdateStatus(id uint, status string) error
	Update(order *models.Order) error
}

// orderRepo implements OrderRepository using GORM.
type orderRepo struct {
	db *gorm.DB
}

// NewOrderRepository creates a new order repository instance.
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepo{db: db}
}

// Create inserts a new order into the database.
func (r *orderRepo) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

// FindAll retrieves all orders from the database.
func (r *orderRepo) FindAll() ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Order("created_at DESC").Find(&orders).Error
	return orders, err
}

// FindByID retrieves an order by its ID.
func (r *orderRepo) FindByID(id uint) (*models.Order, error) {
	var order models.Order
	err := r.db.First(&order, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &order, nil
}

// FindByUserID retrieves all orders for a specific user.
func (r *orderRepo) FindByUserID(userID uint) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&orders).Error
	return orders, err
}

// UpdateStatus updates the status of an order.
func (r *orderRepo) UpdateStatus(id uint, status string) error {
	return r.db.Model(&models.Order{}).Where("id = ?", id).Update("status", status).Error
}

// Update saves changes to an existing order.
func (r *orderRepo) Update(order *models.Order) error {
	return r.db.Save(order).Error
}
