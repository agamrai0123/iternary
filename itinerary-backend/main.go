package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/itinerary-backend/itinerary"
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
	config, err := itinerary.LoadConfig("config/config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger := itinerary.NewLogger(config)
	logger.Info("Starting Itinerary Service")

	// Initialize metrics
	// metrics := itinerary.NewMetrics()  // DISABLED - not used

	// Initialize database
	db, err := itinerary.NewDatabase(config, logger)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	logger.Info("Database connection successful")

	// Initialize service
	svc := itinerary.NewService(db, logger)
	logger.Info("Service initialized")

	// Initialize auth service
	authService := itinerary.NewAuthService(db, logger)
	logger.Info("Auth service initialized")

	// Initialize metrics
	metrics := itinerary.NewMetrics()
	logger.Info("Metrics initialized")

	// Initialize router with all services
	router := itinerary.SetupRoutes(svc, logger, metrics, authService)
	logger.Info("Routes configured")

	// Run server
	port := config.Server.Port
	logger.Info("Starting server on port " + port)
	if err := router.Run(":" + port); err != nil {
		logger.Error("Server failed to start: " + err.Error())
		os.Exit(1)
	}
}
