package itinerary

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up all API and web routes
func SetupRoutes(service *Service, logger *Logger, metrics *Metrics, authService *AuthService) *gin.Engine {
	router := gin.New()

	// Create metrics middleware
	metricsMiddleware := NewMetricsMiddleware(metrics, logger)

	// Create auth middleware
	authMiddleware := NewAuthMiddleware(authService, logger)

	// Apply middleware stack with proper order
	router.Use(metricsMiddleware.PanicRecoveryMiddleware()) // Panic recovery first
	router.Use(logger.RequestLogger())                      // Request logging
	router.Use(metricsMiddleware.MetricsHandler())          // Metrics collection
	router.Use(metricsMiddleware.ErrorHandlerMiddleware())  // Error handling
	router.Use(logger.ErrorLogger())                        // Error logging

	// Set template functions before loading templates
	router.SetFuncMap(TemplateFuncs())

	// Load templates
	router.LoadHTMLGlob("templates/*.html")

	// Serve static files (CSS, JS, images)
	router.Static("/static", "./static")

	// Initialize handlers
	handlers := NewHandlers(service, logger, metrics)
	authHandlers := NewAuthHandlers(service, authService, logger, metrics)

	// ==================== API Routes (JSON) ====================
	// Health & Metrics
	router.GET("/api/health", metricsMiddleware.HealthCheckEndpoint())
	router.GET("/api/metrics", metricsMiddleware.MetricsEndpoint())

	// ==================== Authentication API Routes ====================
	// Auth endpoints (no authentication required)
	router.POST("/auth/login", authHandlers.Login)
	router.POST("/auth/register", authHandlers.Register)
	router.POST("/auth/logout", authMiddleware.RequireAuth(), authHandlers.Logout)
	router.GET("/auth/profile", authMiddleware.RequireAuth(), authHandlers.GetProfile)
	router.PUT("/auth/profile", authMiddleware.RequireAuth(), authHandlers.UpdateProfile)

	// ==================== Cities API (no auth required) ====================
	router.GET("/api/cities", handlers.GetCities)
	router.GET("/api/cities/:cityId", handlers.GetCityByID)

	// ==================== Trip Posts API (no auth required for viewing) ====================
	router.GET("/api/trip-posts", handlers.GetTripPosts)
	router.GET("/api/trip-posts/:postId", handlers.GetTripPostByID)
	router.GET("/api/cities/:cityId/trip-posts", handlers.GetCityTripPosts)
	router.POST("/api/trip-posts/:postId/like", handlers.LikeTripPost)
	router.POST("/api/trip-posts/:postId/save", authMiddleware.RequireAuth(), handlers.SaveTripPost)
	router.POST("/api/trip-posts/:postId/add-to-itinerary", authMiddleware.RequireAuth(), handlers.AddTripPostToItinerary)

	// ==================== User Trips API (requires authentication) ====================
	router.POST("/api/user-trips", authMiddleware.RequireAuth(), handlers.CreateUserTrip)
	router.GET("/api/user-trips", authMiddleware.RequireAuth(), handlers.GetUserTrips)
	router.GET("/api/user-trips/:tripId", authMiddleware.RequireAuth(), handlers.GetUserTripByID)
	router.PUT("/api/user-trips/:tripId", authMiddleware.RequireAuth(), handlers.UpdateUserTrip)
	router.DELETE("/api/user-trips/:tripId", authMiddleware.RequireAuth(), handlers.DeleteUserTrip)

	// ==================== Trip Segments API (requires authentication) ====================
	router.POST("/api/user-trips/:tripId/segments", authMiddleware.RequireAuth(), handlers.AddTripSegment)
	router.PUT("/api/user-trips/:tripId/segments/:segmentId", authMiddleware.RequireAuth(), handlers.UpdateTripSegment)
	router.DELETE("/api/user-trips/:tripId/segments/:segmentId", authMiddleware.RequireAuth(), handlers.DeleteTripSegment)
	router.POST("/api/user-trips/:tripId/segments/:segmentId/photos", authMiddleware.RequireAuth(), handlers.AddTripPhoto)
	router.POST("/api/user-trips/:tripId/segments/:segmentId/review", authMiddleware.RequireAuth(), handlers.AddTripReview)

	// ==================== Group Collaboration Routes (Phase A) ====================
	RegisterGroupRoutes(router, service, authMiddleware, logger)

	return router
}
