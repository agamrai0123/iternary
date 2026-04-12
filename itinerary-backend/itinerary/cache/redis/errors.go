package redis

import "errors"

// Cache error definitions
var (
	ErrCacheMiss    = errors.New("cache miss: key not found")
	ErrKeyNotFound  = errors.New("key not found in cache")
	ErrNotConnected = errors.New("redis client not connected")
	ErrInvalidTTL   = errors.New("invalid TTL value")
)
