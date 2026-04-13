# Phase 2 Sprint 1: Quick Start Guide

**Goal:** Complete MFA + OAuth foundation in 1-2 weeks  
**Target:** Ready to test with 80%+ feature coverage

---

## ⚡ Quick Start (First 30 Minutes)

### Step 1: Create Feature Branch
```bash
cd d:\Learn\iternary\itinerary-backend
git checkout -b feature/phase2-mfa-oauth
```

### Step 2: Add Required Dependencies
```bash
# Go to project root
cd itinerary-backend

# Add dependencies
go get github.com/pquerna/otp
go get golang.org/x/oauth2
go get golang.org/x/oauth2/github
go get golang.org/x/oauth2/google
go get github.com/skip2/go-qrcode  # for QR code generation

# Update go.mod
go mod tidy

# Verify build still works
go build
```

### Step 3: Create Sprint 1 Project Structure
```bash
# Create new directories
mkdir -p itinerary/auth/mfa
mkdir -p itinerary/auth/oauth
mkdir -p itinerary/validation
mkdir -p itinerary/handlers/oauth
mkdir -p itinerary/handlers/mfa
mkdir -p migrations
mkdir -p tests/integration/auth

# Initialize files
touch itinerary/auth/mfa/totp.go
touch itinerary/auth/mfa/models.go
touch itinerary/auth/oauth/manager.go
touch itinerary/auth/oauth/models.go
touch itinerary/handlers/oauth/handlers.go
touch itinerary/handlers/mfa/handlers.go
touch itinerary/validation/schemas.go
touch migrations/002_add_mfa_tables.sql
touch tests/integration/auth/mfa_test.go
```

---

## 📝 Day 1: MFA Data Models & Database

### File 1: `itinerary/auth/mfa/models.go`

```go
package mfa

import (
	"time"
)

// Config represents MFA configuration for a user
type Config struct {
	ID           string     `json:"id"`
	UserID       string     `json:"user_id"`
	MFAType      string     `json:"mfa_type"` // "totp", "sms", "email"
	Enabled      bool       `json:"enabled"`
	SecretHash   string     `json:"-"` // Encrypted secret
	BackupCodes  string     `json:"-"` // Encrypted JSON array
	CreatedAt    time.Time  `json:"created_at"`
	LastUsedAt   *time.Time `json:"last_used_at"`
	VerifiedAt   *time.Time `json:"verified_at"` // When MFA was confirmed
}

// SetupResponse returned during MFA setup
type SetupResponse struct {
	Secret          string   `json:"secret"`
	QRCode          string   `json:"qr_code"`          // Data URI image
	BackupCodes     []string `json:"backup_codes"`     // 10 one-time codes
	BackupCodesSHA  string   `json:"backup_codes_sha"` // Client verifies hash
}

// VerifyRequest for code verification
type VerifyRequest struct {
	Code string `json:"code" binding:"required,len=6"`
}

// VerifyResponse after successful verification
type VerifyResponse struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	SessionID string `json:"session_id,omitempty"` // New session after MFA verify
}

// BackupCodeUseRecord tracks backup code usage
type BackupCodeUseRecord struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	UsedAt    time.Time `json:"used_at"`
	IPAddress string    `json:"ip_address"`
}
```

### File 2: `migrations/002_add_mfa_tables.sql`

```sql
-- MFA Configuration Table
CREATE TABLE IF NOT EXISTS mfa_configs (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    mfa_type TEXT NOT NULL CHECK(mfa_type IN ('totp', 'sms', 'email')),
    enabled BOOLEAN NOT NULL DEFAULT FALSE,
    secret_hash TEXT NOT NULL, -- Encrypted using AES-256
    backup_codes TEXT NOT NULL, -- Encrypted JSON array of hashes
    created_at TIMESTAMP NOT NULL,
    verified_at TIMESTAMP,
    last_used_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, mfa_type)
);

-- MFA Verification Attempts
CREATE TABLE IF NOT EXISTS mfa_attempts (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    mfa_type TEXT NOT NULL,
    code_type TEXT NOT NULL CHECK(code_type IN ('totp', 'backup')),
    success BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    ip_address TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Backup Code Usage
CREATE TABLE IF NOT EXISTS backup_code_usage (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    used_at TIMESTAMP NOT NULL,
    ip_address TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Indexes for performance
CREATE INDEX idx_mfa_configs_user_id ON mfa_configs(user_id);
CREATE INDEX idx_mfa_attempts_user_id ON mfa_attempts(user_id);
CREATE INDEX idx_mfa_attempts_created ON mfa_attempts(created_at);
CREATE INDEX idx_backup_usage_user_id ON backup_code_usage(user_id);
```

**Action Required:**
1. Save both files
2. Run migration: `sqlite3 itinerary.db < migrations/002_add_mfa_tables.sql`
3. Verify tables created: `sqlite3 itinerary.db ".tables"`

---

## 📝 Day 2: TOTP Implementation

### File 3: `itinerary/auth/mfa/totp.go`

```go
package mfa

import (
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	qrcode "github.com/skip2/go-qrcode"
)

// TOTPManager handles Time-based One-Time Passwords
type TOTPManager struct {
	issuer string // displayed in authenticator app
}

// NewTOTPManager creates a new TOTP manager
func NewTOTPManager(issuer string) *TOTPManager {
	return &TOTPManager{issuer: issuer}
}

// GenerateSecret creates a new TOTP secret
func (tm *TOTPManager) GenerateSecret(userEmail string) (string, error) {
	// Generate 32-byte random secret
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      tm.issuer,
		AccountName: userEmail,
		Period:      30,              // 30 seconds
		SecretSize:  32,              // 256 bits
		Algorithm:   otp.AlgorithmSHA256,
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate secret: %w", err)
	}

	return key.Secret(), nil
}

// GetQRCode returns a PNG QR code as data URI
func (tm *TOTPManager) GetQRCode(userEmail string, secret string) (string, error) {
	// Create provisioning URI
	provisioningURI := fmt.Sprintf(
		"otpauth://totp/%s:%s?secret=%s&issuer=%s&period=30&digits=6&algorithm=SHA256",
		tm.issuer,
		userEmail,
		base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString([]byte(secret)),
		tm.issuer,
	)

	// Generate QR code
	qrc, err := qrcode.Encode(provisioningURI, qrcode.Medium, 256)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Return as data URI
	dataURI := fmt.Sprintf("data:image/png;base64,%s", base32.StdEncoding.EncodeToString(qrc))
	return dataURI, nil
}

// VerifyCode validates a 6-digit TOTP code
// Allows up to 1 time window discrepancy (±30 seconds)
func (tm *TOTPManager) VerifyCode(secret string, code string) (bool, error) {
	if len(code) != 6 {
		return false, fmt.Errorf("invalid code length")
	}

	// Use totp.ValidateCustom for flexibility
	valid, err := totp.ValidateCustom(
		code,
		secret,
		time.Now(),
		totp.ValidateOpts{
			Period:    30,
			Digits:    6,
			Algorithm: otp.AlgorithmSHA256,
		},
	)

	if err != nil {
		return false, fmt.Errorf("TOTP validation error: %w", err)
	}

	return valid, nil
}

// GenerateBackupCodes creates 10 recovery codes
// Returns plain text (to display once) and SHA256 hashes (to store)
func (tm *TOTPManager) GenerateBackupCodes() ([]string, []string, error) {
	plainCodes := make([]string, 10)
	hashedCodes := make([]string, 10)

	for i := 0; i < 10; i++ {
		// Generate 8-character code
		code := fmt.Sprintf("%08x", i*12345+time.Now().UnixNano())[:8]
		plainCodes[i] = code

		// Hash for storage
		hash := sha256.Sum256([]byte(code))
		hashedCodes[i] = fmt.Sprintf("%x", hash)
	}

	return plainCodes, hashedCodes, nil
}

// VerifyBackupCode checks if code matches any hash
// Returns true if valid, and removes it from the list
func (tm *TOTPManager) VerifyBackupCode(hashedCodes []string, code string) (bool, int, error) {
	for i, hashedCode := range hashedCodes {
		hash := sha256.Sum256([]byte(code))
		if fmt.Sprintf("%x", hash) == hashedCode {
			return true, i, nil
		}
	}
	return false, -1, nil
}

// HashSecret for secure storage
func (tm *TOTPManager) HashSecret(secret string) string {
	hash := sha256.Sum256([]byte(secret))
	return fmt.Sprintf("%x", hash)
}
```

**Action Required:**
1. Save the file
2. Build to verify no errors: `go build ./...`

---

## 📝 Day 3: MFA API Handlers

### File 4: `itinerary/handlers/mfa/handlers.go`

```go
package mfahandlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"itinerary/auth/mfa"
	"itinerary/database"
	"itinerary/utils"
)

type MFAHandler struct {
	db    *database.Database
	totp  *mfa.TOTPManager
	logger interface{} // Your logger
}

// NewMFAHandler creates new handler
func NewMFAHandler(db *database.Database, totp *mfa.TOTPManager) *MFAHandler {
	return &MFAHandler{
		db:   db,
		totp: totp,
	}
}

// StartSetup initiates MFA setup (returns secret + QR code)
// POST /api/v1/mfa/setup/start
func (h *MFAHandler) StartSetup(c *gin.Context) {
	// Get user from JWT context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDStr := userID.(string)

	// Generate secret
	secret, err := h.totp.GenerateSecret(userIDStr + "@itinerary.app")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate secret"})
		return
	}

	// Generate QR code
	qrcode, err := h.totp.GetQRCode(userIDStr+"@itinerary.app", secret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
		return
	}

	// Generate backup codes
	plainCodes, hashedCodes, err := h.totp.GenerateBackupCodes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate backup codes"})
		return
	}

	// Store hashes in temporary session (expires in 5 minutes)
	tempKey := "mfa_setup_" + userIDStr
	// Note: In production, store in session/cache, not memory
	// c.Session.Set(tempKey, hashedCodes)

	response := mfa.SetupResponse{
		Secret:      secret,
		QRCode:      qrcode,
		BackupCodes: plainCodes,
		BackupCodesSHA: utils.SHA256Hash(utils.JSON(hashedCodes)),
	}

	c.JSON(http.StatusOK, response)
}

// VerifyAndConfirm verifies MFA code and confirms setup
// POST /api/v1/mfa/setup/confirm
func (h *MFAHandler) VerifyAndConfirm(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req mfa.VerifyRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Verify code
	valid, err := h.totp.VerifyCode("USER_SECRET_FROM_SESSION", req.Code)
	if err != nil || !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid code"})
		return
	}

	// Save MFA config to database
	config := &mfa.Config{
		ID:         utils.GenerateUUID(),
		UserID:     userID.(string),
		MFAType:    "totp",
		Enabled:    true,
		SecretHash: "HASHED_SECRET",
		BackupCodes: "ENCRYPTED_HASHES",
		CreatedAt: time.Now(),
		VerifiedAt: &[]time.Time{time.Now()}[0],
	}

	// Save to database
	if err := h.db.SaveMFAConfig(config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save MFA config"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "MFA enabled successfully",
	})
}

// VerifyLogin verifies MFA code during login
// POST /api/v1/mfa/verify
func (h *MFAHandler) VerifyLogin(c *gin.Context) {
	// This is called after email/password verification succeeds
	// but before session is created

	var req mfa.VerifyRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Get temp login session
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user_id"})
		return
	}

	// Verify code or backup code
	// ... verification logic ...

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// DisableMFA removes MFA
// DELETE /api/v1/mfa
func (h *MFAHandler) DisableMFA(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.db.DeleteMFAConfig(userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to disable MFA"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "MFA disabled"})
}
```

**Action Required:**
1. Save the file
2. Update routing file to register these handlers

---

## 🧪 Day 4-5: Testing & OAuth Setup

### Start OAuth Provider Configuration

**GitHub OAuth Setup (15 minutes):**
1. Go to https://github.com/settings/developers
2. Click "New OAuth App"
3. Fill form:
   - App name: `Itinerary Dev`
   - Homepage URL: `http://localhost:8080`
   - Authorization callback URL: `http://localhost:8080/api/v1/oauth/callback/github`
4. Note: Client ID and Secret
5. Add to `.env.dev`:
   ```env
   GITHUB_OAUTH_CLIENT_ID=your_client_id
   GITHUB_OAUTH_CLIENT_SECRET=your_client_secret
   ```

**Google OAuth Setup (15 minutes):**
1. Go to https://console.cloud.google.com/
2. Create new OAuth 2.0 credentials
3. Configure:
   - Authorized JS origins: `http://localhost:8080`
   - Authorized redirect URIs: `http://localhost:8080/api/v1/oauth/callback/google`
4. Note: Client ID and Secret
5. Add to `.env.dev`

---

## ✅ Day 5: Integration Testing

### Quick Test Checklist

```bash
# 1. Build succeeds
go build ./...

# 2. Run unit tests (when created)
go test ./itinerary/auth/mfa/...

# 3. Manual API test
curl -X POST http://localhost:8080/api/v1/mfa/setup/start \
  -H "Authorization: Bearer YOUR_TOKEN"

# 4. Verify database tables
sqlite3 itinerary.db "SELECT name FROM sqlite_master WHERE type='table' AND name LIKE 'mfa%';"
```

---

## 📊 Progress Tracking

| Day | Task | Status | Est. Hours |
|-----|------|--------|-----------|
| 1 | MFA Models + Database | ⬜ | 2 |
| 2 | TOTP Implementation | ⬜ | 3 |
| 3 | MFA Handlers | ⬜ | 3 |
| 4 | OAuth Setup | ⬜ | 4 |
| 5 | Integration Testing | ⬜ | 4 |

---

## 🎯 Definition of Done

Sprint 1 is complete when:

✅ **MFA Module**
- [ ] Models created and tested
- [ ] Database tables created
- [ ] TOTP generation working
- [ ] TOTP verification working
- [ ] Backup codes functional
- [ ] Setup endpoint working
- [ ] Verify endpoint working

✅ **OAuth Groundwork**
- [ ] GitHub OAuth app configured
- [ ] Google OAuth app configured
- [ ] OAuth manager skeleton created
- [ ] OAuth models defined

✅ **Testing**
- [ ] Unit tests for TOTP (80%+ coverage)
- [ ] Integration test blueprint created
- [ ] Manual testing complete (all scenarios passed)

✅ **Documentation**
- [ ] MFA user guide created
- [ ] OAuth developer guide started
- [ ] API docs updated

---

## 🚀 Next Steps (Sprint 2)

1. **Weeks 2-3:**
   - Complete OAuth callback handlers
   - Implement account linking
   - Add OAuth login flow

2. **Week 4:**
   - API validation framework
   - Request/response documentation
   - Load testing

Your Phase 2 Sprint 1 is now ready to execute!
