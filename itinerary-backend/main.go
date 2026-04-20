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
	metrics := itinerary.NewMetrics()

	// Initialize database
	db, err := itinerary.NewDatabase(config, logger)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	logger.Info("Database connection successful")

	// Initialize service
	svc := itinerary.NewService(db, logger)

	// Initialize auth service
	authService := itinerary.NewAuthService(db, logger)

	// Initialize MFA and OAuth components
	totpMgr := itinerary.NewTOTPManager("Iternary")
	logger.Info("TOTP manager initialized")

	oauthMgr := itinerary.NewOAuthManager()

	// Register OAuth providers from environment variables
	githubClientID := os.Getenv("GITHUB_OAUTH_CLIENT_ID")
	githubClientSecret := os.Getenv("GITHUB_OAUTH_CLIENT_SECRET")
	googleClientID := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
	microsoftClientID := os.Getenv("MICROSOFT_OAUTH_CLIENT_ID")
	microsoftClientSecret := os.Getenv("MICROSOFT_OAUTH_CLIENT_SECRET")
	oauthRedirectURL := os.Getenv("OAUTH_REDIRECT_URL")

	if oauthRedirectURL == "" {
		oauthRedirectURL = "http://localhost:8080/api/v1/oauth/callback"
	}

	oauthMgr.RegisterOAuthProviders(
		logger,
		githubClientID, githubClientSecret,
		googleClientID, googleClientSecret,
		microsoftClientID, microsoftClientSecret,
		oauthRedirectURL,
	)

	// Initialize router with MFA and OAuth managers
	router := itinerary.SetupRoutes(svc, logger, metrics, authService, totpMgr, oauthMgr)

	// Run server
	port := config.Server.Port
	logger.Info("Starting server on port " + port)
	if err := router.Run(":" + port); err != nil {
		logger.Error("Server failed to start: " + err.Error())
		os.Exit(1)
	}
}
