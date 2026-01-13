package main

import (
	"log"
	"os"

	"farmer-to-buyer-portal/internal/config"
	"farmer-to-buyer-portal/internal/db"
	"farmer-to-buyer-portal/internal/models"
	"farmer-to-buyer-portal/internal/routes"
	"farmer-to-buyer-portal/internal/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// Set Gin to debug mode to show routes
	gin.SetMode(gin.DebugMode)

	// Load config first to ensure .env is loaded before reading PORT
	cfg := config.Load()

	// Initialize JWT secret
	utils.InitJWT(cfg)

	conn, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("could not start server: %v", err)
	}

	// Auto-migrate models
	if err := conn.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{}, &models.OrderItem{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	router := routes.SetupRouter(conn)

	// Print registered routes
	log.Println("INFO: Registered routes:")
	for _, route := range router.Routes() {
		log.Printf("  %s %s", route.Method, route.Path)
	}

	// PORT is loaded from environment (can be set in .env or system env)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("INFO: Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("server stopped unexpectedly: %v", err)
	}
}
