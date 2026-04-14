# Phase 2 Sprint 1 - Integration Guide

**Purpose:** How to integrate the Sprint 1 components into your main application  
**Status:** All modules created and ready for integration

---

## 📋 Integration Checklist

### Step 1: Initialize MFA Components in main.go

```go
package main

import (
    "itinerary/auth/mfa"
    mfahandlers "itinerary/handlers/mfa"
)

func init() {
    // Initialize TOTP Manager
    totpManager := mfa.NewTOTPManager("Itinerary")
    
    // Initialize MFA Handler
    mfaHandler := mfahandlers.NewHandler(totpManager)
    
    // Register MFA endpoints
    // (see route registration below)
}
```

### Step 2: Initialize OAuth Components in main.go

```go
import (
    "itinerary/auth/oauth"
    oauthhandlers "itinerary/handlers/oauth"
    "os"
)

func init() {
    // Initialize OAuth Manager
    oauthManager := oauth.NewManager()
    
    // Register GitHub OAuth
    oauthManager.RegisterGitHubProvider(
        os.Getenv("GITHUB_OAUTH_CLIENT_ID"),
        os.Getenv("GITHUB_OAUTH_CLIENT_SECRET"),
        os.Getenv("GITHUB_OAUTH_REDIRECT_URI"),
    )
    
    // Register Google OAuth
    oauthManager.RegisterGoogleProvider(
        os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
        os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
        os.Getenv("GOOGLE_OAUTH_REDIRECT_URI"),
    )
    
    // Initialize OAuth Handler
    oauthHandler := oauthhandlers.NewHandler(oauthManager)
    
    // Register OAuth endpoints
    // (see route registration below)
}
```

### Step 3: Register MFA Routes

In your routes file (e.g., `routes/routes.go`):

```go
import (
    mfahandlers "itinerary/handlers/mfa"
)

func SetupMFARoutes(router *gin.Engine, mfaHandler *mfahandlers.Handler) {
    mfa := router.Group("/api/v1/mfa")
    {
        // MFA Setup
        mfa.POST("/setup/start", mfaHandler.StartSetup)
        mfa.POST("/setup/confirm", mfaHandler.VerifyAndConfirm)
        
        // MFA Login
        mfa.POST("/verify", mfaHandler.VerifyLogin)
        
        // MFA Management
        mfa.DELETE("", mfaHandler.DisableMFA)
        mfa.GET("/status", mfaHandler.GetMFAStatus)
        mfa.POST("/backup-codes/regenerate", mfaHandler.RegenerateBackupCodes)
    }
}
```

### Step 4: Register OAuth Routes

```go
import (
    oauthhandlers "itinerary/handlers/oauth"
)

func SetupOAuthRoutes(router *gin.Engine, oauthHandler *oauthhandlers.Handler) {
    oauth := router.Group("/api/v1/oauth")
    {
        // OAuth Authorization
        oauth.GET("/authorize/:provider", oauthHandler.GetAuthURL)
        oauth.GET("/callback/:provider", oauthHandler.HandleCallback)
    }
    
    auth := router.Group("/api/v1/auth")
    {
        // Account Linking
        auth.POST("/link-account", oauthHandler.LinkAccount)
        auth.DELETE("/linked-accounts/:provider", oauthHandler.UnlinkAccount)
        auth.GET("/linked-accounts", oauthHandler.GetLinkedAccounts)
    }
}
```

### Step 5: Set Up Database Migration

```bash
# Run the migration
sqlite3 itinerary.db < migrations/002_add_mfa_oauth.sql

# Verify tables
sqlite3 itinerary.db ".tables"
# Output should include: mfa_configs, mfa_attempts, backup_code_usage, linked_accounts, oauth_states
```

### Step 6: Add Environment Variables

Create/update `.env.production`:

```env
# OAuth - GitHub
GITHUB_OAUTH_CLIENT_ID=your_github_client_id
GITHUB_OAUTH_CLIENT_SECRET=your_github_client_secret
GITHUB_OAUTH_REDIRECT_URI=https://yourdomain.com/api/v1/oauth/callback/github

# OAuth - Google
GOOGLE_OAUTH_CLIENT_ID=your_google_client_id
GOOGLE_OAUTH_CLIENT_SECRET=your_google_client_secret
GOOGLE_OAUTH_REDIRECT_URI=https://yourdomain.com/api/v1/oauth/callback/google

# MFA
MFA_ISSUER=Itinerary
```

---

## 🔌 Database Integration TODOs

### In `itinerary/database/database.go`, add these methods:

```go
// MFA Configuration Methods
func (db *Database) SaveMFAConfig(config *mfa.Config) error {
    query := `
        INSERT OR REPLACE INTO mfa_configs 
        (id, user_id, mfa_type, enabled, secret_hash, backup_codes, created_at, verified_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `
    _, err := db.Exec(query,
        config.ID,
        config.UserID,
        config.MFAType,
        config.Enabled,
        config.SecretHash,
        config.BackupCodes,
        config.CreatedAt,
        config.VerifiedAt,
    )
    return err
}

func (db *Database) GetMFAConfig(userID string) (*mfa.Config, error) {
    query := `SELECT id, user_id, mfa_type, enabled, secret_hash, 
              backup_codes, created_at, verified_at, last_used_at 
              FROM mfa_configs WHERE user_id = ?`
    
    config := &mfa.Config{}
    err := db.QueryRow(query, userID).Scan(
        &config.ID,
        &config.UserID,
        &config.MFAType,
        &config.Enabled,
        &config.SecretHash,
        &config.BackupCodes,
        &config.CreatedAt,
        &config.VerifiedAt,
        &config.LastUsedAt,
    )
    return config, err
}

func (db *Database) DeleteMFAConfig(userID string) error {
    _, err := db.Exec("DELETE FROM mfa_configs WHERE user_id = ?", userID)
    return err
}

// OAuth LinkedAccount Methods
func (db *Database) CreateLinkedAccount(account *oauth.LinkedAccount) error {
    query := `
        INSERT INTO linked_accounts 
        (id, user_id, provider, provider_id, email, name, avatar, linked_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `
    _, err := db.Exec(query,
        account.ID,
        account.UserID,
        account.Provider,
        account.ProviderID,
        account.Email,
        account.Name,
        account.Avatar,
        account.LinkedAt,
    )
    return err
}

func (db *Database) GetLinkedAccounts(userID string) ([]oauth.LinkedAccount, error) {
    query := `SELECT id, user_id, provider, provider_id, email, name, avatar, linked_at 
              FROM linked_accounts WHERE user_id = ?`
    
    rows, err := db.Query(query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    accounts := []oauth.LinkedAccount{}
    for rows.Next() {
        account := oauth.LinkedAccount{}
        err := rows.Scan(
            &account.ID,
            &account.UserID,
            &account.Provider,
            &account.ProviderID,
            &account.Email,
            &account.Name,
            &account.Avatar,
            &account.LinkedAt,
        )
        if err != nil {
            return nil, err
        }
        accounts = append(accounts, account)
    }
    return accounts, nil
}

func (db *Database) DeleteLinkedAccount(userID, provider string) error {
    _, err := db.Exec(
        "DELETE FROM linked_accounts WHERE user_id = ? AND provider = ?",
        userID,
        provider,
    )
    return err
}

// OAuth State Methods
func (db *Database) SaveOAuthState(state *oauth.OAuthState) error {
    query := `
        INSERT INTO oauth_states 
        (state, provider, user_id, redirect_uri, created_at, expires_at)
        VALUES (?, ?, ?, ?, ?, ?)
    `
    _, err := db.Exec(query,
        state.State,
        state.Provider,
        state.UserID,
        state.RedirectURI,
        state.CreatedAt,
        state.ExpiresAt,
    )
    return err
}

func (db *Database) GetOAuthState(stateToken string) (*oauth.OAuthState, error) {
    query := `SELECT state, provider, user_id, redirect_uri, created_at, expires_at 
              FROM oauth_states WHERE state = ?`
    
    state := &oauth.OAuthState{}
    err := db.QueryRow(query, stateToken).Scan(
        &state.State,
        &state.Provider,
        &state.UserID,
        &state.RedirectURI,
        &state.CreatedAt,
        &state.ExpiresAt,
    )
    return state, err
}

func (db *Database) DeleteExpiredOAuthStates() error {
    _, err := db.Exec("DELETE FROM oauth_states WHERE expires_at < datetime('now')")
    return err
}
```

---

## 🧪 Testing the Integration

### Test 1: MFA Setup Flow

```bash
# Start server
go run main.go

# Get JWT token first (from login endpoint)
TOKEN="your_jwt_token_here"

# 1. Initiate MFA setup
curl -X POST http://localhost:8080/api/v1/mfa/setup/start \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json"

# Response will include: secret, qr_code, backup_codes

# 2. Scan QR code with authenticator app (Google Authenticator, Authy, etc.)

# 3. Confirm MFA with code from app
curl -X POST http://localhost:8080/api/v1/mfa/setup/confirm \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"code": "123456"}'
```

### Test 2: OAuth Authorization

```bash
# 1. Get GitHub authorization URL
curl "http://localhost:8080/api/v1/oauth/authorize/github"

# Response: {"auth_url": "https://github.com/login/oauth/authorize?..."}

# 2. User clicks link and authorizes
# 3. GitHub redirects to callback with code
# 4. Your app exchanges code for token
curl "http://localhost:8080/api/v1/oauth/callback/github?code=CODE&state=STATE"
```

### Test 3: Request Validation

```bash
# Test invalid email
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username": "test", "email": "invalid", "password": "short"}'

# Expected: 400 Bad Request with validation errors
```

---

## 🔐 Security Considerations

### 1. Token Storage (for OAuth)
- Store refresh tokens in secure cookies (httpOnly, Secure flags)
- Never store access tokens in localStorage
- Implement token rotation

### 2. MFA Secret Storage
- Encrypt secrets at rest (AES-256)
- Never log plaintext secrets
- Use the `HashSecret()` method for storage

### 3. CSRF Protection
- Validate state tokens on OAuth callback
- Use CSRF tokens for state-changing operations

### 4. Rate Limiting
- Apply rate limiting to MFA verify endpoint
- Lock account after N failed attempts
- Log all MFA attempts for audit

---

## 🚀 Next Steps After Integration

1. **Complete OAuth User Info Retrieval**
   - Implement `getGitHubUserInfo()` in oauth/manager.go
   - Implement `getGoogleUserInfo()` in oauth/manager.go

2. **Add Session Management**
   - Create new session after MFA verification
   - Invalidate old sessions when MFA is enabled

3. **Implement Request Validation**
   - Add validation middleware to routes
   - Validate all incoming requests

4. **Add Comprehensive Tests**
   - Unit tests for TOTP
   - Unit tests for validators
   - Integration tests for flows

5. **Deploy to Staging**
   - Run on staging environment
   - Perform security audit
   - Load testing

---

## 📞 Troubleshooting

### Issue: "module not found" errors
**Solution:** 
```bash
go mod tidy
go mod download
```

### Issue: Database migration fails
**Solution:**
```bash
sqlite3 itinerary.db ".schema"  # Check existing tables
sqlite3 itinerary.db < migrations/002_add_mfa_oauth.sql  # Re-run migration
```

### Issue: OAuth callback returns "unknown provider"
**Solution:**
- Verify provider name matches (github, google, microsoft)
- Check environment variables are loaded
- Verify manager is initialized

### Issue: TOTP codes not validating
**Solution:**
- Check system clock is synchronized
- TOTP allows ±30 second window
- Verify QR code was scanned correctly

---

## 📚 Code Examples

### Using MFA in Your Handlers

```go
import (
    "itinerary/auth/mfa"
)

func LoginWithMFA(c *gin.Context) {
    // 1. Verify email and password (existing code)
    
    // 2. Check if user has MFA enabled
    // mfaConfig, _ := db.GetMFAConfig(userID)
    
    // 3. If MFA enabled, require verification
    // if mfaConfig != nil && mfaConfig.Enabled {
    //     // Store temp session, ask for MFA code
    //     c.JSON(http.StatusAccepted, gin.H{
    //         "message": "MFA verification required",
    //         "user_id": userID,
    //     })
    //     return
    // }
    
    // 4. Create final session
}
```

### Using Validation in Your Handlers

```go
import (
    "itinerary/validation"
)

func RegisterUser(c *gin.Context) {
    var data map[string]interface{}
    c.BindJSON(&data)
    
    // Validate against schema
    validator := validation.NewValidator()
    result := validator.ValidateObject(data, validation.UserRegistrationSchema)
    
    if !result.Valid {
        c.JSON(http.StatusBadRequest, gin.H{
            "errors": result.Errors,
        })
        return
    }
    
    // Process registration
}
```

---

**Integration Guide Complete**

All Sprint 1 components are ready to integrate. Follow this guide to add MFA and OAuth to your Itinerary backend!

