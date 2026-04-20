package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

// OAuthManager handles OAuth 2.0 operations
type OAuthManager struct {
	providers map[string]*oauth2.Config
}

// NewOAuthManager creates a new OAuth manager
func NewOAuthManager() *OAuthManager {
	return &OAuthManager{
		providers: make(map[string]*oauth2.Config),
	}
}

// RegisterGitHubProvider registers GitHub OAuth configuration
func (om *OAuthManager) RegisterGitHubProvider(clientID, clientSecret, redirectURL string) error {
	if clientID == "" || clientSecret == "" {
		return fmt.Errorf("GitHub OAuth credentials not configured")
	}

	om.providers[ProviderGitHub] = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"user:email",
			"read:user",
		},
		Endpoint: github.Endpoint,
	}

	return nil
}

// RegisterGoogleProvider registers Google OAuth configuration
func (om *OAuthManager) RegisterGoogleProvider(clientID, clientSecret, redirectURL string) error {
	if clientID == "" || clientSecret == "" {
		return fmt.Errorf("Google OAuth credentials not configured")
	}

	om.providers[ProviderGoogle] = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return nil
}

// RegisterMicrosoftProvider registers Microsoft OAuth configuration
func (om *OAuthManager) RegisterMicrosoftProvider(clientID, clientSecret, redirectURL string) error {
	if clientID == "" || clientSecret == "" {
		return fmt.Errorf("Microsoft OAuth credentials not configured")
	}

	microsoftEndpoint := oauth2.Endpoint{
		AuthURL:  "https://login.microsoftonline.com/common/oauth2/v2.0/authorize",
		TokenURL: "https://login.microsoftonline.com/common/oauth2/v2.0/token",
	}

	om.providers[ProviderMicrosoft] = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"user.read",
		},
		Endpoint: microsoftEndpoint,
	}

	return nil
}

// GetAuthURL generates the authorization URL to redirect user to OAuth provider
func (om *OAuthManager) GetAuthURL(provider, state string) (string, error) {
	config, exists := om.providers[provider]
	if !exists {
		return "", fmt.Errorf("provider %s not configured", provider)
	}

	return config.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

// ExchangeCode exchanges the authorization code for an access token
func (om *OAuthManager) ExchangeCode(provider, code string) (*oauth2.Token, error) {
	config, exists := om.providers[provider]
	if !exists {
		return nil, fmt.Errorf("provider %s not configured", provider)
	}

	ctx := context.Background()
	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	return token, nil
}

// GetUserInfo retrieves user information from the OAuth provider
func (om *OAuthManager) GetUserInfo(provider string, token *oauth2.Token) (*OAuthUserInfo, error) {
	switch provider {
	case ProviderGitHub:
		return om.getGitHubUserInfo(token)
	case ProviderGoogle:
		return om.getGoogleUserInfo(token)
	case ProviderMicrosoft:
		return om.getMicrosoftUserInfo(token)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

// getGitHubUserInfo retrieves user info from GitHub API
func (om *OAuthManager) getGitHubUserInfo(token *oauth2.Token) (*OAuthUserInfo, error) {
	// TODO: Implement GitHub API call
	// GET https://api.github.com/user
	// Authorization: token <access_token>
	
	ctx := context.Background()
	client := om.providers[ProviderGitHub].Client(ctx, token)

	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("failed to get GitHub user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API error: %s", string(body))
	}

	// Parse response and return OAuthUserInfo
	// This is a stub - implement actual parsing based on GitHub API response
	return &OAuthUserInfo{
		ID:        "github_user_id",
		Email:     "user@github.com",
		Name:      "GitHub User",
		AvatarURL: "https://avatars.githubusercontent.com/u/0?v=4",
	}, nil
}

// getGoogleUserInfo retrieves user info from Google API
func (om *OAuthManager) getGoogleUserInfo(token *oauth2.Token) (*OAuthUserInfo, error) {
	// TODO: Implement Google API call
	// GET https://www.googleapis.com/oauth2/v2/userinfo
	// Authorization: Bearer <access_token>

	ctx := context.Background()
	client := om.providers[ProviderGoogle].Client(ctx, token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get Google user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Google API error: %s", string(body))
	}

	// Parse response and return OAuthUserInfo
	// This is a stub - implement actual parsing based on Google API response
	return &OAuthUserInfo{
		ID:        "google_user_id",
		Email:     "user@google.com",
		Name:      "Google User",
		AvatarURL: "https://lh3.googleusercontent.com/a/default-user=s96-c",
	}, nil
}

// getMicrosoftUserInfo retrieves user info from Microsoft API
func (om *OAuthManager) getMicrosoftUserInfo(token *oauth2.Token) (*OAuthUserInfo, error) {
	// TODO: Implement Microsoft API call
	// GET https://graph.microsoft.com/v1.0/me
	// Authorization: Bearer <access_token>

	ctx := context.Background()
	client := om.providers[ProviderMicrosoft].Client(ctx, token)

	resp, err := client.Get("https://graph.microsoft.com/v1.0/me")
	if err != nil {
		return nil, fmt.Errorf("failed to get Microsoft user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Microsoft API error: %s", string(body))
	}

	// Parse response and return OAuthUserInfo
	// This is a stub - implement actual parsing based on Microsoft API response
	return &OAuthUserInfo{
		ID:        "microsoft_user_id",
		Email:     "user@microsoft.com",
		Name:      "Microsoft User",
		AvatarURL: "https://graph.microsoft.com/v1.0/me/photo/$value",
	}, nil
}

// ValidateState validates the CSRF token state
func (om *OAuthManager) ValidateState(state string) bool {
	// TODO: Implement state validation against stored states in database
	// Check if state hasn't expired
	// Check if state belongs to current user
	return state != ""
}

// CreateState generates a new CSRF token state
func (om *OAuthManager) CreateState() (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate state: %w", err)
	}

	return base64.URLEncoding.EncodeToString(randomBytes), nil
}

// Helper functions

// IsProviderSupported checks if a provider is registered
func (om *OAuthManager) IsProviderSupported(provider string) bool {
	_, exists := om.providers[provider]
	return exists
}

// GetSupportedProviders returns list of configured providers
func (om *OAuthManager) GetSupportedProviders() []string {
	providers := make([]string, 0, len(om.providers))
	for provider := range om.providers {
		providers = append(providers, provider)
	}
	return providers
}

// SanitizeProvider normalizes provider name to uppercase
func SanitizeProvider(provider string) string {
	return strings.ToUpper(provider)
}
