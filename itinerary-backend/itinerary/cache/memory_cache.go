package cache

import (
	"fmt"
	"sync"
	"time"
)

// MemoryCacheEntry represents an entry in the in-memory cache
type MemoryCacheEntry struct {
	Value     interface{}
	ExpiresAt time.Time
}

// MemoryCache is a simple in-memory cache implementation
type MemoryCache struct {
	mu    sync.RWMutex
	store map[string]*MemoryCacheEntry
	quit  chan bool
}

// NewMemoryCache creates a new in-memory cache
func NewMemoryCache() *MemoryCache {
	mc := &MemoryCache{
		store: make(map[string]*MemoryCacheEntry),
		quit:  make(chan bool),
	}

	// Start cleanup goroutine
	go mc.cleanupExpired()

	return mc
}

// Set sets a value in the cache with TTL
func (mc *MemoryCache) Set(key string, value interface{}, ttl time.Duration) error {
	if ttl == 0 {
		ttl = 5 * time.Minute // Default TTL
	}

	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.store[key] = &MemoryCacheEntry{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}

	return nil
}

// Get retrieves a value from the cache
func (mc *MemoryCache) Get(key string) (interface{}, error) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	entry, exists := mc.store[key]
	if !exists {
		return nil, fmt.Errorf("key not found")
	}

	if time.Now().After(entry.ExpiresAt) {
		// Entry expired, but we'll let cleanup goroutine handle removal
		return nil, fmt.Errorf("key expired")
	}

	return entry.Value, nil
}

// Delete removes a value from the cache
func (mc *MemoryCache) Delete(key string) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	delete(mc.store, key)
	return nil
}

// Exists checks if a key exists in the cache
func (mc *MemoryCache) Exists(key string) bool {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	entry, exists := mc.store[key]
	if !exists {
		return false
	}

	if time.Now().After(entry.ExpiresAt) {
		return false
	}

	return true
}

// Clear clears the entire cache
func (mc *MemoryCache) Clear() error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.store = make(map[string]*MemoryCacheEntry)
	return nil
}

// Size returns the number of entries in the cache
func (mc *MemoryCache) Size() int {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	return len(mc.store)
}

// cleanupExpired periodically removes expired entries
func (mc *MemoryCache) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mc.removeExpired()
		case <-mc.quit:
			return
		}
	}
}

// removeExpired removes all expired entries
func (mc *MemoryCache) removeExpired() {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	now := time.Now()
	for key, entry := range mc.store {
		if now.After(entry.ExpiresAt) {
			delete(mc.store, key)
		}
	}
}

// Close stops the cleanup goroutine
func (mc *MemoryCache) Close() error {
	close(mc.quit)
	return nil
}

// GetAll returns all non-expired entries
func (mc *MemoryCache) GetAll() map[string]interface{} {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	result := make(map[string]interface{})
	now := time.Now()

	for key, entry := range mc.store {
		if now.Before(entry.ExpiresAt) {
			result[key] = entry.Value
		}
	}

	return result
}

// GetKeys returns all non-expired keys
func (mc *MemoryCache) GetKeys() []string {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	var keys []string
	now := time.Now()

	for key, entry := range mc.store {
		if now.Before(entry.ExpiresAt) {
			keys = append(keys, key)
		}
	}

	return keys
}

// SetIfNotExists sets a value only if the key doesn't exist
func (mc *MemoryCache) SetIfNotExists(key string, value interface{}, ttl time.Duration) (bool, error) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	entry, exists := mc.store[key]
	if exists && time.Now().Before(entry.ExpiresAt) {
		return false, nil
	}

	if ttl == 0 {
		ttl = 5 * time.Minute
	}

	mc.store[key] = &MemoryCacheEntry{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}

	return true, nil
}

// Increment increments a numeric value
func (mc *MemoryCache) Increment(key string) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	entry, exists := mc.store[key]
	if !exists {
		mc.store[key] = &MemoryCacheEntry{
			Value:     1,
			ExpiresAt: time.Now().Add(5 * time.Minute),
		}
		return nil
	}

	val, ok := entry.Value.(int)
	if !ok {
		return fmt.Errorf("value is not an integer")
	}

	entry.Value = val + 1
	return nil
}

// Decrement decrements a numeric value
func (mc *MemoryCache) Decrement(key string) error {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	entry, exists := mc.store[key]
	if !exists {
		mc.store[key] = &MemoryCacheEntry{
			Value:     -1,
			ExpiresAt: time.Now().Add(5 * time.Minute),
		}
		return nil
	}

	val, ok := entry.Value.(int)
	if !ok {
		return fmt.Errorf("value is not an integer")
	}

	entry.Value = val - 1
	return nil
}
