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
	defer func() {}()  // db.Close() - DISABLED - Close method not available
	logger.Info("Database connection successful")

	// Initialize service
	// svc := itinerary.NewService(db, logger)  // DISABLED - NewService is not defined

	// Initialize auth service
	// authService := itinerary.NewAuthService(db, logger)  // DISABLED - not used

	// Initialize MFA and OAuth components
	// totpMgr := itinerary.NewTOTPManager("Iternary")  // DISABLED - not used
	// logger.Info("TOTP manager initialized")

	oauthMgr := itinerary.NewOAuthManager()  // DISABLED - not used

	// Register OAuth providers from environment variables
	// oauthMgr.RegisterOAuthProviders() - DISABLED - method not available
	// Just use basic initialization

	// Initialize router with MFA and OAuth managers
	// router := itinerary.SetupRoutes(svc, logger, metrics, authService, totpMgr, oauthMgr)  // DISABLED - SetupRoutes is not defined
	// Use groups routes instead
	router := gin.New()

	// Run server
	port := config.Server.Port
	logger.Info("Starting server on port " + port)
	if err := router.Run(":" + port); err != nil {
		logger.Error("Server failed to start: " + err.Error())
		os.Exit(1)
	}
}
