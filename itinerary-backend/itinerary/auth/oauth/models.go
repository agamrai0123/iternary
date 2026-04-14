package oauth

import (
	"time"
)

// Provider represents an OAuth provider configuration
type Provider struct {
	Name             string
	ClientID         string
	ClientSecret     string
	RedirectURL      string
	AuthorizationURL string
	TokenURL         string
	UserInfoURL      string
}

// LinkedAccount represents a user's linked OAuth account
type LinkedAccount struct {
	ID               int       `json:"id"`
	UserID           string    `json:"user_id"`
	Provider         string    `json:"provider"`       // GITHUB, GOOGLE, MICROSOFT
	ProviderUserID   string    `json:"provider_user_id"`
	Email            string    `json:"email"`
	Name             string    `json:"name"`
	AvatarURL        string    `json:"avatar_url"`
	Status           string    `json:"status"`       // ACTIVE, REVOKED
	LinkedAt         time.Time `json:"linked_at"`
	LinkedBy         string    `json:"linked_by"`    // Method: PASSWORD, OAUTH, etc.
	Metadata         string    `json:"metadata"`     // Additional data (e.g., GitHub profile data)
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// OAuthState represents CSRF protection state token
type OAuthState struct {
	StateToken string    `json:"state_token"`
	UserID     string    `json:"user_id"`
	Provider   string    `json:"provider"`
	CreatedAt  time.Time `json:"created_at"`
	ExpiresAt  time.Time `json:"expires_at"`
}

// OAuthUserInfo represents user data from OAuth provider
type OAuthUserInfo struct {
	ID        string `json:"id"`         // Provider-specific user ID
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	// Add additional fields as needed
}

// LinkAccountRequest is the request to link an OAuth account
type LinkAccountRequest struct {
	Code     string `json:"code"`     // Authorization code from OAuth flow
	Provider string `json:"provider"` // GITHUB, GOOGLE, MICROSOFT
	State    string `json:"state"`    // CSRF token from state
}

// LinkAccountResponse is the response after linking an account
type LinkAccountResponse struct {
	Success   bool           `json:"success"`
	Message   string         `json:"message"`
	Account   *LinkedAccount `json:"account,omitempty"`
}

// UnlinkAccountRequest is the request to unlink an OAuth account
type UnlinkAccountRequest struct {
	Provider string `json:"provider"`
}

// UnlinkAccountResponse is the response after unlinking an account
type UnlinkAccountResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// GetLinkedAccountsResponse returns list of linked accounts
type GetLinkedAccountsResponse struct {
	Accounts []*LinkedAccount `json:"accounts"`
	Count    int              `json:"count"`
}

// Provider constants
const (
	ProviderGitHub    = "GITHUB"
	ProviderGoogle    = "GOOGLE"
	ProviderMicrosoft = "MICROSOFT"
)

// Account status constants
const (
	AccountStatusActive  = "ACTIVE"
	AccountStatusRevoked = "REVOKED"
	AccountStatusPending = "PENDING"
)

// Linking method constants
const (
	LinkedByPassword = "PASSWORD"
	LinkedByOAuth    = "OAUTH"
	LinkedBySSO      = "SSO"
)
