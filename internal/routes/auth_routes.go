package routes

import (
	"farmer-to-buyer-portal/internal/handlers"
	"farmer-to-buyer-portal/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes registers authentication routes
func SetupAuthRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
	}

	// Protected route
	rg.GET("/me", middleware.AuthRequired(), handlers.Me)
}
