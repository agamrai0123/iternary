package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware provides authentication middleware for protected routes
type AuthMiddleware struct {
	authService *AuthService
	logger      *Logger
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(authService *AuthService, logger *Logger) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		logger:      logger,
	}
}

// RequireAuth is middleware that requires a valid authentication token
func (am *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header, query parameter, or cookie
		token := c.GetHeader("Authorization")
		if token == "" {
			token = c.Query("token")
		}
		if token == "" {
			// Check for cookie
			token, _ = c.Cookie("token")
		}

		// Remove "Bearer " prefix if present
		if strings.HasPrefix(token, "Bearer ") {
			token = strings.TrimPrefix(token, "Bearer ")
		}

		if token == "" {
			am.logger.Warn("missing_auth_token", "path", c.Request.URL.Path)
			apiErr := NewUnauthorizedError("authentication token is required")
			c.JSON(apiErr.StatusCode, apiErr.ToJSON())
			c.Abort()
			return
		}

		// For MVP: Simple token validation
		// In production, verify JWT signature or check session store
		if len(token) < 20 {
			am.logger.Warn("invalid_token_format", "path", c.Request.URL.Path)
			apiErr := NewUnauthorizedError("invalid authentication token")
			c.JSON(apiErr.StatusCode, apiErr.ToJSON())
			c.Abort()
			return
		}

		// Extract user ID from token (simplified for MVP)
		// In production, decode JWT to get user claims
		userID := extractUserIDFromToken(token)

		// Store token and user ID in context for use in handlers
		c.Set("token", token)
		c.Set("user_id", userID)

		am.logger.Debug("auth_token_valid", "user_id", userID, "path", c.Request.URL.Path)

		c.Next()
	}
}

// OptionalAuth is middleware that optionally checks authentication token
// If token is present, it sets user_id in context. If not, continues without auth.
func (am *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header, query parameter, or cookie
		token := c.GetHeader("Authorization")
		if token == "" {
			token = c.Query("token")
		}
		if token == "" {
			// Check for cookie
			token, _ = c.Cookie("token")
		}

		// Remove "Bearer " prefix if present
		if strings.HasPrefix(token, "Bearer ") {
			token = strings.TrimPrefix(token, "Bearer ")
		}

		if token != "" && len(token) >= 20 {
			// Extract user ID from token (simplified for MVP)
			userID := extractUserIDFromToken(token)
			c.Set("token", token)
			c.Set("user_id", userID)
			am.logger.Debug("optional_auth_token_present", "user_id", userID)
		}

		c.Next()
	}
}

// extractUserIDFromToken extracts user ID from token (simplified for MVP)
// In production, decode JWT to get user claims
func extractUserIDFromToken(token string) string {
	// For MVP, use a simple mapping or decode JWT
	// Since we're using basic auth for MVP, we'll return a demo user ID
	// In production, properly decode the JWT token

	// For now, extract from session storage or use a simplified approach
	// This is a placeholder that returns a demo user ID
	// In production, query the session store with the token

	// For the MVP, we'll use a simplified approach:
	// Token format: "base64-encoded-random-bytes"
	// We'll just return a demo user ID for now

	return "user-001" // Placeholder - should be replaced with proper token validation
}

// GetUserIDFromContext is a helper function to get user ID from context
func GetUserIDFromContext(c *gin.Context) string {
	userID, exists := c.Get("user_id")
	if !exists || userID == nil {
		return ""
	}
	if id, ok := userID.(string); ok {
		return id
	}
	return ""
}

// GetTokenFromContext is a helper function to get token from context
func GetTokenFromContext(c *gin.Context) string {
	token, exists := c.Get("token")
	if !exists || token == nil {
		return ""
	}
	if t, ok := token.(string); ok {
		return t
	}
	return ""
}
