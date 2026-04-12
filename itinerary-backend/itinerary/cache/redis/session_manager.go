package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// Session represents a user session in cache
type Session struct {
	SessionID string                 `json:"session_id"`
	UserID    string                 `json:"user_id"`
	Username  string                 `json:"username"`
	Email     string                 `json:"email"`
	Data      map[string]interface{} `json:"data"`
	CreatedAt time.Time              `json:"created_at"`
	ExpiresAt time.Time              `json:"expires_at"`
}

// SessionManager handles user session caching
type SessionManager struct {
	client *Client
	prefix string
}

// NewSessionManager creates a new session manager
func NewSessionManager(client *Client) *SessionManager {
	return &SessionManager{
		client: client,
		prefix: "session:",
	}
}

// SaveSession saves a user session to cache
func (sm *SessionManager) SaveSession(ctx context.Context, session *Session, ttl time.Duration) error {
	if session.SessionID == "" {
		return fmt.Errorf("session ID cannot be empty")
	}

	if ttl == 0 {
		ttl = 24 * time.Hour // Default session TTL
	}

	session.CreatedAt = time.Now()
	session.ExpiresAt = time.Now().Add(ttl)

	// Serialize session to JSON
	sessionData, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	key := sm.prefix + session.SessionID

	// Save to Redis
	if err := sm.client.Set(ctx, key, string(sessionData), ttl); err != nil {
		return fmt.Errorf("failed to save session to cache: %w", err)
	}

	return nil
}

// GetSession retrieves a session from cache
func (sm *SessionManager) GetSession(ctx context.Context, sessionID string) (*Session, error) {
	if sessionID == "" {
		return nil, fmt.Errorf("session ID cannot be empty")
	}

	key := sm.prefix + sessionID

	// Get from Redis
	sessionData, err := sm.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	// Deserialize from JSON
	session := &Session{}
	if err := json.Unmarshal([]byte(sessionData), session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	// Check if session has expired
	if session.ExpiresAt.Before(time.Now()) {
		// Remove expired session
		_ = sm.InvalidateSession(ctx, sessionID)
		return nil, ErrCacheMiss
	}

	return session, nil
}

// InvalidateSession removes a session from cache
func (sm *SessionManager) InvalidateSession(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return fmt.Errorf("session ID cannot be empty")
	}

	key := sm.prefix + sessionID
	return sm.client.Delete(ctx, key)
}

// InvalidateUserSessions removes all sessions for a user
func (sm *SessionManager) InvalidateUserSessions(ctx context.Context, userID string) error {
	if userID == "" {
		return fmt.Errorf("user ID cannot be empty")
	}

	// Note: This is a placeholder. In production, you'd need to track
	// which sessions belong to which users (e.g., using sorted sets)
	return nil
}

// RefreshSession updates session expiration time
func (sm *SessionManager) RefreshSession(ctx context.Context, sessionID string, ttl time.Duration) error {
	key := sm.prefix + sessionID

	if err := sm.client.Expire(ctx, key, ttl); err != nil {
		return fmt.Errorf("failed to refresh session: %w", err)
	}

	return nil
}

// SessionExists checks if a session exists
func (sm *SessionManager) SessionExists(ctx context.Context, sessionID string) bool {
	key := sm.prefix + sessionID
	count, err := sm.client.Exists(ctx, key)
	return err == nil && count > 0
}

// GetSessionTTL returns remaining time until session expires
func (sm *SessionManager) GetSessionTTL(ctx context.Context, sessionID string) (time.Duration, error) {
	key := sm.prefix + sessionID
	return sm.client.TTL(ctx, key)
}

// SetSessionData sets additional data in session
func (sm *SessionManager) SetSessionData(ctx context.Context, sessionID, dataKey string, value interface{}) error {
	session, err := sm.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}

	if session.Data == nil {
		session.Data = make(map[string]interface{})
	}

	session.Data[dataKey] = value

	// Get remaining TTL
	ttl, err := sm.GetSessionTTL(ctx, sessionID)
	if err != nil {
		ttl = 24 * time.Hour
	}

	return sm.SaveSession(ctx, session, ttl)
}

// GetSessionData retrieves additional data from session
func (sm *SessionManager) GetSessionData(ctx context.Context, sessionID, dataKey string) (interface{}, error) {
	session, err := sm.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if session.Data == nil {
		return nil, fmt.Errorf("no data found for key: %s", dataKey)
	}

	value, ok := session.Data[dataKey]
	if !ok {
		return nil, fmt.Errorf("key not found in session data: %s", dataKey)
	}

	return value, nil
}
