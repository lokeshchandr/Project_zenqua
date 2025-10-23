package service

import (
	"errors"
	"time"

	"productservice/internal/models"
	"productservice/internal/repo"
)

// ProductService defines business logic for product operations.
type ProductService interface {
	CreateProduct(req *models.CreateProductRequest) (*models.Product, error)
	GetAllProducts() ([]models.Product, error)
	GetProductByID(id uint) (*models.Product, error)
	UpdateProduct(id uint, req *models.UpdateProductRequest) (*models.Product, error)
	UpdateStock(id uint, quantity int) (*models.Product, error)
	DeleteProduct(id uint) error
	GetProductsBySellerID(sellerID uint) ([]models.Product, error)
}

// productService implements ProductService.
type productService struct {
	repo repo.ProductRepository
}

// NewProductService creates a new product service instance.
func NewProductService(repo repo.ProductRepository) ProductService {
	return &productService{repo: repo}
}

// CreateProduct creates a new product in the system.
func (s *productService) CreateProduct(req *models.CreateProductRequest) (*models.Product, error) {
	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		SellerID:    req.SellerID,
		Quantity:    req.Quantity,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.Create(product); err != nil {
		return nil, err
	}

	return product, nil
}

// GetAllProducts retrieves all products.
func (s *productService) GetAllProducts() ([]models.Product, error) {
	return s.repo.FindAll()
}

// GetProductsBySellerID retrieves all products for a specific seller.
func (s *productService) GetProductsBySellerID(sellerID uint) ([]models.Product, error) {
	return s.repo.FindBySellerID(sellerID)
}

// GetProductByID retrieves a single product by ID.
func (s *productService) GetProductByID(id uint) (*models.Product, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}
	return product, nil
}

// UpdateProduct updates product information (name, description, price, quantity).
func (s *productService) UpdateProduct(id uint, req *models.UpdateProductRequest) (*models.Product, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}

	// Update only provided fields
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Quantity != nil {
		product.Quantity = *req.Quantity
	}

	product.UpdatedAt = time.Now()

	if err := s.repo.Update(product); err != nil {
		return nil, err
	}

	return product, nil
}

// UpdateStock updates only the stock quantity of a product.
func (s *productService) UpdateStock(id uint, quantity int) (*models.Product, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}

	product.Quantity = quantity
	product.UpdatedAt = time.Now()

	if err := s.repo.Update(product); err != nil {
		return nil, err
	}

	return product, nil
}

// DeleteProduct removes a product from the system.
func (s *productService) DeleteProduct(id uint) error {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("product not found")
	}

	return s.repo.Delete(id)
}
