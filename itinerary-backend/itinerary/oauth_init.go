package itinerary

import (
	"github.com/yourusername/itinerary-backend/itinerary/auth/oauth"
)

// NewOAuthManager creates and returns a new OAuth manager
func NewOAuthManager() *oauth.OAuthManager {
	return oauth.NewOAuthManager()
}

// InitializeOAuth is disabled - common.Logger type is not defined
// TODO: Refactor to use proper logging and move to oauth package
