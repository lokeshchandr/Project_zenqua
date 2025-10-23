package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"productservice/internal/middleware"
	"productservice/internal/models"
	"productservice/internal/service"
)

// ProductHandler holds the product service dependency.
type ProductHandler struct {
	service service.ProductService
}

// NewProductHandler creates a new product handler instance.
func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// ListProducts handles GET /products - retrieves all products (public).
func (h *ProductHandler) ListProducts(c *gin.Context) {
	products, err := h.service.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

// GetProduct handles GET /products/:id - retrieves a single product by ID (public).
func (h *ProductHandler) GetProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	product, err := h.service.GetProductByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product": product})
}

// GetSalerProducts handles GET /allProducts - retrieves all products for a specific seller (Saler only).
func (h *ProductHandler) GetSalerProducts(c *gin.Context) {
	// Get the seller ID from the JWT token (assuming it's included in the claims)
	sellerID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	products, err := h.service.GetProductsBySellerID(sellerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch products"})
		return
	}

	message := "products fetched successfully"
	if len(products) == 0 {
		message = "You have no products yet"
	}

	c.JSON(http.StatusOK, gin.H{
		"Total Items": len(products),
		"Products":    products,
		"message":     message,
	})
}

// // GetSalerProducts - DEPRECATED (use GetMyProducts instead)
// func (h *ProductHandler) GetSalerProducts(c *gin.Context) {
//     h.GetMyProducts(c)
// }

// CreateProduct handles POST /products - creates a new product (Saler only).
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req models.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sellerID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	req.SellerID = sellerID

	product, err := h.service.CreateProduct(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "product created successfully",
		"product": product,
	})
}

// UpdateProduct handles PATCH /products/:id - updates product info (Saler only).
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	sellerID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	existingProduct, err := h.service.GetProductByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	if existingProduct.SellerID != sellerID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Wrong Product ID"})
		return
	}

	var req models.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.service.UpdateProduct(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "product updated successfully",
		"product": product,
	})
}

// UpdateStock handles PATCH /products/:id/stock - updates product stock quantity (Saler only).
func (h *ProductHandler) UpdateStock(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	sellerID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	existingProduct, err := h.service.GetProductByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	if existingProduct.SellerID != sellerID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Product Id is incorrect"})
		return
	}

	var req models.UpdateStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.service.UpdateStock(uint(id), req.Quantity)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "stock updated successfully",
		"product": product,
	})
}

// DeleteProduct handles DELETE /products/:id - deletes a product (Saler only).
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}
	sellerID, err := middleware.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	existingProduct, err := h.service.GetProductByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	if existingProduct.SellerID != sellerID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you can only delete your own products"})
		return
	}

	if err := h.service.DeleteProduct(uint(id)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted successfully"})
}
