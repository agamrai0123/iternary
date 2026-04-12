package redis

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"
)

// QueryCache caches database query results
type QueryCache struct {
	client *Client
	prefix string
}

// CacheEntry represents a cached query result
type CacheEntry struct {
	Query     string      `json:"query"`
	Data      interface{} `json:"data"`
	SavedAt   time.Time   `json:"saved_at"`
	ExpiresAt time.Time   `json:"expires_at"`
}

// NewQueryCache creates a new query cache
func NewQueryCache(client *Client) *QueryCache {
	return &QueryCache{
		client: client,
		prefix: "query:",
	}
}

// generateCacheKey creates a cache key from a query and parameters
func (qc *QueryCache) generateCacheKey(query string, params ...interface{}) string {
	hash := md5.Sum([]byte(fmt.Sprintf("%s:%v", query, params)))
	return qc.prefix + fmt.Sprintf("%x", hash)
}

// Set caches a query result
func (qc *QueryCache) Set(ctx context.Context, query string, data interface{}, ttl time.Duration, params ...interface{}) error {
	if ttl == 0 {
		ttl = 5 * time.Minute // Default cache TTL for queries
	}

	key := qc.generateCacheKey(query, params...)

	entry := CacheEntry{
		Query:     query,
		Data:      data,
		SavedAt:   time.Now(),
		ExpiresAt: time.Now().Add(ttl),
	}

	entryData, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal cache entry: %w", err)
	}

	return qc.client.Set(ctx, key, string(entryData), ttl)
}

// Get retrieves a cached query result
func (qc *QueryCache) Get(ctx context.Context, query string, params ...interface{}) (interface{}, error) {
	key := qc.generateCacheKey(query, params...)

	data, err := qc.client.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var entry CacheEntry
	if err := json.Unmarshal([]byte(data), &entry); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cache entry: %w", err)
	}

	// Check if entry expired
	if entry.ExpiresAt.Before(time.Now()) {
		_ = qc.Invalidate(ctx, query, params...)
		return nil, ErrCacheMiss
	}

	return entry.Data, nil
}

// Invalidate removes a cached query result
func (qc *QueryCache) Invalidate(ctx context.Context, query string, params ...interface{}) error {
	key := qc.generateCacheKey(query, params...)
	return qc.client.Delete(ctx, key)
}

// InvalidatePattern invalidates all queries matching a pattern
func (qc *QueryCache) InvalidatePattern(ctx context.Context, pattern string) error {
	// Note: Requires scanning keys (potential performance impact)
	// This is a placeholder for more sophisticated cache invalidation
	return nil
}

// GetCacheKey returns the cache key for a query (useful for direct operations)
func (qc *QueryCache) GetCacheKey(query string, params ...interface{}) string {
	return qc.generateCacheKey(query, params...)
}

// CacheQueryResult is a convenience method combining Get and Set operations
func (qc *QueryCache) CacheQueryResult(ctx context.Context, query string, ttl time.Duration,
	fetchFunc func() (interface{}, error), params ...interface{}) (interface{}, error) {

	// Try to get from cache first
	if cached, err := qc.Get(ctx, query, params...); err == nil {
		return cached, nil
	}

	// Cache miss, fetch from source
	data, err := fetchFunc()
	if err != nil {
		return nil, err
	}

	// Store in cache
	if err := qc.Set(ctx, query, data, ttl, params...); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("warning: failed to cache query result: %v\n", err)
	}

	return data, nil
}

// Stats

// CacheStats contains cache statistics
type CacheStats struct {
	TotalQueries  int64
	CachedQueries int64
	CacheHits     int64
	CacheMisses   int64
	CacheHitRate  float64
	EvictionCount int64
}

// GetStats would return cache statistics (needs additional tracking)
func (qc *QueryCache) GetStats() CacheStats {
	// Note: This would need additional tracking in the cache
	return CacheStats{}
}
