package security

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// PasswordManager handles password hashing and verification using bcrypt
type PasswordManager struct {
	logger Logger
	cost   int
}

// NewPasswordManager creates a new password manager
// cost: bcrypt cost factor (10-14 recommended, higher = slower but more secure)
func NewPasswordManager(logger Logger, cost int) *PasswordManager {
	if cost < 10 || cost > 14 {
		cost = bcrypt.DefaultCost // 10
	}

	return &PasswordManager{
		logger: logger,
		cost:   cost,
	}
}

// HashPassword hashes a password using bcrypt
func (pm *PasswordManager) HashPassword(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	if len(password) < 6 {
		return "", fmt.Errorf("password must be at least 6 characters")
	}

	if len(password) > 72 { // bcrypt limitation
		return "", fmt.Errorf("password cannot exceed 72 characters")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), pm.cost)
	if err != nil {
		pm.logger.Error("password_hashing_failed", "error", err.Error())
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hash), nil
}

// VerifyPassword verifies a password against its bcrypt hash
func (pm *PasswordManager) VerifyPassword(password, hash string) bool {
	if password == "" || hash == "" {
		return false
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		// Log without revealing whether hash or password was wrong
		pm.logger.Debug("password_verification_failed")
		return false
	}

	return true
}

// ValidatePassword checks if password meets security requirements
func (pm *PasswordManager) ValidatePassword(password string) error {
	if password == "" {
		return fmt.Errorf("password is required")
	}

	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	if len(password) > 72 {
		return fmt.Errorf("password cannot exceed 72 characters")
	}

	// TODO: Add complexity requirements if needed
	// - At least one uppercase letter
	// - At least one lowercase letter
	// - At least one digit
	// - At least one special character

	return nil
}
