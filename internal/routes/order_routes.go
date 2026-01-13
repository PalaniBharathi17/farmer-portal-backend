package routes

import (
	"farmer-to-buyer-portal/internal/handlers"
	"farmer-to-buyer-portal/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupOrderRoutes registers order routes
func SetupOrderRoutes(rg *gin.RouterGroup) {
	orders := rg.Group("/orders")
	orders.Use(middleware.AuthRequired()) // All order routes require authentication
	{
		orders.POST("", handlers.CreateOrder)
		orders.GET("/:id", handlers.GetOrder)
		orders.GET("/buyer/me", handlers.GetBuyerOrders)
		orders.GET("/farmer/me", handlers.GetFarmerOrders)
		orders.PUT("/:id/status", handlers.UpdateOrderStatus)
	}
}
