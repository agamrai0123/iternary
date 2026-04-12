package security

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Session represents a user session
type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	LastUsed  time.Time `json:"last_used"`
}

// SessionStore defines session storage interface
type SessionStore interface {
	SaveSession(ctx context.Context, session *Session) error
	GetSession(ctx context.Context, token string) (*Session, error)
	DeleteSession(ctx context.Context, token string) error
	ValidateSession(ctx context.Context, token string) (string, error) // Returns userID
	RefreshSession(ctx context.Context, token string, newExpiry time.Duration) error
	DeleteUserSessions(ctx context.Context, userID string) error // Logout from all devices
}

// RedisSessionStore implements SessionStore using Redis
type RedisSessionStore struct {
	client redis.Cmdable
	logger Logger
	ttl    time.Duration
	prefix string
}

// NewRedisSessionStore creates a new Redis session store
func NewRedisSessionStore(client redis.Cmdable, logger Logger, ttl time.Duration) *RedisSessionStore {
	if ttl == 0 {
		ttl = 24 * time.Hour // Default 24 hours
	}

	return &RedisSessionStore{
		client: client,
		logger: logger,
		ttl:    ttl,
		prefix: "session:",
	}
}

// SaveSession stores a session in Redis with expiration
func (rss *RedisSessionStore) SaveSession(ctx context.Context, session *Session) error {
	if session == nil || session.Token == "" || session.UserID == "" {
		return fmt.Errorf("invalid session data")
	}

	data, err := json.Marshal(session)
	if err != nil {
		rss.logger.Error("session_marshal_failed", "error", err.Error())
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	key := rss.prefix + session.Token
	
	// Store with TTL
	err = rss.client.Set(ctx, key, data, rss.ttl).Err()
	if err != nil {
		rss.logger.Error("session_save_failed", "error", err.Error())
		return fmt.Errorf("failed to save session: %w", err)
	}

	// Track session for user (for multi-device logout)
	userKey := "user_sessions:" + session.UserID
	err = rss.client.SAdd(ctx, userKey, session.Token).Err()
	if err != nil {
		rss.logger.Warn("session_tracking_failed", "error", err.Error())
		// Don't fail if tracking fails, session is still valid
	}

	rss.logger.Debug("session_saved", "user_id", session.UserID, "expires_in", rss.ttl)
	return nil
}

// GetSession retrieves a session from Redis
func (rss *RedisSessionStore) GetSession(ctx context.Context, token string) (*Session, error) {
	if token == "" {
		return nil, fmt.Errorf("token is empty")
	}

	key := rss.prefix + token
	val, err := rss.client.Get(ctx, key).Result()
	
	if err == redis.Nil {
		tokenHash := rss.hashTokenForLogging(token)
		rss.logger.Warn("session_not_found", "token_hash", tokenHash)
		return nil, fmt.Errorf("session not found or expired")
	}
	
	if err != nil {
		rss.logger.Error("session_retrieval_failed", "error", err.Error())
		return nil, fmt.Errorf("failed to retrieve session: %w", err)
	}

	var session Session
	if err := json.Unmarshal([]byte(val), &session); err != nil {
		rss.logger.Error("session_unmarshal_failed", "error", err.Error())
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	// Check if session has expired
	if time.Now().After(session.ExpiresAt) {
		rss.logger.Warn("session_expired", "user_id", session.UserID)
		// Delete expired session
		rss.client.Del(ctx, key)
		return nil, fmt.Errorf("session expired")
	}

	// Update last used time
	session.LastUsed = time.Now()
	if err := rss.SaveSession(ctx, &session); err != nil {
		// Log but don't fail if update fails
		rss.logger.Warn("session_last_used_update_failed", "error", err.Error())
	}

	return &session, nil
}

// ValidateSession validates a token and returns userID if valid
func (rss *RedisSessionStore) ValidateSession(ctx context.Context, token string) (string, error) {
	session, err := rss.GetSession(ctx, token)
	if err != nil {
		return "", err
	}

	return session.UserID, nil
}

// DeleteSession removes a session from Redis
func (rss *RedisSessionStore) DeleteSession(ctx context.Context, token string) error {
	if token == "" {
		return fmt.Errorf("token is empty")
	}

	key := rss.prefix + token
	
	// First get the session to get userID for cleanup
	session, err := rss.GetSession(ctx, token)
	if err == nil {
		// Remove from user's session set
		userKey := "user_sessions:" + session.UserID
		rss.client.SRem(ctx, userKey, token)
	}

	// Delete the session
	err = rss.client.Del(ctx, key).Err()
	if err != nil {
		rss.logger.Error("session_deletion_failed", "error", err.Error())
		return fmt.Errorf("failed to delete session: %w", err)
	}

	rss.logger.Debug("session_deleted", "token_hash", rss.hashTokenForLogging(token))
	return nil
}

// RefreshSession extends session expiration
func (rss *RedisSessionStore) RefreshSession(ctx context.Context, token string, newExpiry time.Duration) error {
	session, err := rss.GetSession(ctx, token)
	if err != nil {
		return err
	}

	session.ExpiresAt = time.Now().Add(newExpiry)
	return rss.SaveSession(ctx, session)
}

// DeleteUserSessions removes all sessions for a user (logout from all devices)
func (rss *RedisSessionStore) DeleteUserSessions(ctx context.Context, userID string) error {
	if userID == "" {
		return fmt.Errorf("user ID is required")
	}

	userKey := "user_sessions:" + userID
	
	// Get all session tokens for user
	tokens, err := rss.client.SMembers(ctx, userKey).Result()
	if err != nil {
		rss.logger.Error("user_sessions_retrieval_failed", "error", err.Error())
		return fmt.Errorf("failed to retrieve user sessions: %w", err)
	}

	// Delete each session
	for _, token := range tokens {
		key := rss.prefix + token
		rss.client.Del(ctx, key)
	}

	// Delete the user session set
	rss.client.Del(ctx, userKey)

	rss.logger.Info("user_all_sessions_deleted", "user_id", userID, "count", len(tokens))
	return nil
}

// hashTokenForLogging hashes token for secure logging
func (rss *RedisSessionStore) hashTokenForLogging(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:8])
}
