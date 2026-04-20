package mfa

import (
	"net/http"

	"github.com/gin-gonic/gin"
	common "github.com/yourusername/itinerary-backend/itinerary/common"
	mfapkg "github.com/yourusername/itinerary-backend/itinerary/auth/mfa"
)

// Handler handles MFA-related HTTP requests
type Handler struct {
	db         *common.Database
	logger     *common.Logger
	totpMgr    *mfapkg.TOTPManager
}

// NewHandler creates a new MFA handler
func NewHandler(db *common.Database, logger *common.Logger) *Handler {
       return &Handler{
	       db:      db,
	       logger:  logger,
	       totpMgr: mfapkg.NewTOTPManager("Iternary"),
       }
}

// StartSetup handles POST /api/v1/mfa/setup/start
// Initiates MFA setup and returns secret + QR code
func (h *Handler) StartSetup(c *gin.Context) {
	userID := c.GetString("user_id") // From auth middleware
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Get user email from database
	// TODO: Query database to get user email
	userEmail := "user@example.com" // Replace with actual database call

	// Check if user already has MFA enabled
	// TODO: Query database to check existing MFA config
	// If exists, return error

	// Generate new secret
	secret, err := h.totpMgr.GenerateSecret(userEmail)
	if err != nil {
		h.logger.Error("failed to generate TOTP secret", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate secret"})
		return
	}

	// Generate QR code
	qrCode, err := h.totpMgr.GetQRCode(userEmail, secret)
	if err != nil {
		h.logger.Error("failed to generate QR code", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate QR code"})
		return
	}

	response := mfapkg.SetupResponse{
		Secret:  secret,
		QRCode:  qrCode,
		Message: "Scan this QR code with your authenticator app",
	}

	c.JSON(http.StatusOK, response)
}

// VerifyAndConfirm handles POST /api/v1/mfa/setup/confirm
// Verifies the initial TOTP code and enables MFA
func (h *Handler) VerifyAndConfirm(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req mfapkg.VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// TODO: Get the secret from request context or database
	// It should be stored temporarily during setup
	secret := "TEMPORARY_SECRET" // Replace with actual secret

	// Verify code
	valid, err := h.totpMgr.VerifyCode(secret, req.Code)
	if err != nil || !valid {
		h.logger.Warn("invalid MFA code", "user_id", userID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid code"})
		return
	}

	// Generate backup codes
	plainCodes, _, err := h.totpMgr.GenerateBackupCodes()
	if err != nil {
		h.logger.Error("failed to generate backup codes", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate backup codes"})
		return
	}

	// Save MFA config to database
	// TODO: SaveMFAConfig(userID, secretHash, hashedCodes)
	h.logger.Info("MFA enabled for user", "user_id", userID)

	response := mfapkg.VerifyResponse{
		Success:     true,
		Message:     "MFA enabled successfully",
		BackupCodes: plainCodes,
	}

	c.JSON(http.StatusOK, response)
}

// VerifyLogin handles POST /api/v1/mfa/verify
// Verifies TOTP code during login
func (h *Handler) VerifyLogin(c *gin.Context) {
	var req mfapkg.VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// TODO: Get userID from session/context (user in MFA challenge)
	userID := c.GetString("mfa_user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid MFA challenge"})
		return
	}

	// TODO: Query database to get MFA config
	secret := "USER_SECRET" // Replace with actual database call

	// Verify code
	valid, err := h.totpMgr.VerifyCode(secret, req.Code)
	if err != nil {
		h.logger.Warn("invalid MFA code", "user_id", userID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid code"})
		return
	}

	if !valid {
		// Check if it's a backup code
		// TODO: VerifyBackupCode and mark as used
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid code"})
		return
	}

	// TODO: Create session for user
	// TODO: Record MFA attempt as successful
	h.logger.Info("MFA verification successful", "user_id", userID)

	response := mfapkg.VerifyResponse{
		Success: true,
		Message: "MFA verification successful",
	}

	c.JSON(http.StatusOK, response)
}

// GetMFAStatus handles GET /api/v1/mfa/status
// Returns whether user has MFA enabled
func (h *Handler) GetMFAStatus(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// TODO: Query database to get MFA status
	status := false // Replace with actual database call

	c.JSON(http.StatusOK, gin.H{
		"mfa_enabled": status,
		"status":      "success",
	})
}

// DisableMFA handles DELETE /api/v1/mfa
// Disables MFA for the user
func (h *Handler) DisableMFA(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// TODO: Verify password
	// TODO: Delete MFA config from database
	h.logger.Info("MFA disabled for user", "user_id", userID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "MFA disabled successfully",
	})
}

// RegenerateBackupCodes handles POST /api/v1/mfa/backup-codes/regenerate
// Generates new backup codes
func (h *Handler) RegenerateBackupCodes(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Generate new backup codes
	plainCodes, _, err := h.totpMgr.GenerateBackupCodes()
	if err != nil {
		h.logger.Error("failed to generate backup codes", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate backup codes"})
		return
	}

	// TODO: Update backup codes in database
	h.logger.Info("backup codes regenerated", "user_id", userID)

	response := mfapkg.BackupCodesResponse{
		BackupCodes: plainCodes,
		Message:     "New backup codes generated. Please save them in a secure location.",
	}

	c.JSON(http.StatusOK, response)
}

// RegisterMFARoutes registers all MFA routes
func RegisterMFARoutes(router *gin.Engine, handler *Handler) {
	v1 := router.Group("/api/v1")
	{
		mfa := v1.Group("/mfa")
		{
			mfa.POST("/setup/start", handler.StartSetup)
			mfa.POST("/setup/confirm", handler.VerifyAndConfirm)
			mfa.POST("/verify", handler.VerifyLogin)
			mfa.GET("/status", handler.GetMFAStatus)
			mfa.DELETE("", handler.DisableMFA)
			mfa.POST("/backup-codes/regenerate", handler.RegenerateBackupCodes)
		}
	}
}
