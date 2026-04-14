package mfa

import "time"

// Config represents MFA configuration for a user
type Config struct {
	UserID       string    `json:"user_id"`
	SecretHash   string    `json:"secret_hash"`   // Encrypted TOTP secret
	BackupCodes  string    `json:"backup_codes"` // Encrypted JSON array of backup code hashes
	Status       string    `json:"status"`       // "TOTP", "DISABLED"
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// SetupResponse is returned when MFA setup starts
type SetupResponse struct {
	Secret   string `json:"secret"`    // Base32 encoded secret for manual entry
	QRCode   string `json:"qr_code"`   // Data URI for QR code image
	Message  string `json:"message"`
}

// VerifyRequest is used to verify TOTP codes
type VerifyRequest struct {
	Code string `json:"code"` // 6-digit TOTP code
}

// VerifyResponse returns the result of TOTP verification
type VerifyResponse struct {
	Success      bool          `json:"success"`
	Message      string        `json:"message"`
	BackupCodes  []string      `json:"backup_codes,omitempty"` // Only on setup confirmation
}

// BackupCodesResponse returns the recovery codes to user
type BackupCodesResponse struct {
	BackupCodes []string `json:"backup_codes"`
	Message     string   `json:"message"`
}

// BackupCodeUseRecord tracks which backup codes have been used
type BackupCodeUseRecord struct {
	ID        int       `json:"id"`
	UserID    string    `json:"user_id"`
	CodeHash  string    `json:"code_hash"`   // SHA256 hash of the code
	UsedAt    *time.Time `json:"used_at"`    // NULL if unused
	CreatedAt time.Time `json:"created_at"`
}

// BackupCodeVerifyRequest is used to verify backup codes
type BackupCodeVerifyRequest struct {
	Code string `json:"code"` // Backup code to verify
}

// MFA Status constants
const (
	StatusTOTP     = "TOTP"
	StatusDisabled = "DISABLED"
	StatusChallenge = "CHALLENGE"
)

// TOTP Configuration constants
const (
	TOTPPeriod   = 30          // 30 seconds
	TOTPDigits   = 6           // 6-digit codes
	TOTPAlgorithm = "SHA256"   // SHA256 hashing
	WindowSize   = 1           // ±1 time step window
)

// Backup codes configuration
const (
	BackupCodeCount   = 10
	BackupCodeLength  = 8
)
