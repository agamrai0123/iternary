package itinerary

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthHandlers handles authentication endpoints
type AuthHandlers struct {
	service     *Service
	authService *AuthService
	logger      *Logger
	metrics     *Metrics
}

// NewAuthHandlers creates new auth handlers
func NewAuthHandlers(service *Service, authService *AuthService, logger *Logger, metrics *Metrics) *AuthHandlers {
	return &AuthHandlers{
		service:     service,
		authService: authService,
		logger:      logger,
		metrics:     metrics,
	}
}

// Login handles POST /auth/login
func (ah *AuthHandlers) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ah.metrics.RecordValidationError()
		ah.logger.Warn("invalid_login_payload", "error", err.Error())
		apiErr := NewValidationError("Invalid login format", err.Error())
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	ah.logger.Debug("login_attempt", "email", req.Email)

	// For MVP: Check hardcoded demo users or database
	// In production: Query user database
	user, password := ah.getDemoUser(req.Email)
	if user == nil || !ah.authService.VerifyPassword(req.Password, password) {
		ah.metrics.mu.Lock()
		ah.metrics.AuthorizationErrors++
		ah.metrics.mu.Unlock()

		ah.logger.Warn("login_failed", "email", req.Email, "reason", "invalid_credentials")
		apiErr := NewAPIError(
			ErrUnauthorized,
			"Invalid email or password",
			"Authentication failed",
		)
		apiErr.StatusCode = http.StatusUnauthorized
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	// Create session
	session, err := ah.authService.CreateSession(user.ID, 24*time.Hour)
	if err != nil {
		ah.logger.Error("failed_to_create_session", "error", err.Error(), "email", req.Email)
		apiErr := NewInternalServerError("create_session", err)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	ah.logger.Info("login_successful", "user_id", user.ID, "email", req.Email)

	// Set secure cookie for browser-based requests
	c.SetCookie(
		"token",
		session.Token,
		int(session.ExpiresAt.Sub(time.Now()).Seconds()),
		"/",
		"",
		c.Request.Header.Get("X-Forwarded-Proto") == "https" || c.Request.TLS != nil,
		true,
	)

	c.JSON(http.StatusOK, LoginResponse{
		Token:     session.Token,
		User:      user,
		ExpiresAt: session.ExpiresAt,
	})
}

// Logout handles POST /auth/logout
func (ah *AuthHandlers) Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		apiErr := NewInvalidInputError("Authorization", "token is required")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	ah.logger.Debug("logout_attempt", "token", token[:10]+"...")
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// GetProfile handles GET /auth/profile
func (ah *AuthHandlers) GetProfile(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		apiErr := NewUnauthorizedError("token is required")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	// For MVP: Return demo user
	user, _ := ah.getDemoUser("traveler@example.com")
	if user == nil {
		apiErr := NewNotFoundError("User", "current")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	ah.logger.Debug("profile_retrieved", "user_id", user.ID)
	c.JSON(http.StatusOK, user)
}

// UpdateProfile handles PUT /auth/profile
func (ah *AuthHandlers) UpdateProfile(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		apiErr := NewUnauthorizedError("token is required")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	var req ProfileUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		ah.metrics.RecordValidationError()
		apiErr := NewValidationError("Invalid profile update format", err.Error())
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	// For MVP: Update demo user
	user, _ := ah.getDemoUser("traveler@example.com")
	if user == nil {
		apiErr := NewNotFoundError("User", "current")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	user.FullName = req.FullName
	user.Bio = req.Bio
	user.Avatar = req.Avatar
	user.UpdatedAt = time.Now()

	ah.logger.Info("profile_updated", "user_id", user.ID, "full_name", user.FullName)
	c.JSON(http.StatusOK, user)
}

// getDemoUser returns a demo user for MVP
func (ah *AuthHandlers) getDemoUser(email string) (*AuthUser, string) {
	users := []struct {
		user     *AuthUser
		password string
	}{
		{
			user: &AuthUser{
				ID:        "user-001",
				Username:  "traveler1",
				Email:     "traveler@example.com",
				FullName:  "John Traveler",
				Bio:       "Adventure enthusiast exploring the world",
				Avatar:    "🧳",
				CreatedAt: time.Now().AddDate(0, -1, 0),
				UpdatedAt: time.Now(),
			},
			password: ah.authService.HashPassword("password123"),
		},
		{
			user: &AuthUser{
				ID:        "user-002",
				Username:  "explorer2",
				Email:     "explorer@example.com",
				FullName:  "Jane Explorer",
				Bio:       "Finding hidden gems in every destination",
				Avatar:    "🌍",
				CreatedAt: time.Now().AddDate(0, -2, 0),
				UpdatedAt: time.Now(),
			},
			password: ah.authService.HashPassword("password123"),
		},
	}

	for _, u := range users {
		if u.user.Email == email {
			return u.user, u.password
		}
	}
	return nil, ""
}

// NewUnauthorizedError creates an unauthorized error
func NewUnauthorizedError(details string) *APIError {
	err := NewAPIError(
		ErrUnauthorized,
		"Unauthorized",
		details,
	)
	err.StatusCode = http.StatusUnauthorized
	return err
}
