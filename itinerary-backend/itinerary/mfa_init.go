package itinerary

import (
	"github.com/yourusername/itinerary-backend/itinerary/auth/mfa"
	"github.com/yourusername/itinerary-backend/itinerary/common"
)

// NewTOTPManager creates and returns a new TOTP manager for MFA
func NewTOTPManager(issuer string) *mfa.TOTPManager {
	return mfa.NewTOTPManager(issuer)
}

// InitializeMFA initializes the MFA system
func InitializeMFA(logger *common.Logger) *mfa.TOTPManager {
	logger.Info("Initializing MFA (Multi-Factor Authentication) system")
	totpMgr := mfa.NewTOTPManager("Iternary")
	logger.Info("TOTP manager initialized successfully")
	return totpMgr
}
