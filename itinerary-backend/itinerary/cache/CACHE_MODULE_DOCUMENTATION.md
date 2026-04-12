# Cache Module Documentation

The cache module provides comprehensive caching solutions for the Itinerary application, including both in-memory and Redis-based implementations.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Components](#components)
  - [Memory Cache](#memory-cache)
  - [Query Cache](#query-cache)
  - [Session Cache](#session-cache)
  - [Rate Limiters](#rate-limiters)
- [Usage Examples](#usage-examples)
- [Configuration](#configuration)
- [Best Practices](#best-practices)
- [Performance Considerations](#performance-considerations)

## Overview

The cache module provides multiple caching strategies tailored to different use cases:

- **In-Memory Cache**: Fast, non-persistent cache suitable for development and single-instance deployments
- **Redis Cache**: Distributed cache for multi-instance deployments and persistence
- **Query Cache**: Caches database query results with TTL
- **Session Cache**: Manages user sessions with automatic expiration
- **Rate Limiters**: Implements multiple rate limiting algorithms

## Features

### Core Features
- ✅ Thread-safe operations
- ✅ Automatic cleanup of expired entries
- ✅ Flexible TTL configuration
- ✅ Multiple cache implementations
- ✅ Factory pattern for easy instantiation

### Advanced Features
- ✅ Sliding window rate limiting
- ✅ Token bucket rate limiting
- ✅ Session management with user tracking
- ✅ Query result caching with MD5 hashing
- ✅ Cache statistics tracking

## Installation

Ensure your project has the cache module imported:

```go
import (
    "yourapp/cache"
    "yourapp/cache/redis"
)
```

## Quick Start

### Using In-Memory Cache

```go
package main

import (
    "context"
    "fmt"
    "time"
    "yourapp/cache"
)

func main() {
    // Create memory cache
    memCache := cache.NewMemoryCache()
    defer memCache.Close()

    // Set a value
    memCache.Set("key1", "value1", 5*time.Minute)

    // Get a value
    value, err := memCache.Get("key1")
    if err == nil {
        fmt.Println("Retrieved:", value)
    }

    // Check existence
    exists := memCache.Exists("key1")
    fmt.Println("Key exists:", exists)

    // Delete a value
    memCache.Delete("key1")
}
```

### Using Cache Factory

```go
// Create with builder pattern
cacheInstance, err := cache.NewCacheBuilder().
    WithType(cache.Memory).
    WithDefaultTTL(10 * time.Minute).
    WithMaxSize(5000).
    Build()

if err != nil {
    panic(err)
}
defer cacheInstance.Close()
```

## Components

### Memory Cache

In-memory cache implementation with automatic cleanup.

#### Features
- Non-blocking reads with RWMutex
- Automatic expiration cleanup
- Thread-safe operations
- Support for various data types

#### Methods
- `Set(key string, value interface{}, ttl time.Duration) error`
- `Get(key string) (interface{}, error)`
- `Delete(key string) error`
- `Exists(key string) bool`
- `Clear() error`
- `SetIfNotExists(key, value, ttl) (bool, error)`
- `Increment(key string) error`
- `Decrement(key string) error`
- `GetAll() map[string]interface{}`
- `GetKeys() []string`
- `Size() int`

### Query Cache

Caches database query results with automatic invalidation support.

#### Usage Example

```go
package main

import (
    "context"
    "yourapp/cache/redis"
)

func main() {
    client := redis.NewClient("localhost:6379")
    defer client.Close()

    queryCache := redis.NewQueryCache(client)

    // Cache query result
    ctx := context.Background()
    err := queryCache.Set(ctx, 
        "SELECT * FROM users WHERE id=?", 
        userData, 
        5*time.Minute, 
        123) // userID as parameter

    // Retrieve cached result
    cached, err := queryCache.Get(ctx, 
        "SELECT * FROM users WHERE id=?", 
        123)

    // Cache with automatic fetch
    result, err := queryCache.CacheQueryResult(ctx,
        "SELECT * FROM users WHERE active=?",
        5*time.Minute,
        func() (interface{}, error) {
            return database.QueryActiveUsers()
        },
        true)
}
```

#### Parameters Used for Cache Keys
- Query string
- any additional parameters passed to Get/Set

### Session Cache

Manages user sessions with automatic expiration.

#### Features
- Automatic user session tracking
- Session data manipulation
- Session extension
- TTL management

#### Usage Example

```go
package main

import (
    "context"
    "yourapp/cache/redis"
)

func main() {
    client := redis.NewClient("localhost:6379")
    sessionCache := redis.NewSessionCache(client)

    ctx := context.Background()

    // Create session
    session, err := sessionCache.Create(ctx, 
        "session123", 
        "user456", 
        24*time.Hour)

    // Set session value
    err = sessionCache.SetValue(ctx, "session123", "cart_items", []string{"item1", "item2"})

    // Get session value
    cartItems, err := sessionCache.GetValue(ctx, "session123", "cart_items")

    // Extend session
    err = sessionCache.Extend(ctx, "session123", 24*time.Hour)

    // Check existence
    exists, err := sessionCache.Exists(ctx, "session123")

    // Delete session
    err = sessionCache.Delete(ctx, "session123")
}
```

### Rate Limiters

Three rate limiting strategies are implemented:

#### 1. Simple Rate Limiter

Basic fixed-window rate limiting.

```go
limiter := redis.NewRateLimiter(client)

config := redis.RateLimitConfig{
    MaxRequests:   100,
    WindowSize:    time.Minute,
    BlockDuration: time.Second,
}

allowed, err := limiter.IsAllowed(ctx, "user:123", config)
if !allowed {
    return errors.New("rate limit exceeded")
}

remaining, err := limiter.GetRemaining(ctx, "user:123", 100)
```

#### 2. Sliding Window Rate Limiter

More accurate rate limiting using time-based windows.

```go
sliding := redis.NewSlidingWindowLimiter(client)

allowed, err := sliding.IsAllowed(ctx, "user:123", 100, time.Minute)
if !allowed {
    fmt.Println("Rate limit exceeded")
}
```

#### 3. Token Bucket Rate Limiter

Fine-grained rate limiting for varying request costs.

```go
bucket := redis.NewTokenBucket(client, 100.0, 10.0) // 100 tokens, 10/sec refill

allowed, err := bucket.AllowRequest(ctx, "user:123", 5.0) // Requires 5 tokens
if !allowed {
    fmt.Println("Insufficient tokens")
}

available, err := bucket.GetAvailableTokens(ctx, "user:123")
fmt.Printf("Available tokens: %f\n", available)
```

## Usage Examples

### Example 1: Caching API Responses

```go
func getUser(ctx context.Context, userID int) (*User, error) {
    cache := getCacheInstance()
    
    // Try cache first
    key := fmt.Sprintf("user:%d", userID)
    if cached, err := cache.Get(key); err == nil {
        if user, ok := cached.(*User); ok {
            return user, nil
        }
    }

    // Fetch from database
    user, err := db.GetUser(ctx, userID)
    if err != nil {
        return nil, err
    }

    // Cache for 1 hour
    cache.Set(key, user, time.Hour)
    return user, nil
}
```

### Example 2: Rate Limiting API Endpoints

```go
func rateLimitMiddleware(next http.Handler) http.Handler {
    limiter := getRateLimiter()

    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        clientID := r.RemoteAddr
        config := redis.RateLimitConfig{
            MaxRequests:   100,
            WindowSize:    time.Minute,
        }

        allowed, _ := limiter.IsAllowed(r.Context(), clientID, config)
        if !allowed {
            http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
            return
        }

        next.ServeHTTP(w, r)
    })
}
```

### Example 3: Session Management

```go
func loginHandler(w http.ResponseWriter, r *http.Request) {
    sessionCache := getSessionCache()
    
    // Authenticate user
    user, err := authenticateUser(r)
    if err != nil {
        http.Error(w, "Authentication failed", http.StatusUnauthorized)
        return
    }

    // Create session
    sessionID := generateSessionID()
    session, err := sessionCache.Create(r.Context(), sessionID, user.ID, 24*time.Hour)
    if err != nil {
        http.Error(w, "Session creation failed", http.StatusInternalServerError)
        return
    }

    // Set in cookie
    http.SetCookie(w, &http.Cookie{
        Name:  "sessionID",
        Value: sessionID,
        Path:  "/",
    })

    json.NewEncoder(w).Encode(session)
}
```

## Configuration

### Creating a Custom Configuration

```go
config := &cache.Config{
    Type:       cache.Memory,
    DefaultTTL: 10 * time.Minute,
    MaxSize:    50000,
}

factory := cache.NewFactory(config)
cacheInstance, err := factory.Create()
```

### Redis Configuration

```go
config := &cache.Config{
    Type:      cache.Redis,
    RedisURL:  "redis://localhost:6379",
    RedisDB:   0,
    PoolSize:  20,
}

factory := cache.NewFactory(config)
redisCache, err := factory.Create()
```

## Best Practices

### 1. Choose the Right Cache Type
- Use **Memory Cache** for single-instance apps or development
- Use **Redis** for distributed systems or persistent cache

### 2. Set Appropriate TTLs
```go
// Short TTL for frequently changing data
cache.Set("user_preferences", data, 5*time.Minute)

// Longer TTL for static data
cache.Set("country_list", data, 24*time.Hour)
```

### 3. Implement Graceful Degradation
```go
val, err := cache.Get("key")
if err != nil {
    // Fallback to database
    val, err = database.Get(key)
    if err != nil {
        return handleError(err)
    }
}
```

### 4. Monitor Cache Hit Rates
```go
stats := queryCache.GetStats()
fmt.Printf("Cache hit rate: %.2f%%\n", stats.CacheHitRate*100)
```

### 5. Clean Up Resources
```go
defer cache.Close()
defer sessionCache.Delete(ctx, sessionID)
```

## Performance Considerations

### Memory Cache
- ✅ Very fast read/write operations
- ✅ No network latency
- ⚠️ Memory bound
- ⚠️ Single instance only

### Redis Cache
- ✅ Distributed across instances
- ✅ Persistent across restarts
- ✅ Scalable
- ⚠️ Network latency
- ⚠️ Requires Redis infrastructure

### Optimization Tips
1. Use appropriate TTL to balance freshness and performance
2. Implement cache invalidation strategies
3. Monitor cache size and clean up effectively
4. Use compression for large values
5. Batch operations when possible

## Troubleshooting

### Issue: Cache misses on every request
**Solution**: Check TTL configuration and clock synchronization

### Issue: Memory cache growing indefinitely
**Solution**: Enable cleanup goroutine or reduce TTL

### Issue: Redis connection failures
**Solution**: Implement retry logic and connection pooling

### Issue: Session data corruption
**Solution**: Ensure JSON marshaling/unmarshaling is correct

## Migration Guide

### From No Cache to Memory Cache

```go
// Before
result, _ := db.Query(sql)

// After
result, _ := cache.CacheQueryResult(ctx, sql, 5*time.Minute, 
    func() (interface{}, error) {
        return db.Query(sql)
    })
```

### From Memory Cache to Redis

```go
// Update configuration only
config.Type = cache.Redis
config.RedisURL = "redis://localhost:6379"
// No code changes needed due to interface abstraction
```
