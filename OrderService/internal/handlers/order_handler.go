package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"orderservice/internal/middleware"
	"orderservice/internal/models"
	"orderservice/internal/service"
)

// OrderHandler holds the order service dependency.
type OrderHandler struct {
	service service.OrderService
}

// NewOrderHandler creates a new order handler instance.
func NewOrderHandler(service service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

// CreateOrder handles POST /orders - creates a new order (User only).
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req models.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract user ID from JWT
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Create order
	order, err := h.service.CreateOrder(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "order created successfully",
		"order":   order,
	})
}

// GetOrders handles GET /orders - get orders based on role.
// Users see only their orders; Admins see all orders.
func (h *OrderHandler) GetOrders(c *gin.Context) {
	role, err := middleware.GetRole(c)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "role not found"})
		return
	}

	// Admin/Super Admin can see all orders
	if role == "saler" || role == "superadmin" {
		orders, err := h.service.GetAllOrders()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch orders"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"total":  len(orders),
			"orders": orders,
		})
		return
	}

	// Regular users see only their orders
	userID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	orders, err := h.service.GetOrdersByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch your orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":  len(orders),
		"orders": orders,
	})
}

// GetOrder handles GET /orders/:id - get a specific order.
func (h *OrderHandler) GetOrder(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	order, err := h.service.GetOrderByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Authorization check: Users can only view their own orders
	userID, _ := middleware.GetUserID(c)
	role, _ := middleware.GetRole(c)

	if role != "saler" && role != "superadmin" && order.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you can only view your own orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": order})
}

// UpdateOrderStatus handles PATCH /orders/:id/status - updates order status (Admin only).
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	var req models.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateOrderStatus(uint(id), req.Status); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order status updated successfully"})
}
