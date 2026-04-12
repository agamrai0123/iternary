# Cache Module Integration Guide

This guide shows how to integrate the cache module into your Itinerary backend application.

## Table of Contents
- [Installation](#installation)
- [Configuration](#configuration)
- [Integration with HTTP Middleware](#integration-with-http-middleware)
- [Integration with Database Layer](#integration-with-database-layer)
- [Integration with Authentication](#integration-with-authentication)
- [Performance Optimization](#performance-optimization)
- [Monitoring and Debugging](#monitoring-and-debugging)
- [Migration Path](#migration-path)

## Installation

### Step 1: Add Cache Package

Ensure the cache package is in your project structure:

```
itinerary-backend/
├── itinerary/
│   ├── cache/
│   │   ├── memory_cache.go
│   │   ├── factory.go
│   │   ├── examples.go
│   │   ├── cache_test.go
│   │   └── redis/
│   │       ├── client.go
│   │       ├── query_cache.go
│   │       ├── session_cache.go
│   │       ├── rate_limiter.go
│   │       └── module.go
```

### Step 2: Import in your application

```go
import (
    "yourapp/itinerary/cache"
    "yourapp/itinerary/cache/redis"
)
```

## Configuration

### Application-Level Configuration

Create a cache configuration in your app initialization:

```go
// main.go
package main

import (
    "time"
    "yourapp/cache"
)

func initializeCache() cache.CacheInterface {
    // For development
    return cache.NewMemoryCache()
    
    // Or using builder pattern
    // return cache.NewCacheBuilder().
    //     WithType(cache.Memory).
    //     WithDefaultTTL(5 * time.Minute).
    //     WithMaxSize(50000).
    //     Build()
}

func main() {
    cacheInstance := initializeCache()
    defer cacheInstance.Close()
    
    // Use cache in your application
}
```

### Environment-Based Configuration

```go
func initializeCache() cache.CacheInterface {
    env := os.Getenv("ENVIRONMENT")
    
    if env == "production" {
        // Use Redis in production
        return cache.NewCacheBuilder().
            WithType(cache.Redis).
            WithRedis("redis://prod-redis:6379", 0, 20).
            Build()
    } else {
        // Use memory cache in development
        return cache.NewCacheBuilder().
            WithType(cache.Memory).
            WithDefaultTTL(5 * time.Minute).
            Build()
    }
}
```

## Integration with HTTP Middleware

### 1. Rate Limiting Middleware

```go
// middleware/rate_limiter.go
package middleware

import (
    "net/http"
    "yourapp/cache/redis"
)

func RateLimitMiddleware(limiter *redis.RateLimiter) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            clientID := getClientID(r) // IP address or user ID
            
            config := redis.RateLimitConfig{
                MaxRequests:   100,
                WindowSize:    time.Minute,
                BlockDuration: time.Second,
            }
            
            allowed, err := limiter.IsAllowed(r.Context(), clientID, config)
            if err != nil {
                http.Error(w, "Internal error", http.StatusInternalServerError)
                return
            }
            
            if !allowed {
                w.Header().Set("Retry-After", "60")
                http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}

func getClientID(r *http.Request) string {
    // Try to get user ID from context
    if userID, ok := r.Context().Value("userID").(string); ok {
        return "user:" + userID
    }
    // Fallback to IP address
    return "ip:" + r.RemoteAddr
}
```

### 2. Cache Control Middleware

```go
// middleware/cache_control.go
package middleware

import (
    "net/http"
    "time"
)

func CacheControlMiddleware(duration time.Duration) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Only cache GET requests
            if r.Method != "GET" {
                next.ServeHTTP(w, r)
                return
            }
            
            w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", int(duration.Seconds())))
            w.Header().Set("Expires", time.Now().Add(duration).Format(http.TimeFormat))
            
            next.ServeHTTP(w, r)
        })
    }
}
```

## Integration with Database Layer

### 1. Query Result Caching

```go
// repository/user_repository.go
package repository

import (
    "context"
    "fmt"
    "time"
    "yourapp/cache"
    "yourapp/cache/redis"
    "yourapp/models"
)

type UserRepository struct {
    db         *sql.DB
    queryCache *redis.QueryCache
}

func NewUserRepository(db *sql.DB, queryCache *redis.QueryCache) *UserRepository {
    return &UserRepository{
        db:         db,
        queryCache: queryCache,
    }
}

func (ur *UserRepository) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
    query := "SELECT id, name, email, created_at FROM users WHERE id = ?"
    
    // Try cache first
    cached, err := ur.queryCache.Get(ctx, query, userID)
    if err == nil {
        if user, ok := cached.(*models.User); ok {
            return user, nil
        }
    }
    
    // Fetch from database
    var user models.User
    err = ur.db.QueryRowContext(ctx, query, userID).Scan(
        &user.ID,
        &user.Name,
        &user.Email,
        &user.CreatedAt,
    )
    
    if err != nil {
        return nil, err
    }
    
    // Cache for 1 hour
    ur.queryCache.Set(ctx, query, &user, time.Hour, userID)
    
    return &user, nil
}

func (ur *UserRepository) GetAllUsers(ctx context.Context) ([]*models.User, error) {
    query := "SELECT id, name, email, created_at FROM users"
    
    return ur.queryCache.CacheQueryResult(ctx,
        query,
        time.Hour,
        func() (interface{}, error) {
            rows, err := ur.db.QueryContext(ctx, query)
            if err != nil {
                return nil, err
            }
            defer rows.Close()
            
            var users []*models.User
            for rows.Next() {
                var user models.User
                err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
                if err != nil {
                    return nil, err
                }
                users = append(users, &user)
            }
            
            return users, rows.Err()
        },
    )
}

func (ur *UserRepository) InvalidateUserCache(ctx context.Context, userID int) error {
    query := "SELECT id, name, email, created_at FROM users WHERE id = ?"
    return ur.queryCache.Invalidate(ctx, query, userID)
}
```

### 2. Cache Invalidation on Updates

```go
func (ur *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
    query := "UPDATE users SET name = ?, email = ? WHERE id = ?"
    
    _, err := ur.db.ExecContext(ctx, query, user.Name, user.Email, user.ID)
    if err != nil {
        return err
    }
    
    // Invalidate cache for this user
    return ur.InvalidateUserCache(ctx, user.ID)
}
```

## Integration with Authentication

### 1. Session Management

```go
// auth/session_manager.go
package auth

import (
    "context"
    "fmt"
    "time"
    "yourapp/cache/redis"
    "yourapp/models"
)

type SessionManager struct {
    sessionCache *redis.SessionCache
}

func NewSessionManager(sessionCache *redis.SessionCache) *SessionManager {
    return &SessionManager{
        sessionCache: sessionCache,
    }
}

func (sm *SessionManager) CreateSession(ctx context.Context, user *models.User) (string, error) {
    sessionID := generateSessionID()
    
    session, err := sm.sessionCache.Create(ctx, sessionID, user.ID, 24*time.Hour)
    if err != nil {
        return "", err
    }
    
    // Store user information
    sm.sessionCache.SetValue(ctx, sessionID, "user_name", user.Name)
    sm.sessionCache.SetValue(ctx, sessionID, "user_email", user.Email)
    sm.sessionCache.SetValue(ctx, sessionID, "user_role", user.Role)
    
    return sessionID, nil
}

func (sm *SessionManager) GetSession(ctx context.Context, sessionID string) (*models.Session, error) {
    session, err := sm.sessionCache.Get(ctx, sessionID)
    if err != nil {
        return nil, fmt.Errorf("session not found")
    }
    
    return &models.Session{
        ID:     session.ID,
        UserID: session.UserID,
        Data:   session.Data,
    }, nil
}

func (sm *SessionManager) InvalidateSession(ctx context.Context, sessionID string) error {
    return sm.sessionCache.Delete(ctx, sessionID)
}

func (sm *SessionManager) RefreshSession(ctx context.Context, sessionID string) error {
    return sm.sessionCache.Extend(ctx, sessionID, 24*time.Hour)
}
```

### 2. Authentication Middleware

```go
// middleware/auth.go
package middleware

import (
    "net/http"
    "context"
    "yourapp/auth"
)

func AuthMiddleware(sm *auth.SessionManager) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            sessionID, err := r.Cookie("sessionID")
            if err != nil {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }
            
            session, err := sm.GetSession(r.Context(), sessionID.Value)
            if err != nil {
                http.Error(w, "Session expired", http.StatusUnauthorized)
                return
            }
            
            // Add session to context
            ctx := context.WithValue(r.Context(), "session", session)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

## Performance Optimization

### 1. Multi-Tier Caching

```go
// cache/multi_tier.go
package cache

import (
    "fmt"
    "time"
)

type MultiTierCache struct {
    l1 CacheInterface // Memory cache (fast, limited size)
    l2 CacheInterface // Redis cache (slower, larger)
}

func NewMultiTierCache(l1, l2 CacheInterface) *MultiTierCache {
    return &MultiTierCache{l1: l1, l2: l2}
}

func (mtc *MultiTierCache) Get(key string) (interface{}, error) {
    // Try L1 first (memory)
    if val, err := mtc.l1.Get(key); err == nil {
        return val, nil
    }
    
    // Try L2 (Redis)
    if val, err := mtc.l2.Get(key); err == nil {
        // Populate L1
        mtc.l1.Set(key, val, 5*time.Minute)
        return val, nil
    }
    
    return nil, fmt.Errorf("key not found")
}

func (mtc *MultiTierCache) Set(key string, value interface{}, ttl time.Duration) error {
    mtc.l1.Set(key, value, ttl)
    return mtc.l2.Set(key, value, ttl)
}
```

### 2. Cache Warming

```go
// cache/warmer.go
package cache

type CacheWarmer struct {
    cache CacheInterface
}

func NewCacheWarmer(cache CacheInterface) *CacheWarmer {
    return &CacheWarmer{cache: cache}
}

func (cw *CacheWarmer) WarmCache(ctx context.Context, repo *repository.Repository) error {
    // Load frequently accessed data
    countries, err := repo.GetAllCountries(ctx)
    if err != nil {
        return err
    }
    
    for _, country := range countries {
        key := fmt.Sprintf("country:%d", country.ID)
        cw.cache.Set(key, country, 24*time.Hour)
    }
    
    return nil
}
```

## Monitoring and Debugging

### 1. Cache Metrics

```go
// cache/metrics.go
package cache

type CacheMetrics struct {
    TotalRequests  int64
    CacheHits      int64
    CacheMisses    int64
    EvictionCount  int64
}

func (cm *CacheMetrics) RecordHit() {
    cm.CacheHits++
    cm.TotalRequests++
}

func (cm *CacheMetrics) RecordMiss() {
    cm.CacheMisses++
    cm.TotalRequests++
}

func (cm *CacheMetrics) HitRate() float64 {
    if cm.TotalRequests == 0 {
        return 0
    }
    return float64(cm.CacheHits) / float64(cm.TotalRequests)
}
```

### 2. Debug Logging

```go
// cache/logger.go
package cache

import (
    "log"
    "fmt"
)

type CacheLogger struct {
    enabled bool
}

func NewCacheLogger(enabled bool) *CacheLogger {
    return &CacheLogger{enabled: enabled}
}

func (cl *CacheLogger) LogGet(key string, hit bool) {
    if !cl.enabled {
        return
    }
    status := "miss"
    if hit {
        status = "hit"
    }
    log.Printf("[CACHE] GET %s: %s\n", key, status)
}

func (cl *CacheLogger) LogSet(key string, ttl time.Duration) {
    if !cl.enabled {
        return
    }
    log.Printf("[CACHE] SET %s (TTL: %v)\n", key, ttl)
}

func (cl *CacheLogger) LogDelete(key string) {
    if !cl.enabled {
        return
    }
    log.Printf("[CACHE] DELETE %s\n", key)
}
```

## Migration Path

### Phase 1: Development (Week 1)
- Use in-memory cache
- Focus on functionality
- No persistence needed

### Phase 2: Testing (Week 2)
- Introduce Redis for testing
- Verify cache invalidation
- Performance testing

### Phase 3: Production (Week 3+)
- Deploy Redis cluster
- Configure replication
- Monitor cache performance

### Migration Example

```go
// Start with memory cache
cache, _ := cache.NewCacheBuilder().
    WithType(cache.Memory).
    Build()

// Later, switch to Redis by just changing one line
cache, _ := cache.NewCacheBuilder().
    WithType(cache.Redis).
    WithRedis("redis://redis:6379", 0, 20).
    Build()

// Code using cache doesn't change!
```

## Troubleshooting

### Issue: Cache not working in production

**Solution**: Verify Redis connection and ensure cache initialization happens before handlers start.

```go
func main() {
    cache, err := initializeCache()
    if err != nil {
        log.Fatalf("Failed to initialize cache: %v", err)
    }
    defer cache.Close()
    
    // Initialize handlers after cache
    setupHandlers(cache)
}
```

### Issue: Memory usage increasing

**Solution**: Check TTL settings and implement cleanup goroutine.

```go
cache := cache.NewMemoryCache()
// Cleanup runs automatically every 1 minute
defer cache.Close()
```

### Issue: Cache not invalidating on updates

**Solution**: Always invalidate cache when data is updated.

```go
// UPDATE pattern
err := updateUserInDB(user)
if err == nil {
    queryCache.Invalidate(ctx, userQuery, userID)
}
```

## Summary

The cache module provides:
- ✅ Flexible caching strategies
- ✅ Easy integration with existing code
- ✅ Minimal performance overhead
- ✅ Production-ready features
- ✅ Comprehensive testing support

For more examples, see [examples.go](./examples.go)
