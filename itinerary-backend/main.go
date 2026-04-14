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
	// MFA manager
	totpMgr := itinerary.NewTOTPManager("Iternary")
	logger.Info("TOTP manager initialized")

	// OAuth manager
	oauthMgr := itinerary.NewOAuthManager()

	// Register OAuth providers from environment variables
	githubClientID := os.Getenv("GITHUB_OAUTH_CLIENT_ID")
	githubClientSecret := os.Getenv("GITHUB_OAUTH_CLIENT_SECRET")
	if githubClientID != "" && githubClientSecret != "" {
		githubRedirectURL := os.Getenv("OAUTH_REDIRECT_URL")
		if githubRedirectURL == "" {
			githubRedirectURL = "http://localhost:8080/api/v1/oauth/callback/github"
		}
		if err := oauthMgr.RegisterGitHubProvider(githubClientID, githubClientSecret, githubRedirectURL); err != nil {
			logger.Warn("failed to register GitHub OAuth: " + err.Error())
		} else {
			logger.Info("GitHub OAuth provider registered")
		}
	}

	googleClientID := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
	if googleClientID != "" && googleClientSecret != "" {
		googleRedirectURL := os.Getenv("OAUTH_REDIRECT_URL")
		if googleRedirectURL == "" {
			googleRedirectURL = "http://localhost:8080/api/v1/oauth/callback/google"
		}
		if err := oauthMgr.RegisterGoogleProvider(googleClientID, googleClientSecret, googleRedirectURL); err != nil {
			logger.Warn("failed to register Google OAuth: " + err.Error())
		} else {
			logger.Info("Google OAuth provider registered")
		}
	}

	// Initialize router
	router := itinerary.SetupRoutes(svc, logger, metrics, authService, totpMgr, oauthMgr)

	// Run server
	port := config.Server.Port
	logger.Info("Starting server on port " + port)
	if err := router.Run(":" + port); err != nil {
		logger.Error("Server failed to start: " + err.Error())
		os.Exit(1)
	}
}
