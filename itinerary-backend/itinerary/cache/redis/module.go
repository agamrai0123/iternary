// Package redis provides Redis-based caching utilities for the Itinerary application
package redis

import (
	"errors"
)

var (
	ErrCacheMiss = errors.New("cache miss")
	ErrNilValue  = errors.New("nil value")
)

// Cache is the main cache interface for all implementations
type Cache interface {
	Get(key string) (string, error)
	Set(key string, value string, ttl int) error
	Delete(key string) error
	Exists(key string) bool
	Increment(key string) error
	Decrement(key string) error
}

// CacheModule contains all Redis cache utilities
type CacheModule struct {
	Query   *QueryCache
	Session *SessionCache
	Limiter *RateLimiter
	Sliding *SlidingWindowLimiter
	Tokens  *TokenBucket
}

// NewCacheModule creates a new cache module with all utilities
func NewCacheModule(client *Client) *CacheModule {
	return &CacheModule{
		Query:   NewQueryCache(client),
		Session: NewSessionCache(client),
		Limiter: NewRateLimiter(client),
		Sliding: NewSlidingWindowLimiter(client),
		Tokens:  NewTokenBucket(client, 100, 1), // Default: 100 tokens, 1 per second
	}
}
