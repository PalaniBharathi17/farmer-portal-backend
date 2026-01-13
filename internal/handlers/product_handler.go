package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"farmer-to-buyer-portal/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateProductRequest represents the request payload for creating a product
type CreateProductRequest struct {
	CropName     string  `json:"crop_name" binding:"required"`
	Quantity     float64 `json:"quantity" binding:"required,gt=0"`
	Unit         string  `json:"unit" binding:"required"`
	PricePerUnit float64 `json:"price_per_unit" binding:"required,gt=0"`
	State        string  `json:"state" binding:"required"`
	City         string  `json:"city" binding:"required"`
	Pincode      string  `json:"pincode" binding:"required"`
}

// UpdateProductRequest represents the request payload for updating a product
type UpdateProductRequest struct {
	Quantity     *float64 `json:"quantity"`
	PricePerUnit *float64 `json:"price_per_unit"`
	Status       *string  `json:"status"`
}

// ProductResponse represents the product data in API responses
type ProductResponse struct {
	ID           string  `json:"id"`
	FarmerID     string  `json:"farmer_id"`
	CropName     string  `json:"crop_name"`
	Quantity     float64 `json:"quantity"`
	Unit         string  `json:"unit"`
	PricePerUnit float64 `json:"price_per_unit"`
	State        string  `json:"state"`
	City         string  `json:"city"`
	Pincode      string  `json:"pincode"`
	Status       string  `json:"status"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

// toProductResponse converts a Product model to ProductResponse
func toProductResponse(p models.Product) ProductResponse {
	return ProductResponse{
		ID:           p.ID,
		FarmerID:     p.FarmerID,
		CropName:     p.CropName,
		Quantity:     p.Quantity,
		Unit:         p.Unit,
		PricePerUnit: p.PricePerUnit,
		State:        p.State,
		City:         p.City,
		Pincode:      p.Pincode,
		Status:       p.Status,
		CreatedAt:    p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// CreateProduct handles POST /api/v1/products (farmer only)
func CreateProduct(c *gin.Context) {
	// Check if user is farmer
	role := c.MustGet("role").(string)
	if role != "farmer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only farmers can create products"})
		return
	}

	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(string)

	product := models.Product{
		FarmerID:     userID,
		CropName:     req.CropName,
		Quantity:     req.Quantity,
		Unit:         req.Unit,
		PricePerUnit: req.PricePerUnit,
		State:        req.State,
		City:         req.City,
		Pincode:      req.Pincode,
		Status:       "active",
	}

	if err := db.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, toProductResponse(product))
}

// GetProducts handles GET /api/v1/products (public with filters)
func GetProducts(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	query := db.Model(&models.Product{}).Where("status = ?", "active")

	// Apply filters
	if cropName := c.Query("crop_name"); cropName != "" {
		query = query.Where("crop_name LIKE ?", "%"+cropName+"%")
	}
	if pincode := c.Query("pincode"); pincode != "" {
		query = query.Where("pincode = ?", pincode)
	}
	if city := c.Query("city"); city != "" {
		query = query.Where("city = ?", city)
	}
	if state := c.Query("state"); state != "" {
		query = query.Where("state = ?", state)
	}
	if minPrice := c.Query("min_price"); minPrice != "" {
		if min, err := strconv.ParseFloat(minPrice, 64); err == nil {
			query = query.Where("price_per_unit >= ?", min)
		}
	}
	if maxPrice := c.Query("max_price"); maxPrice != "" {
		if max, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			query = query.Where("price_per_unit <= ?", max)
		}
	}

	// Sort by created_at desc
	query = query.Order("created_at DESC")

	var products []models.Product
	if err := query.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	responses := make([]ProductResponse, len(products))
	for i, p := range products {
		responses[i] = toProductResponse(p)
	}

	c.JSON(http.StatusOK, responses)
}

// GetProduct handles GET /api/v1/products/:id (public)
func GetProduct(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	productID := c.Param("id")

	var product models.Product
	if err := db.Where("id = ?", productID).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, toProductResponse(product))
}

// GetMyProducts handles GET /api/v1/products/me (farmer only)
func GetMyProducts(c *gin.Context) {
	// Check if user is farmer
	role := c.MustGet("role").(string)
	if role != "farmer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only farmers can access this endpoint"})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(string)

	var products []models.Product
	if err := db.Where("farmer_id = ?", userID).Order("created_at DESC").Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	responses := make([]ProductResponse, len(products))
	for i, p := range products {
		responses[i] = toProductResponse(p)
	}

	c.JSON(http.StatusOK, responses)
}

// UpdateProduct handles PUT /api/v1/products/:id (farmer only, owner only)
func UpdateProduct(c *gin.Context) {
	// Check if user is farmer
	role := c.MustGet("role").(string)
	if role != "farmer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only farmers can update products"})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	productID := c.Param("id")
	userID := c.MustGet("user_id").(string)

	// Check if product exists and belongs to the user
	var product models.Product
	if err := db.Where("id = ? AND farmer_id = ?", productID, userID).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found or you don't have permission to update it"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update only provided fields
	updates := make(map[string]interface{})
	if req.Quantity != nil {
		if *req.Quantity <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Quantity must be greater than 0"})
			return
		}
		updates["quantity"] = *req.Quantity
	}
	if req.PricePerUnit != nil {
		if *req.PricePerUnit <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Price per unit must be greater than 0"})
			return
		}
		updates["price_per_unit"] = *req.PricePerUnit
	}
	if req.Status != nil {
		validStatuses := []string{"active", "closed", "sold"}
		valid := false
		for _, s := range validStatuses {
			if *req.Status == s {
				valid = true
				break
			}
		}
		if !valid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Status must be one of: active, closed, sold"})
			return
		}
		updates["status"] = *req.Status
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No fields to update"})
		return
	}

	if err := db.Model(&product).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	// Reload product to get updated values
	db.Where("id = ?", productID).First(&product)
	c.JSON(http.StatusOK, toProductResponse(product))
}

// DeleteProduct handles DELETE /api/v1/products/:id (farmer only, owner only)
func DeleteProduct(c *gin.Context) {
	// Check if user is farmer
	role := c.MustGet("role").(string)
	if role != "farmer" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only farmers can delete products"})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	productID := c.Param("id")
	userID := c.MustGet("user_id").(string)

	// Check if product exists and belongs to the user
	var product models.Product
	if err := db.Where("id = ? AND farmer_id = ?", productID, userID).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found or you don't have permission to delete it"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if err := db.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
