package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// SessionCache manages user sessions with Redis backend
type SessionCache struct {
	client *Client
	prefix string
}

// NewSessionCache creates a new session cache
func NewSessionCache(client *Client) *SessionCache {
	return &SessionCache{
		client: client,
		prefix: "session:",
	}
}

// Create creates a new session
func (sc *SessionCache) Create(ctx context.Context, sessionID, userID string, duration time.Duration) (*Session, error) {
	if duration == 0 {
		duration = 24 * time.Hour // Default session duration
	}

	session := &Session{
		ID:        sessionID,
		UserID:    userID,
		Data:      make(map[string]interface{}),
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(duration),
	}

	key := sc.prefix + sessionID
	sessionData, err := json.Marshal(session)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal session: %w", err)
	}

	if err := sc.client.Set(ctx, key, string(sessionData), duration); err != nil {
		return nil, err
	}

	return session, nil
}

// Get retrieves a session
func (sc *SessionCache) Get(ctx context.Context, sessionID string) (*Session, error) {
	key := sc.prefix + sessionID

	data, err := sc.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var session Session
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return &session, nil
}

// Update updates session data
func (sc *SessionCache) Update(ctx context.Context, sessionID string, data map[string]interface{}) error {
	session, err := sc.Get(ctx, sessionID)
	if err != nil {
		return err
	}

	// Merge new data
	for k, v := range data {
		session.Data[k] = v
	}

	key := sc.prefix + sessionID
	sessionData, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	ttl := time.Until(session.ExpiresAt)
	return sc.client.Set(ctx, key, string(sessionData), ttl)
}

// Delete deletes a session
func (sc *SessionCache) Delete(ctx context.Context, sessionID string) error {
	key := sc.prefix + sessionID
	return sc.client.Delete(ctx, key)
}

// Extend extends a session's expiration
func (sc *SessionCache) Extend(ctx context.Context, sessionID string, duration time.Duration) error {
	session, err := sc.Get(ctx, sessionID)
	if err != nil {
		return err
	}

	session.ExpiresAt = time.Now().Add(duration)

	key := sc.prefix + sessionID
	sessionData, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	return sc.client.Set(ctx, key, string(sessionData), duration)
}

// SetValue sets a value in session data
func (sc *SessionCache) SetValue(ctx context.Context, sessionID, key string, value interface{}) error {
	session, err := sc.Get(ctx, sessionID)
	if err != nil {
		return err
	}

	session.Data[key] = value

	sessionKey := sc.prefix + sessionID
	sessionData, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	ttl := time.Until(session.ExpiresAt)
	return sc.client.Set(ctx, sessionKey, string(sessionData), ttl)
}

// GetValue gets a value from session data
func (sc *SessionCache) GetValue(ctx context.Context, sessionID, key string) (interface{}, error) {
	session, err := sc.Get(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	value, exists := session.Data[key]
	if !exists {
		return nil, fmt.Errorf("key not found in session")
	}

	return value, nil
}

// DeleteValue deletes a value from session data
func (sc *SessionCache) DeleteValue(ctx context.Context, sessionID, key string) error {
	session, err := sc.Get(ctx, sessionID)
	if err != nil {
		return err
	}

	delete(session.Data, key)

	sessionKey := sc.prefix + sessionID
	sessionData, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	ttl := time.Until(session.ExpiresAt)
	return sc.client.Set(ctx, sessionKey, string(sessionData), ttl)
}

// Exists checks if a session exists
func (sc *SessionCache) Exists(ctx context.Context, sessionID string) (bool, error) {
	_, err := sc.Get(ctx, sessionID)
	if err == ErrCacheMiss {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetAllUserSessions gets all sessions for a user (requires key scanning)
func (sc *SessionCache) GetAllUserSessions(ctx context.Context, userID string) ([]*Session, error) {
	// Note: This would require Redis key scanning and filtering
	// Placeholder for more sophisticated session management
	return []*Session{}, nil
}

// InvalidateUserSessions invalidates all sessions for a user
func (sc *SessionCache) InvalidateUserSessions(ctx context.Context, userID string) error {
	// Note: This would require Redis key scanning
	// Placeholder for more sophisticated session invalidation
	return nil
}
