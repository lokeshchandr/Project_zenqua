package repo

import (
	"errors"

	"gorm.io/gorm"

	"productservice/internal/models"
)

// ProductRepository defines the interface for product data operations.
type ProductRepository interface {
	Create(product *models.Product) error
	FindAll() ([]models.Product, error)
	FindByID(id uint) (*models.Product, error)
	Update(product *models.Product) error
	Delete(id uint) error
	FindBySellerID(sellerID uint) ([]models.Product, error)
}

// productRepo implements ProductRepository using GORM.
type productRepo struct {
	db *gorm.DB
}

// NewProductRepository creates a new product repository instance.
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepo{db: db}
}

// Create inserts a new product into the database.
func (r *productRepo) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

// FindAll retrieves all products from the database.
func (r *productRepo) FindAll() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Order("created_at DESC").Find(&products).Error
	return products, err
}

// FindByID retrieves a product by its ID.
func (r *productRepo) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if not found
		}
		return nil, err
	}
	return &product, nil
}

// Update saves changes to an existing product.
func (r *productRepo) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

// Delete removes a product from the database by ID.
func (r *productRepo) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

// FindBySellerID retrieves all products for a specific seller.
func (r *productRepo) FindBySellerID(sellerID uint) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Where("seller_id = ?", sellerID).Order("created_at DESC").Find(&products).Error
	return products, err
}
