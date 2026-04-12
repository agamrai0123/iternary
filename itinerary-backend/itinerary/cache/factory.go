package cache

import (
	"fmt"
	"time"
)

// CacheType defines the type of cache to use
type CacheType string

const (
	// Memory-based cache
	Memory CacheType = "memory"
	// Redis-based cache
	Redis CacheType = "redis"
)

// Config contains cache configuration
type Config struct {
	Type       CacheType
	DefaultTTL time.Duration
	MaxSize    int
	RedisURL   string
	RedisDB    int
	PoolSize   int
}

// DefaultConfig returns the default cache configuration
func DefaultConfig() *Config {
	return &Config{
		Type:       Memory,
		DefaultTTL: 5 * time.Minute,
		MaxSize:    10000,
		RedisDB:    0,
		PoolSize:   10,
	}
}

// CacheInterface defines the cache interface
type CacheInterface interface {
	Set(key string, value interface{}, ttl time.Duration) error
	Get(key string) (interface{}, error)
	Delete(key string) error
	Exists(key string) bool
	Clear() error
	Close() error
}

// Factory creates cache instances
type Factory struct {
	config *Config
}

// NewFactory creates a new cache factory
func NewFactory(config *Config) *Factory {
	if config == nil {
		config = DefaultConfig()
	}
	return &Factory{config: config}
}

// Create creates a cache instance based on configuration
func (f *Factory) Create() (CacheInterface, error) {
	switch f.config.Type {
	case Memory:
		return f.createMemoryCache()
	case Redis:
		return f.createRedisCache()
	default:
		return nil, fmt.Errorf("unknown cache type: %s", f.config.Type)
	}
}

// createMemoryCache creates an in-memory cache
func (f *Factory) createMemoryCache() (CacheInterface, error) {
	return NewMemoryCache(), nil
}

// createRedisCache creates a Redis cache
func (f *Factory) createRedisCache() (CacheInterface, error) {
	// Note: This is a placeholder
	// In the actual implementation, you would initialize a Redis connection
	return nil, fmt.Errorf("Redis cache not yet implemented")
}

// CacheBuilder provides a fluent interface for building cache instances
type CacheBuilder struct {
	config *Config
}

// NewCacheBuilder creates a new cache builder
func NewCacheBuilder() *CacheBuilder {
	return &CacheBuilder{
		config: DefaultConfig(),
	}
}

// WithType sets the cache type
func (cb *CacheBuilder) WithType(cacheType CacheType) *CacheBuilder {
	cb.config.Type = cacheType
	return cb
}

// WithDefaultTTL sets the default TTL
func (cb *CacheBuilder) WithDefaultTTL(ttl time.Duration) *CacheBuilder {
	cb.config.DefaultTTL = ttl
	return cb
}

// WithMaxSize sets the maximum cache size
func (cb *CacheBuilder) WithMaxSize(size int) *CacheBuilder {
	cb.config.MaxSize = size
	return cb
}

// WithRedis configures Redis settings
func (cb *CacheBuilder) WithRedis(url string, db int, poolSize int) *CacheBuilder {
	cb.config.Type = Redis
	cb.config.RedisURL = url
	cb.config.RedisDB = db
	cb.config.PoolSize = poolSize
	return cb
}

// Build builds the cache instance
func (cb *CacheBuilder) Build() (CacheInterface, error) {
	factory := NewFactory(cb.config)
	return factory.Create()
}

// CacheManager provides a centralized cache management interface
type CacheManager struct {
	caches map[string]CacheInterface
	mu     chan bool // Simple mutex
}

// NewCacheManager creates a new cache manager
func NewCacheManager() *CacheManager {
	return &CacheManager{
		caches: make(map[string]CacheInterface),
		mu:     make(chan bool, 1),
	}
}

// Register registers a named cache
func (cm *CacheManager) Register(name string, cache CacheInterface) error {
	cm.mu <- true
	defer func() { <-cm.mu }()

	if _, exists := cm.caches[name]; exists {
		return fmt.Errorf("cache '%s' already registered", name)
	}

	cm.caches[name] = cache
	return nil
}

// Get retrieves a registered cache
func (cm *CacheManager) Get(name string) (CacheInterface, error) {
	cm.mu <- true
	defer func() { <-cm.mu }()

	cache, exists := cm.caches[name]
	if !exists {
		return nil, fmt.Errorf("cache '%s' not found", name)
	}

	return cache, nil
}

// Close closes all registered caches
func (cm *CacheManager) Close() error {
	cm.mu <- true
	defer func() { <-cm.mu }()

	for name, cache := range cm.caches {
		if err := cache.Close(); err != nil {
			return fmt.Errorf("error closing cache '%s': %w", name, err)
		}
	}

	return nil
}

// GetAll returns all registered caches
func (cm *CacheManager) GetAll() map[string]CacheInterface {
	cm.mu <- true
	defer func() { <-cm.mu }()

	result := make(map[string]CacheInterface)
	for name, cache := range cm.caches {
		result[name] = cache
	}

	return result
}
