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

	// ==================== Health & Monitoring Routes (NO AUTH) ====================
	// These are critical for Kubernetes and infrastructure monitoring
	router.GET("/health", handlers.HealthCheckHandler)       // Liveness probe
	router.GET("/ready", handlers.ReadinessHandler)         // Readiness probe
	router.GET("/live", handlers.LivenessHandler)          // Startup probe
	router.GET("/status", handlers.StatusHandler)          // Detailed status
	router.GET("/metrics", handlers.MetricsHandler)        // Prometheus metrics
	
	// Health check aliases for compatibility
	router.GET("/healthz", handlers.HealthCheckHandler)
	router.GET("/readyz", handlers.ReadinessHandler)
	router.GET("/livez", handlers.LivenessHandler)

	// ==================== Web Routes (HTML Pages) ====================
	// Auth pages (no authentication required)
	router.GET("/login", handlers.LoginPage)
	router.GET("/", handlers.LoginPage) // Redirect home to login

	// Protected pages (require authentication)
	// Note: Dashboard temporarily allows unauthenticated access for development
	router.GET("/dashboard", handlers.Dashboard)
	router.GET("/plan-trip", authMiddleware.RequireAuth(), handlers.PlanTripPage)
	router.GET("/my-trips", authMiddleware.RequireAuth(), handlers.MyTripsPage)
	router.GET("/my-trips/:id", authMiddleware.RequireAuth(), handlers.MyTripDetail)
	router.GET("/community", authMiddleware.OptionalAuth(), handlers.CommunityPage)

	// Legacy pages (kept for backward compatibility)
	router.GET("/destination/:id", handlers.DestinationDetail)
	router.GET("/itinerary/:id", handlers.ItineraryDetail)
	router.GET("/create", handlers.CreateItineraryPage)
	router.POST("/create", handlers.CreateItinerarySubmit)
	router.GET("/search", handlers.SearchPage)

	// ==================== API Routes (JSON) ====================
	// Health & Metrics
	router.GET("/api/health", metricsMiddleware.HealthCheckEndpoint())
	router.GET("/api/metrics", metricsMiddleware.MetricsEndpoint())

	// Destination API (no auth required)
	router.GET("/api/destinations", handlers.GetDestinations)

	// Itinerary API (no auth required)
	router.GET("/api/destinations/:destinationId/itineraries", handlers.GetItinerariesByDestination)
	router.GET("/api/itineraries/:itineraryId", handlers.GetItineraryDetail)
	router.POST("/api/itineraries", handlers.CreateItinerary)

	// Like API (no auth required for MVP)
	router.POST("/api/itineraries/:itineraryId/like", handlers.LikeItinerary)

	// Comment API (no auth required for MVP)
	router.POST("/api/itineraries/:itineraryId/comments", handlers.CommentOnItinerary)

	// User Trip API (requires authentication)
	router.POST("/api/user-trips", authMiddleware.RequireAuth(), handlers.CreateUserTrip)
	router.GET("/api/user-trips/:id", authMiddleware.RequireAuth(), handlers.GetUserTrip)
	router.PUT("/api/user-trips/:id", authMiddleware.RequireAuth(), handlers.UpdateUserTrip)
	router.DELETE("/api/user-trips/:id", authMiddleware.RequireAuth(), handlers.DeleteUserTrip)
	router.GET("/api/user-trips", authMiddleware.RequireAuth(), handlers.ListUserTrips)
	router.POST("/api/user-trips/:id/segments", authMiddleware.RequireAuth(), handlers.AddTripSegment)
	router.POST("/api/trip-segments/:id/photos", authMiddleware.RequireAuth(), handlers.AddTripPhoto)
	router.POST("/api/trip-segments/:id/review", authMiddleware.RequireAuth(), handlers.AddTripReview)
	router.POST("/api/user-trips/:id/publish", authMiddleware.RequireAuth(), handlers.PublishUserTrip)

	// ==================== Authentication API Routes ====================
	// Auth endpoints (no authentication required)
	router.POST("/auth/login", authHandlers.Login)
	router.POST("/auth/logout", authHandlers.Logout)
	router.GET("/auth/profile", authHandlers.GetProfile)
	router.PUT("/auth/profile", authHandlers.UpdateProfile)

	// ==================== Group Collaboration Routes (Phase A) ====================
	RegisterGroupRoutes(router, service, authMiddleware, logger)

	return router
}
