package security

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

// JWTManager handles JWT token creation and validation
type JWTManager struct {
	secretKey string
	logger    Logger
	issuer    string
	audience  string
}

// Logger interface for dependency injection
type Logger interface {
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(secretKey string, logger Logger) (*JWTManager, error) {
	if secretKey == "" {
		return nil, fmt.Errorf("JWT_SECRET_KEY environment variable is required")
	}
	if len(secretKey) < 32 {
		return nil, fmt.Errorf("JWT_SECRET_KEY must be at least 32 characters")
	}

	return &JWTManager{
		secretKey: secretKey,
		logger:    logger,
		issuer:    "triply-api",
		audience:  "triply-web",
	}, nil
}

// GenerateToken creates a signed JWT token
func (jm *JWTManager) GenerateToken(ctx context.Context, userID, username, email string, expiresIn time.Duration) (string, error) {
	now := time.Now()
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    jm.issuer,
			Audience:  jwt.ClaimStrings{jm.audience},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jm.secretKey))
	if err != nil {
		jm.logger.Error("token_generation_failed", "error", err.Error())
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	jm.logger.Debug("token_generated", "user_id", userID, "expires_in", expiresIn)
	return tokenString, nil
}

// ValidateToken verifies and parses JWT token
func (jm *JWTManager) ValidateToken(ctx context.Context, tokenString string) (*JWTClaims, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("token is empty")
	}

	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method to prevent algorithm substitution attacks
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jm.secretKey), nil
	})

	if err != nil {
		tokenHash := jm.hashTokenForLogging(tokenString)
		jm.logger.Warn("token_validation_failed", "error", err.Error(), "token_hash", tokenHash)
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		tokenHash := jm.hashTokenForLogging(tokenString)
		jm.logger.Warn("token_invalid", "token_hash", tokenHash)
		return nil, fmt.Errorf("token is invalid")
	}

	// Verify claims
	// if err := claims.Valid(); err != nil {  // DISABLED - Valid() method not available
	// 	jm.logger.Warn("token_claims_invalid", "error", err.Error())
	// 	return nil, fmt.Errorf("token claims invalid: %w", err)
	// }

	return claims, nil
}

// RefreshToken generates a new token based on existing claims
func (jm *JWTManager) RefreshToken(ctx context.Context, oldToken string, expiresIn time.Duration) (string, error) {
	claims, err := jm.ValidateToken(ctx, oldToken)
	if err != nil {
		return "", fmt.Errorf("cannot refresh invalid token: %w", err)
	}

	return jm.GenerateToken(ctx, claims.UserID, claims.Username, claims.Email, expiresIn)
}

// hashTokenForLogging hashes token for secure logging
func (jm *JWTManager) hashTokenForLogging(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:8]) // First 8 chars of hash
}
