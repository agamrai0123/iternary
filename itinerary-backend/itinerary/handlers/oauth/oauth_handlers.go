package oauth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	common "github.com/yourusername/itinerary-backend/itinerary/common"
	oauthpkg "github.com/yourusername/itinerary-backend/itinerary/auth/oauth"
)

// Handler handles OAuth-related HTTP requests
type Handler struct {
	db       *common.Database
	logger   *common.Logger
	oauthMgr *oauthpkg.OAuthManager
}

// NewHandler creates a new OAuth handler
func NewHandler(db *common.Database, logger *common.Logger, oauthMgr *oauthpkg.OAuthManager) *Handler {
       return &Handler{
	       db:       db,
	       logger:   logger,
	       oauthMgr: oauthMgr,
       }
}

// GetAuthURL handles GET /api/v1/oauth/authorize/:provider
// Redirects user to OAuth provider for authorization
func (h *Handler) GetAuthURL(c *gin.Context) {
	provider := c.Param("provider")
	if provider == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider required"})
		return
	}

	// Sanitize provider name
	provider = oauthpkg.SanitizeProvider(provider)

	// Check if provider is supported
	if !h.oauthMgr.IsProviderSupported(provider) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported provider"})
		return
	}

	// Generate state token for CSRF protection
	state, err := h.oauthMgr.CreateState()
	if err != nil {
		h.logger.Error("failed to create state token", "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create state"})
		return
	}

	// TODO: Store state in database with expiration (15 minutes)
	// TODO: Associate state with current user

	// Get authorization URL
	authURL, err := h.oauthMgr.GetAuthURL(provider, state)
	if err != nil {
		h.logger.Error("failed to get auth URL", "provider", provider, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate auth URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"auth_url": authURL,
		"state":    state,
	})
}

// HandleCallback handles GET /api/v1/oauth/callback/:provider
// Processes OAuth callback from provider
func (h *Handler) HandleCallback(c *gin.Context) {
	provider := c.Param("provider")
	if provider == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider required"})
		return
	}

	provider = oauthpkg.SanitizeProvider(provider)

	// Get authorization code and state from query params
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		h.logger.Warn("missing authorization code", "provider", provider)
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing authorization code"})
		return
	}

	// TODO: Validate state token
	if state == "" {
		h.logger.Warn("missing state token", "provider", provider)
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing state token"})
		return
	}

	// Exchange code for token
	token, err := h.oauthMgr.ExchangeCode(provider, code)
	if err != nil {
		h.logger.Error("failed to exchange code", "provider", provider, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to exchange code"})
		return
	}

	// Get user info from provider
	userInfo, err := h.oauthMgr.GetUserInfo(provider, token)
	if err != nil {
		h.logger.Error("failed to get user info", "provider", provider, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user info"})
		return
	}

	// TODO: Check if account already linked
	// TODO: Link to existing user or create new user
	// TODO: Create session

	h.logger.Info("OAuth callback processed", "provider", provider, "user_id", userInfo.ID)

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"message":    "OAuth authentication successful",
		"user_info":  userInfo,
	})
}

// LinkAccount handles POST /api/v1/auth/link-account
// Links an OAuth account to the current user
func (h *Handler) LinkAccount(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req oauthpkg.LinkAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	provider := oauthpkg.SanitizeProvider(req.Provider)

	if !h.oauthMgr.IsProviderSupported(provider) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported provider"})
		return
	}

	// TODO: Validate state token

	// Exchange code for token
	token, err := h.oauthMgr.ExchangeCode(provider, req.Code)
	if err != nil {
		h.logger.Error("failed to exchange code", "provider", provider, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to link account"})
		return
	}

	// Get user info
	_, err = h.oauthMgr.GetUserInfo(provider, token)
	if err != nil {
		h.logger.Error("failed to get user info", "provider", provider, "error", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user info"})
		return
	}

	// TODO: Check if account already linked to another user
	// TODO: Create LinkedAccount in database
	h.logger.Info("account linked", "user_id", userID, "provider", provider)

	response := oauthpkg.LinkAccountResponse{
		Success: true,
		Message: "Account linked successfully",
	}

	c.JSON(http.StatusOK, response)
}

// UnlinkAccount handles DELETE /api/v1/auth/linked-accounts/:provider
// Unlinks an OAuth account from the current user
func (h *Handler) UnlinkAccount(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	provider := c.Param("provider")
	if provider == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider required"})
		return
	}

	provider = oauthpkg.SanitizeProvider(provider)

	// TODO: Delete LinkedAccount from database
	h.logger.Info("account unlinked", "user_id", userID, "provider", provider)

	response := oauthpkg.UnlinkAccountResponse{
		Success: true,
		Message: "Account unlinked successfully",
	}

	c.JSON(http.StatusOK, response)
}

// GetLinkedAccounts handles GET /api/v1/auth/linked-accounts
// Returns list of linked OAuth accounts
func (h *Handler) GetLinkedAccounts(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// TODO: Query database to get linked accounts for user
	var linkedAccounts []*oauthpkg.LinkedAccount
	// linkedAccounts = ... database query ...

	response := oauthpkg.GetLinkedAccountsResponse{
		Accounts: linkedAccounts,
		Count:    len(linkedAccounts),
	}

	c.JSON(http.StatusOK, response)
}

// RegisterOAuthRoutes registers all OAuth routes
func RegisterOAuthRoutes(router *gin.Engine, handler *Handler) {
	v1 := router.Group("/api/v1")
	{
		// OAuth authorization flow
		oauth := v1.Group("/oauth")
		{
			oauth.GET("/authorize/:provider", handler.GetAuthURL)
			oauth.GET("/callback/:provider", handler.HandleCallback)
		}

		// Account linking (auth required)
		auth := v1.Group("/auth")
		{
			auth.POST("/link-account", handler.LinkAccount)
			auth.DELETE("/linked-accounts/:provider", handler.UnlinkAccount)
			auth.GET("/linked-accounts", handler.GetLinkedAccounts)
		}
	}
}
