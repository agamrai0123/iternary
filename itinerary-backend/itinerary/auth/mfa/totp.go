package mfa

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image/png"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

// TOTPManager handles TOTP operations
type TOTPManager struct {
	// Config for TOTP generation
	issuer string // Typically the app name like "Iternary"
}

// NewTOTPManager creates a new TOTP manager
func NewTOTPManager(issuer string) *TOTPManager {
	return &TOTPManager{
		issuer: issuer,
	}
}

// GenerateSecret creates a new random TOTP secret
// Returns a base32-encoded secret that can be manually entered into authenticator apps
func (tm *TOTPManager) GenerateSecret(userEmail string) (string, error) {
	// Create TOTP key using pquerna/otp
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      tm.issuer,
		AccountName: userEmail,
		Algorithm:   otp.AlgorithmSHA256,
		Period:      TOTPPeriod,
		Digits:      otp.Digits(TOTPDigits),
		SecretSize:  32, // 256 bits
	})

	if err != nil {
		return "", fmt.Errorf("failed to generate TOTP secret: %w", err)
	}

	// Return the Secret (base32-encoded)
	return key.Secret(), nil
}

// GetQRCode generates a QR code for the secret
// Returns a data URI that can be displayed as an image
func (tm *TOTPManager) GetQRCode(userEmail, secret string) (string, error) {
	// Create TOTP Key from secret
	key, err := otp.NewKeyFromURL(fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s&algorithm=%s&digits=%d&period=%d",
		url.QueryEscape(tm.issuer),
		url.QueryEscape(userEmail),
		secret,
		url.QueryEscape(tm.issuer),
		TOTPAlgorithm,
		TOTPDigits,
		TOTPPeriod,
	))

	if err != nil {
		return "", fmt.Errorf("failed to create TOTP key: %w", err)
	}

	// Generate QR code image
	qrImage, err := qrcode.New(key.URL(), qrcode.High)
	if err != nil {
		return "", fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Convert to PNG and create data URI
	img, err := qrImage.Image(300)
	if err != nil {
		return "", fmt.Errorf("failed to create QR image: %w", err)
	}

	// Convert image to base64 data URI
	buf := new(strings.Builder)
	if err := png.Encode(buf, img); err != nil {
		return "", fmt.Errorf("failed to encode image: %w", err)
	}

	// Create data URI
	dataURI := fmt.Sprintf("data:image/png;base64,%s", base64Encode([]byte(buf.String())))
	return dataURI, nil
}

// VerifyCode validates a 6-digit TOTP code
// Uses a time window of ±1 step (±30 seconds) for tolerance
func (tm *TOTPManager) VerifyCode(secret string, code string) (bool, error) {
	// Remove spaces and hyphens from code
	code = strings.ReplaceAll(code, " ", "")
	code = strings.ReplaceAll(code, "-", "")

	// Validate code format (exactly 6 digits)
	if len(code) != TOTPDigits {
		return false, fmt.Errorf("code must be exactly %d digits", TOTPDigits)
	}

	// Validate code is numeric
	if _, err := strconv.Atoi(code); err != nil {
		return false, fmt.Errorf("code must be numeric")
	}

	// Verify TOTP code with tolerance window
	valid, err := totp.ValidateCustom(code, secret, time.Now(), totp.ValidateOpts{
		Period:    TOTPPeriod,
		Skew:      uint(WindowSize),
		Digits:    otp.Digits(TOTPDigits),
		Algorithm: otp.AlgorithmSHA256,
	})

	if err != nil {
		return false, fmt.Errorf("failed to validate TOTP code: %w", err)
	}

	return valid, nil
}

// GenerateBackupCodes creates 10 recovery codes
// Returns plaintext codes and their SHA256 hashes for storage
func (tm *TOTPManager) GenerateBackupCodes() ([]string, []string, error) {
	plainCodes := make([]string, BackupCodeCount)
	hashedCodes := make([]string, BackupCodeCount)

	for i := 0; i < BackupCodeCount; i++ {
		// Generate random 6-byte value and encode as hex (12 character code)
		randomBytes := make([]byte, BackupCodeLength)
		_, err := rand.Read(randomBytes)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate random bytes: %w", err)
		}

		// Create readable backup code (uppercase hex with hyphen for readability)
		code := strings.ToUpper(hex.EncodeToString(randomBytes))
		plainCodes[i] = code

		// Hash the code for storage
		hashedCodes[i] = tm.HashSecret(code)
	}

	return plainCodes, hashedCodes, nil
}

// VerifyBackupCode validates a backup code against stored hashes
func (tm *TOTPManager) VerifyBackupCode(hashedCodes []string, code string) (bool, error) {
	// Remove spaces and hyphens from code for user-friendly input
	code = strings.ToUpper(code)
	code = strings.ReplaceAll(code, " ", "")
	code = strings.ReplaceAll(code, "-", "")

	// Validate code format
	if len(code) != 16 { // 8 bytes = 16 hex chars
		return false, fmt.Errorf("backup code must be 16 characters")
	}

	// Hash the provided code and compare with stored hashes
	providedHash := tm.HashSecret(code)

	for _, storedHash := range hashedCodes {
		if providedHash == storedHash {
			return true, nil
		}
	}

	return false, nil
}

// HashSecret hashes a secret using SHA256
func (tm *TOTPManager) HashSecret(secret string) string {
	hash := sha256.Sum256([]byte(secret))
	return hex.EncodeToString(hash[:])
}

// Helper function to base64 encode
func base64Encode(data []byte) string {
	return base32.StdEncoding.EncodeToString(data)
}

// BackupCodesToJSON converts backup codes to JSON string for storage
func BackupCodesToJSON(codes []string) (string, error) {
	data, err := json.Marshal(codes)
	if err != nil {
		return "", fmt.Errorf("failed to marshal backup codes: %w", err)
	}
	return string(data), nil
}

// BackupCodesFromJSON converts JSON string back to backup codes array
func BackupCodesFromJSON(jsonStr string) ([]string, error) {
	var codes []string
	err := json.Unmarshal([]byte(jsonStr), &codes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal backup codes: %w", err)
	}
	return codes, nil
}
