# Phase 2 Sprint 1: MFA & OAuth Foundation

## Sprint Goals
1. Implement Time-based One-Time Password (TOTP) MFA
2. Set up OAuth 2.0 provider integration (GitHub + Google)
3. Create API validation framework
4. Establish testing infrastructure for new features

**Sprint Duration:** 1-2 weeks  
**Target Completion:** Full MFA + OAuth support ready for testing

---

## Task Breakdown

### Task 1: TOTP/MFA Foundation (3 days)

#### 1.1 Create MFA Data Models
**File:** `itinerary/auth/mfa.go`

```go
package auth

// MFAConfig represents MFA configuration for a user
type MFAConfig struct {
    ID           string    `json:"id"`
    UserID       string    `json:"user_id"`
    MFAType      string    `json:"mfa_type"` // "totp", "sms", "email"
    Enabled      bool      `json:"enabled"`
    Secret       string    `json:"-"` // Encrypted
    BackupCodes  []string  `json:"-"` // Encrypted
    CreatedAt    time.Time `json:"created_at"`
    LastUsedAt   *time.Time `json:"last_used_at"`
}

// MFASetupResponse for initial MFA setup
type MFASetupResponse struct {
    Secret      string   `json:"secret"`
    QRCode      string   `json:"qr_code"` // Data URI
    BackupCodes []string `json:"backup_codes"`
}

// MFAVerifyRequest for verification
type MFAVerifyRequest struct {
    Code string `json:"code" binding:"required,len=6"`
}
```

#### 1.2 TOTP Implementation
**File:** `itinerary/security/totp.go`

```go
package security

// TOTPManager handles Time-based OTP generation and verification
type TOTPManager struct {
    secret string
    logger Logger
}

// GenerateSecret creates a new TOTP secret
func (tm *TOTPManager) GenerateSecret() (string, error)

// GetQRCode returns QR code for authenticator apps
func (tm *TOTPManager) GetQRCode(username, issuer string) (string, error)

// VerifyCode validates a TOTP code
func (tm *TOTPManager) VerifyCode(secret, code string) (bool, error)

// GenerateBackupCodes creates recovery codes
func (tm *TOTPManager) GenerateBackupCodes() ([]string, error)

// VerifyBackupCode checks and invalidates a backup code
func (tm *TOTPManager) VerifyBackupCode(hashedCodes []string, code string) (bool, error)
```

#### 1.3 MFA Endpoints
**File:** `itinerary/handlers/mfa_handlers.go`

```
POST   /api/v1/auth/mfa/setup          - Start MFA setup
POST   /api/v1/auth/mfa/verify         - Verify MFA code
POST   /api/v1/auth/mfa/confirm        - Complete MFA setup
DELETE /api/v1/auth/mfa/disable        - Disable MFA
POST   /api/v1/auth/mfa/backup-codes   - Regenerate backup codes
```

#### 1.4 Database Schema
**File:** `migrations/002_add_mfa.sql`

```sql
CREATE TABLE mfa_configs (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL UNIQUE,
    mfa_type TEXT NOT NULL, -- 'totp', 'sms', 'email'
    enabled BOOLEAN DEFAULT FALSE,
    secret TEXT NOT NULL, -- encrypted
    backup_codes TEXT NOT NULL, -- encrypted JSON array
    created_at TIMESTAMP,
    last_used_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE mfa_attempts (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    attempt_type TEXT NOT NULL, -- 'success', 'failure'
    code_type TEXT NOT NULL, -- 'totp', 'backup'
    created_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

---

### Task 2: OAuth 2.0 Integration (3-4 days)

#### 2.1 OAuth Data Models
**File:** `itinerary/auth/oauth.go`

```go
package auth

// OAuthProvider represents external OAuth provider
type OAuthProvider struct {
    ID            string
    Name          string // "github", "google"
    ClientID      string
    ClientSecret  string
    RedirectURL   string
    Scopes        []string
}

// LinkedAccount represents a linked OAuth account
type LinkedAccount struct {
    ID           string    `json:"id"`
    UserID       string    `json:"user_id"`
    Provider     string    `json:"provider"` // "github", "google"
    ProviderID   string    `json:"provider_id"` // ID from provider
    Email        string    `json:"email"`
    Name         string    `json:"name"`
    Avatar       string    `json:"avatar"`
    LinkedAt     time.Time `json:"linked_at"`
}

// OAuthState for secure OAuth flow
type OAuthState struct {
    State       string
    Provider    string
    RedirectURI string
    CreatedAt   time.Time
    ExpiresAt   time.Time
}
```

#### 2.2 OAuth Manager
**File:** `itinerary/auth/oauth_manager.go`

```go
package auth

type OAuthManager struct {
    providers map[string]*OAuthProvider
    logger    Logger
    db        Database
}

// GetAuthURL returns the authorization URL for provider
func (om *OAuthManager) GetAuthURL(provider string, state string) (string, error)

// HandleCallback processes OAuth callback
func (om *OAuthManager) HandleCallback(provider string, code string, state string) (*LinkedAccount, error)

// LinkAccount connects OAuth account to user
func (om *OAuthManager) LinkAccount(userID string, account *LinkedAccount) error

// UnlinkAccount removes OAuth account
func (om *OAuthManager) UnlinkAccount(userID string, provider string) error

// FindUserByOAuthAccount finds user from OAuth account
func (om *OAuthManager) FindUserByOAuthAccount(provider string, providerID string) (*User, error)
```

#### 2.3 OAuth Endpoints
**File:** `itinerary/handlers/oauth_handlers.go`

```
GET    /api/v1/oauth/authorize/:provider   - Start OAuth flow
GET    /api/v1/oauth/callback/:provider    - OAuth callback
POST   /api/v1/auth/link-account           - Link OAuth to existing account
DELETE /api/v1/auth/linked-accounts/:id    - Unlink OAuth account
GET    /api/v1/auth/linked-accounts        - List linked accounts
```

#### 2.4 OAuth Configuration
**File:** `.env.production` additions

```env
# GitHub OAuth
GITHUB_OAUTH_CLIENT_ID=your-github-client-id
GITHUB_OAUTH_CLIENT_SECRET=your-github-client-secret
GITHUB_OAUTH_REDIRECT_URI=https://yourdomain.com/api/v1/oauth/callback/github

# Google OAuth
GOOGLE_OAUTH_CLIENT_ID=your-google-client-id
GOOGLE_OAUTH_CLIENT_SECRET=your-google-client-secret
GOOGLE_OAUTH_REDIRECT_URI=https://yourdomain.com/api/v1/oauth/callback/google
```

---

### Task 3: API Validation Framework (2 days)

#### 3.1 Validation Schemas
**File:** `itinerary/validation/schemas.go`

```go
package validation

type Schema struct {
    Fields map[string]*FieldSchema
}

type FieldSchema struct {
    Required bool
    Type     string      // "string", "email", "number", "uuid"
    MinLen   int
    MaxLen   int
    Pattern  string      // regex
    Enum     []string
}

// Predefined schemas
var (
    UserRegistrationSchema
    LoginSchema
    MFAVerifySchema
    // ... more schemas
)
```

#### 3.2 Validation Middleware
**File:** `itinerary/middleware/request_validation.go`

```go
package middleware

// ValidateRequest validates incoming JSON request
func ValidateRequest(schema *validation.Schema) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Parse and validate request body
        // Return 400 with detailed errors if validation fails
        c.Next()
    }
}
```

#### 3.3 Usage in Handlers
```go
// In route registration
router.POST("/api/v1/auth/register",
    middleware.ValidateRequest(validation.UserRegistrationSchema),
    authHandler.Register,
)
```

---

### Task 4: Testing Infrastructure (2 days)

#### 4.1 Unit Tests
**File:** `itinerary/auth/mfa_test.go`

```go
package auth

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestTOTPGeneration(t *testing.T) {
    // Test TOTP secret generation
}

func TestTOTPVerification(t *testing.T) {
    // Test TOTP code verification
}

func TestBackupCodeGeneration(t *testing.T) {
    // Test backup code generation
}

func TestMFASetup(t *testing.T) {
    // Test complete MFA setup flow
}
```

#### 4.2 Integration Tests
**File:** `itinerary/integration_tests/mfa_test.go`

```go
package integration_tests

import (
    "testing"
)

func TestMFASetupFlow(t *testing.T) {
    // 1. User requests MFA setup
    // 2. System returns secret + QR code
    // 3. User scans QR code
    // 4. User provides TOTP code
    // 5. System confirms MFA enabled
}

func TestOAuthFlow(t *testing.T) {
    // 1. User clicks "Login with GitHub"
    // 2. Redirected to GitHub auth
    // 3. GitHub redirects back with code
    // 4. Mobile app exchanges code for token
    // 5. User logged in
}
```

---

## Implementation Checklist

### Phase 1: MFA Setup
- [ ] Create MFA data models
- [ ] Implement TOTP generation/verification
- [ ] Create MFA database tables
- [ ] Build MFA API endpoints
- [ ] Write unit tests
- [ ] Write integration tests

### Phase 2: OAuth Integration
- [ ] Create OAuth data models
- [ ] Implement OAuth manager
- [ ] Set up GitHub OAuth app
- [ ] Set up Google OAuth app
- [ ] Create OAuth endpoints
- [ ] Handle OAuth callbacks
- [ ] Test OAuth flows

### Phase 3: API Validation
- [ ] Define validation schemas
- [ ] Implement validation middleware
- [ ] Apply validation to endpoints
- [ ] Test validation edge cases

### Phase 4: Testing & Documentation
- [ ] Set up test environment
- [ ] Run full test suite
- [ ] Document MFA setup for users
- [ ] Document OAuth for users
- [ ] Create deployment guide

---

## Development Setup

### 1. Add Dependencies to go.mod
```bash
cd itinerary-backend
go get github.com/pquerna/otp
go get golang.org/x/oauth2
go get golang.org/x/oauth2/github
go get golang.org/x/oauth2/google
go mod tidy
```

### 2. Create New Directories
```bash
mkdir -p itinerary/auth
mkdir -p itinerary/handlers
mkdir -p itinerary/validation
mkdir -p migrations
mkdir -p itinerary/integration_tests
```

### 3. Set Up OAuth Apps

**GitHub:**
1. Go to https://github.com/settings/developers
2. Create new OAuth App
3. Set Authorization callback URL
4. Note Client ID and Secret

**Google:**
1. Go to https://console.cloud.google.com/
2. Create OAuth 2.0 credentials
3. Set Authorized redirect URIs
4. Note Client ID and Secret

### 4. Update Environment
```bash
cp .env.production .env.development
# Add OAuth credentials to .env.development
```

---

## Testing Procedures

### Manual Testing Checklist

**MFA Testing:**
- [ ] User can start MFA setup
- [ ] QR code displays correctly
- [ ] TOTP code verification works
- [ ] Backup codes are generated and work
- [ ] User cannot login without MFA after enabling
- [ ] Backup codes disable themselves after use
- [ ] Old TOTP codes are rejected (time window)

**OAuth Testing:**
- [ ] GitHub login flow completes
- [ ] Google login flow completes
- [ ] Existing user can link OAuth
- [ ] New user can create account via OAuth
- [ ] User data (email, name) imports correctly
- [ ] Can unlink OAuth account
- [ ] Cannot link same OAuth twice

**Validation Testing:**
- [ ] Invalid email rejected
- [ ] Weak passwords rejected
- [ ] Missing required fields caught
- [ ] Oversized payloads rejected
- [ ] Valid requests accepted

---

## Success Criteria

✅ **MFA Complete When:**
- TOTP setup and verification working
- Backup codes functional
- All tests passing (90%+ coverage)
- No security issues in code review

✅ **OAuth Complete When:**
- GitHub login working end-to-end
- Google login working end-to-end
- Account linking working
- No token leaks in logs or errors

✅ **Validation Complete When:**
- All endpoints validated
- Edge cases handled
- Performance acceptable (<50ms overhead)
- Documentation complete

---

## Time Estimates

| Task | Estimate | Actual |
|------|----------|--------|
| MFA Implementation | 16 hours | - |
| OAuth Integration | 12 hours | - |
| API Validation | 8 hours | - |
| Testing | 8 hours | - |
| Documentation | 4 hours | - |
| **Total** | **48 hours** | - |

---

## Resources & References

### TOTP/MFA
- https://github.com/pquerna/otp
- https://tools.ietf.org/html/rfc6238 (TOTP RFC)
- https://tools.ietf.org/html/rfc4648 (Base32 encoding)

### OAuth 2.0
- https://pkg.go.dev/golang.org/x/oauth2
- https://docs.github.com/en/developers/apps/building-oauth-apps
- https://developers.google.com/identity/protocols/oauth2

### Testing
- https://github.com/stretchr/testify
- https://golang.org/pkg/testing/

---

**Sprint 1 Kickoff Ready**
Starting: April 13, 2026
Target Completion: April 25-27, 2026
