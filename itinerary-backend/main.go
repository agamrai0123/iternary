package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/itinerary-backend/itinerary/auth"
	"github.com/yourusername/itinerary-backend/itinerary/config"
	"github.com/yourusername/itinerary-backend/itinerary/middleware"
	"github.com/yourusername/itinerary-backend/itinerary/routes"
	"github.com/yourusername/itinerary-backend/itinerary/service"
	"github.com/yourusername/itinerary-backend/itinerary/utils"
)

func main() {
	// Set Gin mode from environment variable (required for production)
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = "debug" // Default to debug for development
	}
	gin.SetMode(ginMode)
	log.Printf("Gin mode set to: %s", ginMode)

	// Load configuration
	cfg, err := config.LoadConfig("config/config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger := utils.NewLogger(cfg)
	logger.Info("Starting Itinerary Service")


	// Initialize metrics
	metrics := middleware.NewMetrics()

	// Initialize database
	db, err := service.NewDatabase(cfg, logger)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	logger.Info("Database connection successful")

	// Initialize service
	svc := service.NewService(db, logger)

	// Initialize auth service
	authService := auth.NewAuthService(db, logger)

	// Initialize router
	router := routes.SetupRoutes(svc, logger, metrics, authService)

	// Run server
	port := config.Server.Port
	logger.Info("Starting server on port " + port)
	if err := router.Run(":" + port); err != nil {
		logger.Error("Server failed to start: " + err.Error())
		os.Exit(1)
	}
}
