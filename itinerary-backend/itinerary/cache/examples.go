package cache

import (
	"fmt"
	"log"
	"time"
)

// ExampleMemoryCacheUsage demonstrates in-memory cache usage
func ExampleMemoryCacheUsage() {
	cache := NewMemoryCache()
	defer cache.Close()

	// Set values
	cache.Set("user:1", map[string]string{"name": "Alice", "email": "alice@example.com"}, 5*time.Minute)
	cache.Set("user:2", map[string]string{"name": "Bob", "email": "bob@example.com"}, 5*time.Minute)

	// Get values
	if val, err := cache.Get("user:1"); err == nil {
		fmt.Printf("Retrieved user: %v\n", val)
	}

	// Check existence
	if exists := cache.Exists("user:1"); exists {
		fmt.Println("User 1 is cached")
	}

	// Get all keys
	keys := cache.GetKeys()
	fmt.Printf("Cached keys: %v\n", keys)

	// Delete
	cache.Delete("user:2")
}

// ExampleCacheFactoryUsage demonstrates factory pattern
func ExampleCacheFactoryUsage() {
	// Using builder
	cache, err := NewCacheBuilder().
		WithType(Memory).
		WithDefaultTTL(10 * time.Minute).
		WithMaxSize(5000).
		Build()

	if err != nil {
		log.Fatal(err)
	}
	defer cache.Close()

	cache.Set("app_config", map[string]string{"version": "1.0.0"}, time.Hour)
}

// ExampleQueryCacheUsage demonstrates query caching
func ExampleQueryCacheUsage() {
	// Initialize Redis client
	// client := redis.NewClient("localhost:6379")
	// defer client.Close()

	// queryCache := redis.NewQueryCache(client)
	// ctx := context.Background()

	// // Cache a query result
	// query := "SELECT * FROM users WHERE id = ?"
	// userID := 123

	// // First call - fetches from database
	// var userData interface{}
	// // userData, err := database.Query(query, userID)

	// err := queryCache.Set(ctx, query, userData, 5*time.Minute, userID)
	// if err != nil {
	//     log.Printf("Error caching query: %v\n", err)
	// }

	// // Subsequent calls - fetches from cache
	// cached, err := queryCache.Get(ctx, query, userID)
	// if err == nil {
	//     fmt.Printf("Cached result: %v\n", cached)
	// }

	// // Invalidate cache
	// err = queryCache.Invalidate(ctx, query, userID)
	// if err != nil {
	//     log.Printf("Error invalidating cache: %v\n", err)
	// }
}

// ExampleSessionCacheUsage demonstrates session management
func ExampleSessionCacheUsage() {
	// Initialize Redis client
	// client := redis.NewClient("localhost:6379")
	// sessionCache := redis.NewSessionCache(client)
	// ctx := context.Background()

	// // Create session
	// session, err := sessionCache.Create(ctx, "session:abc123", "user:456", 24*time.Hour)
	// if err != nil {
	//     log.Fatal(err)
	// }

	// // Set session values
	// err = sessionCache.SetValue(ctx, "session:abc123", "cart_items", []string{"item1", "item2"})
	// err = sessionCache.SetValue(ctx, "session:abc123", "preferences", map[string]string{"theme": "dark"})

	// // Get session values
	// cartItems, err := sessionCache.GetValue(ctx, "session:abc123", "cart_items")
	// fmt.Printf("Cart items: %v\n", cartItems)

	// // Extend session
	// err = sessionCache.Extend(ctx, "session:abc123", 24*time.Hour)

	// // Check session existence
	// exists, err := sessionCache.Exists(ctx, "session:abc123")
	// fmt.Printf("Session exists: %v\n", exists)

	// // Delete session
	// err = sessionCache.Delete(ctx, "session:abc123")
}

// ExampleRateLimiterUsage demonstrates rate limiting
func ExampleRateLimiterUsage() {
	// Initialize Redis client
	// client := redis.NewClient("localhost:6379")
	// limiter := redis.NewRateLimiter(client)
	// ctx := context.Background()

	// config := redis.RateLimitConfig{
	//     MaxRequests:   100,
	//     WindowSize:    time.Minute,
	//     BlockDuration: time.Second,
	// }

	// for i := 0; i < 150; i++ {
	//     allowed, err := limiter.IsAllowed(ctx, "user:789", config)
	//     if !allowed {
	//         fmt.Printf("Request %d: Rate limit exceeded\n", i)
	//     } else {
	//         fmt.Printf("Request %d: Allowed\n", i)
	//     }

	//     remaining, _ := limiter.GetRemaining(ctx, "user:789", 100)
	//     fmt.Printf("Remaining requests: %d\n", remaining)
	// }
}

// ExampleSlidingWindowLimiterUsage demonstrates sliding window rate limiting
func ExampleSlidingWindowLimiterUsage() {
	// Initialize Redis client
	// client := redis.NewClient("localhost:6379")
	// sliding := redis.NewSlidingWindowLimiter(client)
	// ctx := context.Background()

	// // Allow 50 requests per minute
	// for i := 0; i < 100; i++ {
	//     allowed, err := sliding.IsAllowed(ctx, "api:user:789", 50, time.Minute)
	//     if !allowed {
	//         fmt.Printf("Request %d: Rate limit exceeded (sliding window)\n", i)
	//     } else {
	//         fmt.Printf("Request %d: Allowed\n", i)
	//     }
	//     time.Sleep(100 * time.Millisecond)
	// }
}

// ExampleTokenBucketUsage demonstrates token bucket rate limiting
func ExampleTokenBucketUsage() {
	// Initialize Redis client
	// client := redis.NewClient("localhost:6379")
	// bucket := redis.NewTokenBucket(client, 100.0, 5.0) // 100 tokens, 5/sec refill
	// ctx := context.Background()

	// // Expensive operations
	// operations := []struct {
	//     name   string
	//     tokens float64
	// }{
	//     {"small_operation", 1.0},
	//     {"medium_operation", 5.0},
	//     {"large_operation", 20.0},
	// }

	// for _, op := range operations {
	//     allowed, err := bucket.AllowRequest(ctx, "api:premium", op.tokens)
	//     if !allowed {
	//         fmt.Printf("%s: Insufficient tokens\n", op.name)
	//     } else {
	//         fmt.Printf("%s: Allowed\n", op.name)
	//     }

	//     available, _ := bucket.GetAvailableTokens(ctx, "api:premium")
	//     fmt.Printf("Available tokens: %f\n", available)
	// }
}

// ExampleCacheManagerUsage demonstrates cache manager usage
func ExampleCacheManagerUsage() {
	manager := NewCacheManager()
	defer manager.Close()

	// Create different cache instances
	queryCache := NewMemoryCache()
	sessionCache := NewMemoryCache()

	// Register caches
	manager.Register("queries", queryCache)
	manager.Register("sessions", sessionCache)

	// Use caches
	if cache, err := manager.Get("queries"); err == nil {
		cache.Set("key1", "value1", time.Hour)
	}

	// List all caches
	allCaches := manager.GetAll()
	fmt.Printf("Registered caches: %v\n", len(allCaches))
}

// ExampleCacheWithFallback demonstrates cache with fallback
func ExampleCacheWithFallback(cache CacheInterface, key string) (interface{}, error) {
	// Try cache
	value, err := cache.Get(key)
	if err == nil {
		fmt.Println("Cache hit")
		return value, nil
	}

	// Fallback to source
	fmt.Println("Cache miss - fetching from source")
	// value, err := fetchFromSource(key)
	// if err != nil {
	//     return nil, err
	// }

	// Store in cache for future use
	// cache.Set(key, value, 5*time.Minute)

	return value, nil
}

// ExampleConcurrentCacheUsage demonstrates concurrent cache operations
func ExampleConcurrentCacheUsage() {
	cache := NewMemoryCache()
	defer cache.Close()

	// Set values concurrently
	for i := 0; i < 100; i++ {
		go func(id int) {
			key := fmt.Sprintf("key:%d", id)
			cache.Set(key, id, 5*time.Minute)
		}(i)
	}

	// Get values concurrently
	for i := 0; i < 100; i++ {
		go func(id int) {
			key := fmt.Sprintf("key:%d", id)
			if val, err := cache.Get(key); err == nil {
				fmt.Printf("Got %s = %v\n", key, val)
			}
		}(i)
	}

	time.Sleep(1 * time.Second)
	fmt.Printf("Final cache size: %d\n", cache.Size())
}

// ExampleMultiTierCache demonstrates multi-tier caching
type MultiTierCacheExample struct {
	l1Cache CacheInterface // In-memory
	l2Cache CacheInterface // Redis
}

func (mtc *MultiTierCacheExample) Get(key string) (interface{}, error) {
	// Try L1 first
	if val, err := mtc.l1Cache.Get(key); err == nil {
		fmt.Println("L1 cache hit")
		return val, nil
	}

	// Try L2
	if val, err := mtc.l2Cache.Get(key); err == nil {
		fmt.Println("L2 cache hit")
		mtc.l1Cache.Set(key, val, 5*time.Minute) // Populate L1
		return val, nil
	}

	fmt.Println("Cache miss")
	return nil, fmt.Errorf("key not found in any cache tier")
}

func (mtc *MultiTierCacheExample) Set(key string, value interface{}) error {
	mtc.l1Cache.Set(key, value, 5*time.Minute)
	mtc.l2Cache.Set(key, value, time.Hour)
	return nil
}
