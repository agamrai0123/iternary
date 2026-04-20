package itinerary

import (
"github.com/gin-gonic/gin"
"github.com/yourusername/itinerary-backend/itinerary/auth/mfa"
"github.com/yourusername/itinerary-backend/itinerary/auth/oauth"
"github.com/yourusername/itinerary-backend/itinerary/common"
mfahandlers "github.com/yourusername/itinerary-backend/itinerary/handlers/mfa"
oauthhandlers "github.com/yourusername/itinerary-backend/itinerary/handlers/oauth"
)

// SetupRoutes sets up all API and web routes
func SetupRoutes(service *Service, logger *common.Logger, metrics *Metrics, authService *AuthService, totpMgr *mfa.TOTPManager, oauthMgr *oauth.OAuthManager) *gin.Engine {
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

// ==================== Web Routes (HTML Pages) ====================
// Auth pages (no authentication required)
router.GET("/login", handlers.LoginPage)
router.GET("/", handlers.LoginPage) // Redirect home to login

// Protected pages (require authentication)
router.GET("/dashboard", authMiddleware.OptionalAuth(), handlers.DasheAuth(), handlers.GetUserTrip)
router.PUT("/api/user-trips/:id", authMideteUserTrip)
router.GET("/api/userreAuth(), handlers.AddTripSegment)
router.POST("/api/trip-segments/:id/photos", authMiddleware.RequireAuth(), handlers.AddTripPhoto)
router.POST("/api/trip-segments/:id/review", authMiddleware.RequireAuth(), handlers.AddTripReview)
router.POST("/api/user-trips/:id/publish", authMiddleware.RequireAuth(), handlers.PublishUserTrip)

// ==================== Authent.Login)
router.POST("/auth/logout", authHandlers.Logout)
router.GET("/auth/profile", authHandlers.GetProfile)
router.PUT("/auth/profile", authHandlers.UpdateProfile)

// ==================== Group Collaboration Routes (Phase A) =====nt 1) ====================
mfaHandler := mfahandlers.NewHandler((*common.Database)(service.db), logger)
mfahandlers.RegisterMFARoutes(router, mfaHandler)
logger.Info("MFA routes registered")

// ====================logger, oauthMgr)
oauthhandlers.RegisterOAuthRoutes(router, oauthHandler)
logger.Info("OAuth routes registered")

return router
}
