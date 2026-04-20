package itinerary

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/itinerary-backend/itinerary/common"
)

// ==================== AUTH HANDLERS ====================

// AuthHandlers handles authentication-related requests
type AuthHandlers struct {
	service     *Service
	authService *AuthService
	logger      *common.Logger
	metrics     *Metrics
}

// NewAuthHandlers creates a new auth handlers instance
func NewAuthHandlers(service *Service, authService *AuthService, logger *common.Logger, metrics *Metrics) *AuthHandlers {
	return &AuthHandlers{
		service:     service,
		authService: authService,
		logger:      logger,
		metrics:     metrics,
	}
}

// Login handles POST /auth/login
func (h *AuthHandlers) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid_login_request", "error", err.Error())
		h.metrics.RecordValidationError()
		apiErr := NewInvalidInputError("request_body", "invalid email or password")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	h.logger.Debug("login_attempt", "email", req.Email)

	// Get user by email from database
	user, err := h.service.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		h.logger.Warn("user_not_found", "email", req.Email)
		h.metrics.RecordValidationError()
		apiErr := NewAuthenticationError("invalid email or password")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	// Verify password
	if !h.authService.VerifyPassword(user.PasswordHash, req.Password) {
		h.logger.Warn("invalid_password", "email", req.Email)
		h.metrics.RecordValidationError()
		apiErr := NewAuthenticationError("invalid email or password")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	// Create session
	session, err := h.authService.CreateSession(user.ID, 24*time.Hour)
	if err != nil {
		h.logger.Error("failed_to_create_session", "error", err.Error())
		apiErr := NewInternalServerError("session_creation", err)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	// Store session in database
	if err := h.service.CreateSession(session); err != nil {
		h.logger.Error("failed_to_store_session", "error", err.Error())
		apiErr := NewInternalServerError("session_storage", err)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	h.logger.Info("user_logged_in", "user_id", user.ID)

	resp := LoginResponse{
		Token: session.Token,
		User: &AuthUser{
			ID:       user.ID,
			Email:    user.Email,
			Username: user.Username,
			FullName: user.FullName,
			Avatar:   user.Avatar,
		},
		ExpiresAt: session.ExpiresAt,
	}

	c.JSON(http.StatusOK, resp)
}

// Register handles POST /auth/register
func (h *AuthHandlers) Register(c *gin.Context) {
	var req SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid_signup_request", "error", err.Error())
		h.metrics.RecordValidationError()
		apiErr := NewInvalidInputError("request_body", "all fields are required")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	h.logger.Debug("signup_attempt", "email", req.Email)

	// Check if user exists
	existingUser, _ := h.service.GetUserByEmail(req.Email)
	if existingUser != nil {
		h.logger.Warn("user_already_exists", "email", req.Email)
		h.metrics.RecordValidationError()
		apiErr := NewInvalidInputError("email", "email already registered")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	// Create new user
	user := &AuthUser{
		ID:       uuid.New().String(),
		Email:    req.Email,
		Username: req.Username,
		FullName: req.FullName,
	}

	// Hash password
	passwordHash := h.authService.HashPassword(req.Password)

	// Store user in database
	if err := h.service.CreateUser(user, passwordHash); err != nil {
		h.logger.Error("failed_to_create_user", "error", err.Error())
		apiErr := NewInternalServerError("user_creation", err)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	// Create session
	session, err := h.authService.CreateSession(user.ID, 24*time.Hour)
	if err != nil {
		h.logger.Error("failed_to_create_session", "error", err.Error())
		apiErr := NewInternalServerError("session_creation", err)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	// Store session
	if err := h.service.CreateSession(session); err != nil {
		h.logger.Error("failed_to_store_session", "error", err.Error())
		apiErr := NewInternalServerError("session_storage", err)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	h.logger.Info("user_registered", "user_id", user.ID, "email", user.Email)
	h.logger.Debug("login_success", "user_id", user.ID)

	resp := LoginResponse{
		Token: session.Token,
		User:  user,
		ExpiresAt: session.ExpiresAt,
	}

	c.JSON(http.StatusCreated, resp)
}

// GetProfile handles GET /auth/profile
func (h *AuthHandlers) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		apiErr := NewAuthenticationError("user not authenticated")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	h.logger.Debug("get_profile", "user_id", userID)

	user, err := h.service.GetUserByID(userID.(string))
	if err != nil || user == nil {
		h.logger.Warn("user_not_found", "user_id", userID)
		apiErr := NewAuthenticationError("user not found")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateProfile handles PUT /auth/profile
func (h *AuthHandlers) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		apiErr := NewAuthenticationError("user not authenticated")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	var req ProfileUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid_profile_update_request", "error", err.Error())
		apiErr := NewInvalidInputError("request_body", "invalid request")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	h.logger.Debug("update_profile", "user_id", userID)

	if err := h.service.UpdateProfile(userID.(string), &req); err != nil {
		h.logger.Error("failed_to_update_profile", "error", err.Error())
		apiErr := NewInternalServerError("profile_update", err)
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	user, _ := h.service.GetUserByID(userID.(string))
	c.JSON(http.StatusOK, user)
}

// Logout handles POST /auth/logout
func (h *AuthHandlers) Logout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		apiErr := NewAuthenticationError("user not authenticated")
		c.JSON(apiErr.StatusCode, apiErr.ToJSON())
		return
	}

	h.logger.Debug("logout", "user_id", userID)

	// Invalidate session
	if err := h.service.InvalidateSession(c.GetHeader("Authorization")); err != nil {
		h.logger.Warn("failed_to_invalidate_session", "error", err.Error())
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}
