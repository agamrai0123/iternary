package redis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// Client wraps the Redis client with connection pooling and error handling
type Client struct {
	client    *redis.Client
	config    *Config
	mu        sync.RWMutex
	connected bool
	stats     *Stats
	ctx       context.Context
	cancel    context.CancelFunc
}

// Config holds Redis connection configuration
type Config struct {
	Host     string
	Port     int
	Password string
	DB       int
	// Connection pool settings
	MaxRetries   int
	PoolSize     int
	MinIdleConns int
	MaxConnAge   time.Duration
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	DialTimeout  time.Duration
	// Cache settings
	DefaultTTL     time.Duration
	MaxMemory      string // e.g., "512mb"
	EvictionPolicy string // e.g., "allkeys-lru"
}

// Stats tracks Redis operation metrics
type Stats struct {
	mu          sync.RWMutex
	Gets        int64
	Sets        int64
	Deletes     int64
	Hits        int64
	Misses      int64
	Errors      int64
	HitRate     float64
	AvgLatency  time.Duration
	LastUpdated time.Time
}

// NewClient creates a new Redis client with connection pooling
func NewClient(config *Config) *Client {
	ctx, cancel := context.WithCancel(context.Background())

	redisClient := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password:     config.Password,
		DB:           config.DB,
		MaxRetries:   config.MaxRetries,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		MaxConnAge:   config.MaxConnAge,
		IdleTimeout:  config.IdleTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		DialTimeout:  config.DialTimeout,
	})

	return &Client{
		client:    redisClient,
		config:    config,
		connected: false,
		stats: &Stats{
			LastUpdated: time.Now(),
		},
		ctx:    ctx,
		cancel: cancel,
	}
}

// Connect verifies the Redis connection
func (c *Client) Connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Use provided context or fallback to client context
	ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := c.client.Ping(ctxWithTimeout).Err(); err != nil {
		return fmt.Errorf("redis connection failed: %w", err)
	}

	c.connected = true
	return nil
}

// IsConnected returns the connection status
func (c *Client) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.connected
}

// Set stores a key-value pair with optional TTL
func (c *Client) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	c.recordSet()

	if ttl == 0 {
		ttl = c.config.DefaultTTL
	}

	err := c.client.Set(ctx, key, value, ttl).Err()
	if err != nil {
		c.recordError()
		return fmt.Errorf("redis set failed for key %s: %w", key, err)
	}

	return nil
}

// Get retrieves a value by key
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	c.recordGet()

	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		c.recordMiss()
		return "", ErrCacheMiss
	}
	if err != nil {
		c.recordError()
		return "", fmt.Errorf("redis get failed for key %s: %w", key, err)
	}

	c.recordHit()
	return val, nil
}

// GetBytes retrieves a byte value by key
func (c *Client) GetBytes(ctx context.Context, key string) ([]byte, error) {
	val, err := c.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	return []byte(val), nil
}

// Delete removes a key from cache
func (c *Client) Delete(ctx context.Context, keys ...string) error {
	c.recordDelete()

	if len(keys) == 0 {
		return nil
	}

	err := c.client.Del(ctx, keys...).Err()
	if err != nil {
		c.recordError()
		return fmt.Errorf("redis delete failed: %w", err)
	}

	return nil
}

// Exists checks if keys exist in cache
func (c *Client) Exists(ctx context.Context, keys ...string) (int64, error) {
	count, err := c.client.Exists(ctx, keys...).Result()
	if err != nil {
		c.recordError()
		return 0, fmt.Errorf("redis exists failed: %w", err)
	}
	return count, nil
}

// Expire sets expiration on a key
func (c *Client) Expire(ctx context.Context, key string, ttl time.Duration) error {
	ok, err := c.client.Expire(ctx, key, ttl).Result()
	if err != nil {
		c.recordError()
		return fmt.Errorf("redis expire failed: %w", err)
	}
	if !ok {
		return ErrKeyNotFound
	}
	return nil
}

// TTL returns the remaining TTL for a key
func (c *Client) TTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := c.client.TTL(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, ErrKeyNotFound
		}
		c.recordError()
		return 0, fmt.Errorf("redis ttl failed: %w", err)
	}
	return ttl, nil
}

// Increment increments a numeric key
func (c *Client) Increment(ctx context.Context, key string) (int64, error) {
	val, err := c.client.Incr(ctx, key).Result()
	if err != nil {
		c.recordError()
		return 0, fmt.Errorf("redis increment failed: %w", err)
	}
	return val, nil
}

// Decrement decrements a numeric key
func (c *Client) Decrement(ctx context.Context, key string) (int64, error) {
	val, err := c.client.Decr(ctx, key).Result()
	if err != nil {
		c.recordError()
		return 0, fmt.Errorf("redis decrement failed: %w", err)
	}
	return val, nil
}

// Append appends value to key
func (c *Client) Append(ctx context.Context, key, value string) (int64, error) {
	len, err := c.client.Append(ctx, key, value).Result()
	if err != nil {
		c.recordError()
		return 0, fmt.Errorf("redis append failed: %w", err)
	}
	return len, nil
}

// Flush clears all data from current database
func (c *Client) Flush(ctx context.Context) error {
	err := c.client.FlushDB(ctx).Err()
	if err != nil {
		c.recordError()
		return fmt.Errorf("redis flush failed: %w", err)
	}
	return nil
}

// Close closes the Redis connection
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cancel()
	c.connected = false
	return c.client.Close()
}

// Stats retrieval methods

// GetStats returns current cache statistics
func (c *Client) GetStats() Stats {
	c.stats.mu.RLock()
	defer c.stats.mu.RUnlock()

	stats := *c.stats

	// Calculate hit rate
	total := stats.Hits + stats.Misses
	if total > 0 {
		stats.HitRate = float64(stats.Hits) / float64(total) * 100
	}

	return stats
}

// ResetStats resets all statistics
func (c *Client) ResetStats() {
	c.stats.mu.Lock()
	defer c.stats.mu.Unlock()

	c.stats = &Stats{
		LastUpdated: time.Now(),
	}
}

// Internal stats recording methods

func (c *Client) recordGet() {
	c.stats.mu.Lock()
	defer c.stats.mu.Unlock()
	c.stats.Gets++
}

func (c *Client) recordSet() {
	c.stats.mu.Lock()
	defer c.stats.mu.Unlock()
	c.stats.Sets++
}

func (c *Client) recordDelete() {
	c.stats.mu.Lock()
	defer c.stats.mu.Unlock()
	c.stats.Deletes++
}

func (c *Client) recordHit() {
	c.stats.mu.Lock()
	defer c.stats.mu.Unlock()
	c.stats.Hits++
}

func (c *Client) recordMiss() {
	c.stats.mu.Lock()
	defer c.stats.mu.Unlock()
	c.stats.Misses++
}

func (c *Client) recordError() {
	c.stats.mu.Lock()
	defer c.stats.mu.Unlock()
	c.stats.Errors++
}
