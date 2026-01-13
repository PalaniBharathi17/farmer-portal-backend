package handlers

import (
	"errors"
	"net/http"

	"farmer-to-buyer-portal/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateOrderRequest represents the request payload for creating an order
type CreateOrderRequest struct {
	ProductID    string  `json:"product_id" binding:"required"`
	Quantity     float64 `json:"quantity" binding:"required,gt=0"`
	DeliveryMode string  `json:"delivery_mode" binding:"required,oneof=pickup courier"`
}

// UpdateOrderStatusRequest represents the request payload for updating order status
type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=accepted rejected shipped delivered"`
}

// OrderItemResponse represents an order item in API responses
type OrderItemResponse struct {
	ID           string  `json:"id"`
	ProductID    string  `json:"product_id"`
	Quantity     float64 `json:"quantity"`
	PricePerUnit float64 `json:"price_per_unit"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

// OrderResponse represents an order in API responses
type OrderResponse struct {
	ID           string             `json:"id"`
	BuyerID      string             `json:"buyer_id"`
	FarmerID     string             `json:"farmer_id"`
	Status       string             `json:"status"`
	DeliveryMode string             `json:"delivery_mode"`
	TotalAmount  float64            `json:"total_amount"`
	CreatedAt    string             `json:"created_at"`
	UpdatedAt    string             `json:"updated_at"`
	OrderItems   []OrderItemResponse `json:"order_items"`
}

// toOrderItemResponse converts an OrderItem model to OrderItemResponse
func toOrderItemResponse(oi models.OrderItem) OrderItemResponse {
	return OrderItemResponse{
		ID:           oi.ID,
		ProductID:    oi.ProductID,
		Quantity:     oi.Quantity,
		PricePerUnit: oi.PricePerUnit,
		CreatedAt:    oi.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    oi.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// toOrderResponse converts an Order model to OrderResponse
func toOrderResponse(order models.Order) OrderResponse {
	items := make([]OrderItemResponse, len(order.OrderItems))
	for i, item := range order.OrderItems {
		items[i] = toOrderItemResponse(item)
	}

	return OrderResponse{
		ID:           order.ID,
		BuyerID:      order.BuyerID,
		FarmerID:     order.FarmerID,
		Status:       order.Status,
		DeliveryMode: order.DeliveryMode,
		TotalAmount:  order.TotalAmount,
		CreatedAt:    order.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    order.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		OrderItems:   items,
	}
}

// CreateOrder handles POST /api/v1/orders (buyer only)
func CreateOrder(c *gin.Context) {
	// Check if user is buyer
	role := c.MustGet("role").(string)
	if role != "buyer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only buyers can place orders"})
		return
	}

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	buyerID := c.MustGet("user_id").(string)

	// Fetch product
	var product models.Product
	if err := db.Where("id = ?", req.ProductID).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Check if product is active
	if product.Status != "active" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product is not available for ordering"})
		return
	}

	// Calculate total amount
	totalAmount := req.Quantity * product.PricePerUnit

	// Create order
	order := models.Order{
		BuyerID:      buyerID,
		FarmerID:     product.FarmerID,
		Status:       "pending",
		DeliveryMode: req.DeliveryMode,
		TotalAmount:  totalAmount,
	}

	if err := db.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Create order item
	orderItem := models.OrderItem{
		OrderID:      order.ID,
		ProductID:    req.ProductID,
		Quantity:     req.Quantity,
		PricePerUnit: product.PricePerUnit,
	}

	if err := db.Create(&orderItem).Error; err != nil {
		// Rollback order creation
		db.Delete(&order)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order item"})
		return
	}

	// Load order with items for response
	var createdOrder models.Order
	if err := db.Preload("OrderItems").Where("id = ?", order.ID).First(&createdOrder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch created order"})
		return
	}

	c.JSON(http.StatusCreated, toOrderResponse(createdOrder))
}

// GetOrder handles GET /api/v1/orders/:id (buyer or farmer can access own order)
func GetOrder(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	orderID := c.Param("id")
	userID := c.MustGet("user_id").(string)
	role := c.MustGet("role").(string)

	var order models.Order
	query := db.Preload("OrderItems").Where("id = ?", orderID)

	// Check ownership based on role
	if role == "buyer" {
		query = query.Where("buyer_id = ?", userID)
	} else if role == "farmer" {
		query = query.Where("farmer_id = ?", userID)
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "Invalid role"})
		return
	}

	if err := query.First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found or you don't have permission to access it"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, toOrderResponse(order))
}

// GetBuyerOrders handles GET /api/v1/orders/buyer/me (buyer only)
func GetBuyerOrders(c *gin.Context) {
	// Check if user is buyer
	role := c.MustGet("role").(string)
	if role != "buyer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only buyers can access this endpoint"})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	buyerID := c.MustGet("user_id").(string)

	var orders []models.Order
	if err := db.Preload("OrderItems").Where("buyer_id = ?", buyerID).Order("created_at DESC").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	responses := make([]OrderResponse, len(orders))
	for i, order := range orders {
		responses[i] = toOrderResponse(order)
	}

	c.JSON(http.StatusOK, responses)
}

// GetFarmerOrders handles GET /api/v1/orders/farmer/me (farmer only)
func GetFarmerOrders(c *gin.Context) {
	// Check if user is farmer
	role := c.MustGet("role").(string)
	if role != "farmer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only farmers can access this endpoint"})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	farmerID := c.MustGet("user_id").(string)

	var orders []models.Order
	if err := db.Preload("OrderItems").Where("farmer_id = ?", farmerID).Order("created_at DESC").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	responses := make([]OrderResponse, len(orders))
	for i, order := range orders {
		responses[i] = toOrderResponse(order)
	}

	c.JSON(http.StatusOK, responses)
}

// UpdateOrderStatus handles PUT /api/v1/orders/:id/status (farmer only)
func UpdateOrderStatus(c *gin.Context) {
	// Check if user is farmer
	role := c.MustGet("role").(string)
	if role != "farmer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only farmers can update order status"})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	orderID := c.Param("id")
	farmerID := c.MustGet("user_id").(string)

	// Check if order exists and belongs to the farmer
	var order models.Order
	if err := db.Where("id = ? AND farmer_id = ?", orderID, farmerID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found or you don't have permission to update it"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	var req UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate status transitions
	currentStatus := order.Status
	newStatus := req.Status

	validTransition := false
	switch currentStatus {
	case "pending":
		if newStatus == "accepted" || newStatus == "rejected" {
			validTransition = true
		}
	case "accepted":
		if newStatus == "shipped" {
			validTransition = true
		}
	case "shipped":
		if newStatus == "delivered" {
			validTransition = true
		}
	case "rejected", "delivered":
		// Cannot transition from rejected or delivered
		validTransition = false
	}

	if !validTransition {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid status transition",
			"current_status": currentStatus,
			"requested_status": newStatus,
			"valid_transitions": getValidTransitions(currentStatus),
		})
		return
	}

	// Update status
	if err := db.Model(&order).Update("status", newStatus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	// Reload order with items for response
	var updatedOrder models.Order
	if err := db.Preload("OrderItems").Where("id = ?", orderID).First(&updatedOrder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated order"})
		return
	}

	c.JSON(http.StatusOK, toOrderResponse(updatedOrder))
}

// getValidTransitions returns valid status transitions for a given status
func getValidTransitions(status string) []string {
	switch status {
	case "pending":
		return []string{"accepted", "rejected"}
	case "accepted":
		return []string{"shipped"}
	case "shipped":
		return []string{"delivered"}
	default:
		return []string{}
	}
}
