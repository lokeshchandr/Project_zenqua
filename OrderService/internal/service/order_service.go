package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"orderservice/internal/models"
	"orderservice/internal/repo"
)

// OrderService defines business logic for order operations.
type OrderService interface {
	CreateOrder(userID uint, req *models.CreateOrderRequest) (*models.Order, error)
	GetAllOrders() ([]models.Order, error)
	GetOrdersByUserID(userID uint) ([]models.Order, error)
	GetOrderByID(id uint) (*models.Order, error)
	UpdateOrderStatus(id uint, status string) error
}

// orderService implements OrderService.
type orderService struct {
	repo              repo.OrderRepository
	productServiceURL string
}

// NewOrderService creates a new order service instance.
func NewOrderService(repo repo.OrderRepository, productServiceURL string) OrderService {
	return &orderService{
		repo:              repo,
		productServiceURL: productServiceURL,
	}
}

// CreateOrder creates a new order with product validation and quantity deduction.
func (s *orderService) CreateOrder(userID uint, req *models.CreateOrderRequest) (*models.Order, error) {
	// 1. Validate product availability from Product Service
	product, err := s.getProductFromService(req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch product: %w", err)
	}

	if product == nil {
		return nil, errors.New("product not found")
	}

	// 2. Check stock availability
	if product.Quantity < req.Quantity {
		return nil, fmt.Errorf("insufficient stock: available=%d, requested=%d", product.Quantity, req.Quantity)
	}

	// 3. Calculate total amount
	totalAmount := product.Price * float64(req.Quantity)

	// 4. Determine order status based on payment method
	status := models.StatusPending
	if req.PaymentMethod == models.PaymentOnline {
		status = models.StatusCompleted
	}

	// 5. Create order
	order := &models.Order{
		UserID:        userID,
		ProductID:     req.ProductID,
		Quantity:      req.Quantity,
		TotalAmount:   totalAmount,
		PaymentMethod: req.PaymentMethod,
		Status:        status,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.repo.Create(order); err != nil {
		return nil, err
	}

	// 6. Deduct product quantity (call Product Service)
	if err := s.deductProductQuantity(req.ProductID, req.Quantity); err != nil {
		// Log error but don't fail the order (can be handled by background job)
		fmt.Printf("Warning: failed to deduct product quantity: %v\n", err)
	}

	return order, nil
}

// GetAllOrders retrieves all orders (Admin view).
func (s *orderService) GetAllOrders() ([]models.Order, error) {
	return s.repo.FindAll()
}

// GetOrdersByUserID retrieves orders for a specific user.
func (s *orderService) GetOrdersByUserID(userID uint) ([]models.Order, error) {
	return s.repo.FindByUserID(userID)
}

// GetOrderByID retrieves a single order by ID.
func (s *orderService) GetOrderByID(id uint) (*models.Order, error) {
	order, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if order == nil {
		return nil, errors.New("order not found")
	}
	return order, nil
}

// UpdateOrderStatus updates the status of an order.
func (s *orderService) UpdateOrderStatus(id uint, status string) error {
	order, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if order == nil {
		return errors.New("order not found")
	}

	return s.repo.UpdateStatus(id, status)
}

// getProductFromService fetches product details from Product Service.
func (s *orderService) getProductFromService(productID uint) (*models.Product, error) {
	url := fmt.Sprintf("%s/products/%d", s.productServiceURL, productID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("product service unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("product service returned status %d", resp.StatusCode)
	}

	var result struct {
		Product models.Product `json:"product"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode product response: %w", err)
	}

	return &result.Product, nil
}

// deductProductQuantity calls Product Service to reduce stock.
func (s *orderService) deductProductQuantity(productID uint, quantity int) error {
	url := fmt.Sprintf("%s/products/%d/stock", s.productServiceURL, productID)

	// Calculate new quantity (fetch current, subtract ordered)
	product, err := s.getProductFromService(productID)
	if err != nil || product == nil {
		return err
	}

	newQuantity := product.Quantity - quantity
	if newQuantity < 0 {
		newQuantity = 0
	}

	// Prepare request body
	bodyStr := fmt.Sprintf(`{"quantity": %d}`, newQuantity)
	req, err := http.NewRequest("PATCH", url, strings.NewReader(bodyStr))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	// For now, skip authentication (in production, use service-to-service auth)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update product quantity: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("product service returned status %d", resp.StatusCode)
	}

	return nil
}
