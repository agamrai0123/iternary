package itinerary

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"
)

// PerformanceOptimizations contains optimization configurations
type PerformanceOptimizations struct {
	// Connection pooling
	MaxOpenConnections int
	MaxIdleConnections int
	ConnMaxLifetime    time.Duration

	// Cache settings
	CacheTTL              time.Duration
	CacheMaxSize          int
	CacheEvictionPolicy   string // "LRU", "LFU", "FIFO"

	// Rate limiting
	MaxRequestsPerSecond float64
	BurstSize            int

	// Request handling
	RequestTimeout      time.Duration
	IdleTimeout         time.Duration
	ReadHeaderTimeout   time.Duration
	MaxHeaderBytes      int

	// Go runtime
	MaxGoroutines       int
	MemoryAllocLimit    uint64
	GCTargetPercent     int // gc percent for GC trigger
}

// DefaultOptimizations returns recommended optimization settings
func DefaultOptimizations() *PerformanceOptimizations {
	return &PerformanceOptimizations{
		// Database connection pooling
		MaxOpenConnections: 25,  // Render default: 20, but allow some buffer
		MaxIdleConnections: 5,   // Keep 5 idle for quick reuse
		ConnMaxLifetime:    5 * time.Minute,

		// Cache optimization
		CacheTTL:            5 * time.Minute,
		CacheMaxSize:        10000,
		CacheEvictionPolicy: "LRU",

		// Rate limiting
		MaxRequestsPerSecond: 1000,
		BurstSize:           100,

		// HTTP timeouts
		RequestTimeout:    30 * time.Second,
		IdleTimeout:       90 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1MB

		// Go runtime
		MaxGoroutines:    250,
		MemoryAllocLimit: 256 * 1024 * 1024, // 256MB
		GCTargetPercent:  75,
	}
}

// ApplyDatabaseOptimizations configures database connection pool
func (po *PerformanceOptimizations) ApplyDatabaseOptimizations(db *sql.DB) {
	db.SetMaxOpenConns(po.MaxOpenConnections)
	db.SetMaxIdleConns(po.MaxIdleConnections)
	db.SetConnMaxLifetime(po.ConnMaxLifetime)

	fmt.Printf("✓ Database optimizations applied:\n")
	fmt.Printf("  - MaxOpenConnections: %d\n", po.MaxOpenConnections)
	fmt.Printf("  - MaxIdleConnections: %d\n", po.MaxIdleConnections)
	fmt.Printf("  - ConnMaxLifetime: %v\n", po.ConnMaxLifetime)
}

// QueryCache provides simple in-memory caching with TTL
type QueryCache struct {
	mu      sync.RWMutex
	data    map[string]cacheEntry
	maxSize int
	ttl     time.Duration
}

type cacheEntry struct {
	value     interface{}
	timestamp time.Time
	hits      int64
}

// NewQueryCache creates a new query cache
func NewQueryCache(maxSize int, ttl time.Duration) *QueryCache {
	cache := &QueryCache{
		data:    make(map[string]cacheEntry),
		maxSize: maxSize,
		ttl:     ttl,
	}

	// Start cleanup goroutine
	go cache.cleanupStaleEntries()

	return cache
}

// Get retrieves value from cache if not expired
func (qc *QueryCache) Get(key string) (interface{}, bool) {
	qc.mu.RLock()
	defer qc.mu.RUnlock()

	entry, exists := qc.data[key]
	if !exists {
		return nil, false
	}

	// Check if expired
	if time.Since(entry.timestamp) > qc.ttl {
		return nil, false
	}

	// Update hit counter
	entry.hits++
	return entry.value, true
}

// Set stores value in cache
func (qc *QueryCache) Set(key string, value interface{}) {
	qc.mu.Lock()
	defer qc.mu.Unlock()

	// Evict if cache is full
	if len(qc.data) >= qc.maxSize {
		qc.evictLRU()
	}

	qc.data[key] = cacheEntry{
		value:     value,
		timestamp: time.Now(),
		hits:      0,
	}
}

// evictLRU evicts least recently used entry
func (qc *QueryCache) evictLRU() {
	var lruKey string
	var lruTime time.Time

	for key, entry := range qc.data {
		if lruTime.IsZero() || entry.timestamp.Before(lruTime) {
			lruKey = key
			lruTime = entry.timestamp
		}
	}

	if lruKey != "" {
		delete(qc.data, lruKey)
	}
}

// cleanupStaleEntries removes expired entries periodically
func (qc *QueryCache) cleanupStaleEntries() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		qc.mu.Lock()
		now := time.Now()
		for key, entry := range qc.data {
			if now.Sub(entry.timestamp) > qc.ttl {
				delete(qc.data, key)
			}
		}
		qc.mu.Unlock()
	}
}

// Size returns current cache size
func (qc *QueryCache) Size() int {
	qc.mu.RLock()
	defer qc.mu.RUnlock()
	return len(qc.data)
}

// Stats returns cache statistics
type CacheStats struct {
	Size     int
	MaxSize  int
	TTL      time.Duration
	HitRatio float64
}

func (qc *QueryCache) Stats() CacheStats {
	qc.mu.RLock()
	defer qc.mu.RUnlock()

	totalHits := int64(0)
	for _, entry := range qc.data {
		totalHits += entry.hits
	}

	hitRatio := float64(0)
	if len(qc.data) > 0 {
		hitRatio = float64(totalHits) / float64(len(qc.data))
	}

	return CacheStats{
		Size:     len(qc.data),
		MaxSize:  qc.maxSize,
		TTL:      qc.ttl,
		HitRatio: hitRatio,
	}
}

// OptimizedQueryExecutor executes queries with caching and optimizations
type OptimizedQueryExecutor struct {
	db    *sql.DB
	cache *QueryCache
	opts  *PerformanceOptimizations
}

// NewOptimizedQueryExecutor creates new executor with optimizations
func NewOptimizedQueryExecutor(db *sql.DB, opts *PerformanceOptimizations) *OptimizedQueryExecutor {
	executor := &OptimizedQueryExecutor{
		db:    db,
		cache: NewQueryCache(opts.CacheMaxSize, opts.CacheTTL),
		opts:  opts,
	}

	// Apply database optimizations
	opts.ApplyDatabaseOptimizations(db)

	return executor
}

// QueryWithCache executes query with caching
func (oqe *OptimizedQueryExecutor) QueryWithCache(ctx context.Context, cacheKey string, query string, args ...interface{}) (*sql.Rows, error) {
	// Check cache first
	if cached, found := oqe.cache.Get(cacheKey); found {
		if rows, ok := cached.(*sql.Rows); ok {
			return rows, nil
		}
	}

	// Execute query with timeout
	ctx, cancel := context.WithTimeout(ctx, oqe.opts.RequestTimeout)
	defer cancel()

	rows, err := oqe.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	// Cache the result
	oqe.cache.Set(cacheKey, rows)

	return rows, nil
}

// HealthCheck verifies optimization settings are working
func (oqe *OptimizedQueryExecutor) HealthCheck() map[string]interface{} {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stats := map[string]interface{}{}

	// Check database connection
	if err := oqe.db.PingContext(ctx); err != nil {
		stats["database"] = "ERROR: " + err.Error()
	} else {
		stats["database"] = "OK"
	}

	// Get connection pool stats
	dbStats := oqe.db.Stats()
	stats["connections"] = map[string]interface{}{
		"open":           dbStats.OpenConnections,
		"in_use":         dbStats.InUse,
		"idle":           dbStats.Idle,
		"wait_count":     dbStats.WaitCount,
		"wait_duration":  dbStats.WaitDuration.String(),
		"max_idle_close": dbStats.MaxIdleClosed,
	}

	// Get cache stats
	cacheStats := oqe.cache.Stats()
	stats["cache"] = map[string]interface{}{
		"size":      cacheStats.Size,
		"max_size":  cacheStats.MaxSize,
		"ttl":       cacheStats.TTL.String(),
		"hit_ratio": fmt.Sprintf("%.2f%%", cacheStats.HitRatio*100),
	}

	return stats
}

// Example usage
/*
func main() {
	// Initialize database
	db, err := sql.Open("postgres", "...")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create optimizer with default settings
	opts := DefaultOptimizations()
	executor := NewOptimizedQueryExecutor(db, opts)

	// Use executor for queries
	rows, err := executor.QueryWithCache(
		context.Background(),
		"get_users_all",
		"SELECT * FROM users",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Check health
	health := executor.HealthCheck()
	fmt.Printf("System Health: %+v\n", health)
}
*/
