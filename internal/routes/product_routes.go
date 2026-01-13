package routes

import (
	"farmer-to-buyer-portal/internal/handlers"
	"farmer-to-buyer-portal/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupProductRoutes registers product routes
func SetupProductRoutes(rg *gin.RouterGroup) {
	products := rg.Group("/products")
	{
		// Public routes
		products.GET("", handlers.GetProducts)
		products.GET("/:id", handlers.GetProduct)

		// Protected routes (farmer only)
		products.POST("", middleware.AuthRequired(), handlers.CreateProduct)
		products.GET("/me", middleware.AuthRequired(), handlers.GetMyProducts)
		products.PUT("/:id", middleware.AuthRequired(), handlers.UpdateProduct)
		products.DELETE("/:id", middleware.AuthRequired(), handlers.DeleteProduct)
	}
}
