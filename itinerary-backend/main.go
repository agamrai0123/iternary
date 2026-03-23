package main

import (
	"log"
	"os"

	"github.com/yourusername/itinerary-backend/itinerary"
)

func main() {
	// Load configuration
	config, err := itinerary.LoadConfig("config/config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger := itinerary.NewLogger(config)
	logger.Info("Starting Itinerary Service")

	// Initialize metrics
	metrics := itinerary.NewMetrics()

	// Initialize database
	db, err := itinerary.NewDatabase(config, logger)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	logger.Info("Database connection successful")

	// Initialize service
	service := itinerary.NewService(db, logger)

	// Initialize auth service
	authService := itinerary.NewAuthService(db, logger)

	// Initialize router
	router := itinerary.SetupRoutes(service, logger, metrics, authService)

	// Run server
	port := config.Server.Port
	logger.Info("Starting server on port " + port)
	if err := router.Run(":" + port); err != nil {
		logger.Error("Server failed to start: " + err.Error())
		os.Exit(1)
	}
}
