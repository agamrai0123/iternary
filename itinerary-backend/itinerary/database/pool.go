package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// ConnectionPoolConfig contains connection pool configuration
type ConnectionPoolConfig struct {
	MaxOpenConnections    int
	MaxIdleConnections    int
	ConnectionMaxLifetime time.Duration
	ConnectionMaxIdleTime time.Duration
}

// DefaultConnectionPoolConfig returns optimized connection pool settings
func DefaultConnectionPoolConfig() *ConnectionPoolConfig {
	return &ConnectionPoolConfig{
		MaxOpenConnections:    25,
		MaxIdleConnections:    5,
		ConnectionMaxLifetime: 5 * time.Minute,
		ConnectionMaxIdleTime: 2 * time.Minute,
	}
}

// HighThroughputConnectionPoolConfig for high-traffic scenarios
func HighThroughputConnectionPoolConfig() *ConnectionPoolConfig {
	return &ConnectionPoolConfig{
		MaxOpenConnections:    100,
		MaxIdleConnections:    20,
		ConnectionMaxLifetime: 5 * time.Minute,
		ConnectionMaxIdleTime: 1 * time.Minute,
	}
}

// LowLatencyConnectionPoolConfig for low-latency priority
func LowLatencyConnectionPoolConfig() *ConnectionPoolConfig {
	return &ConnectionPoolConfig{
		MaxOpenConnections:    50,
		MaxIdleConnections:    15,
		ConnectionMaxLifetime: 3 * time.Minute,
		ConnectionMaxIdleTime: 30 * time.Second,
	}
}

// PoolManager manages database connection pools
type PoolManager struct {
	db  *sql.DB
	cfg *ConnectionPoolConfig
}

// NewPoolManager creates a new pool manager
func NewPoolManager(db *sql.DB) *PoolManager {
	return &PoolManager{
		db:  db,
		cfg: DefaultConnectionPoolConfig(),
	}
}

// Configure applies connection pool settings
func (pm *PoolManager) Configure(cfg *ConnectionPoolConfig) error {
	if cfg == nil {
		cfg = DefaultConnectionPoolConfig()
	}

	pm.db.SetMaxOpenConns(cfg.MaxOpenConnections)
	pm.db.SetMaxIdleConns(cfg.MaxIdleConnections)
	pm.db.SetConnMaxLifetime(cfg.ConnectionMaxLifetime)
	pm.db.SetConnMaxIdleTime(cfg.ConnectionMaxIdleTime)

	pm.cfg = cfg
	log.Printf("Connection pool configured: MaxOpen=%d, MaxIdle=%d, MaxLifetime=%v",
		cfg.MaxOpenConnections, cfg.MaxIdleConnections, cfg.ConnectionMaxLifetime)

	return nil
}

// PoolStats represents connection pool statistics
type PoolStats struct {
	OpenConnections   int32
	InUse             int32
	Idle              int32
	WaitCount         int64
	WaitDuration      time.Duration
	MaxIdleClosed     int64
	MaxLifetimeClosed int64
}

// GetStats retrieves connection pool statistics
func (pm *PoolManager) GetStats() *PoolStats {
	dbStats := pm.db.Stats()

	return &PoolStats{
		OpenConnections:   int32(dbStats.OpenConnections),
		InUse:             int32(dbStats.InUse),
		Idle:              int32(dbStats.Idle),
		WaitCount:         dbStats.WaitCount,
		WaitDuration:      dbStats.WaitDuration,
		MaxIdleClosed:     dbStats.MaxIdleClosed,
		MaxLifetimeClosed: dbStats.MaxLifetimeClosed,
	}
}

// PrintStats prints pool statistics
func (pm *PoolManager) PrintStats() {
	stats := pm.GetStats()
	log.Printf("Connection Pool Stats:")
	log.Printf("  Open Connections: %d", stats.OpenConnections)
	log.Printf("  In Use: %d", stats.InUse)
	log.Printf("  Idle: %d", stats.Idle)
	log.Printf("  Wait Count: %d", stats.WaitCount)
	log.Printf("  Wait Duration: %v", stats.WaitDuration)
	log.Printf("  Max Idle Closed: %d", stats.MaxIdleClosed)
	log.Printf("  Max Lifetime Closed: %d", stats.MaxLifetimeClosed)
}

// HealthCheck verifies the connection pool is healthy
func (pm *PoolManager) HealthCheck() error {
	stats := pm.db.Stats()

	if stats.OpenConnections == 0 {
		return fmt.Errorf("no open connections in pool")
	}

	// Check if we're hitting max connections too often
	if stats.WaitCount > 1000 && stats.WaitDuration > time.Second {
		return fmt.Errorf("connection pool under stress: waitCount=%d, waitDuration=%v", stats.WaitCount, stats.WaitDuration)
	}

	return nil
}

// OptimizeForScenario configures pool based on usage scenario
func (pm *PoolManager) OptimizeForScenario(scenario string) error {
	var cfg *ConnectionPoolConfig

	switch scenario {
	case "high-throughput":
		cfg = HighThroughputConnectionPoolConfig()
	case "low-latency":
		cfg = LowLatencyConnectionPoolConfig()
	case "default":
		cfg = DefaultConnectionPoolConfig()
	default:
		return fmt.Errorf("unknown scenario: %s", scenario)
	}

	return pm.Configure(cfg)
}

// ConnectionPoolMonitor monitors pool health
type ConnectionPoolMonitor struct {
	pm       *PoolManager
	ticker   *time.Ticker
	stopChan chan bool
}

// NewConnectionPoolMonitor creates a new pool monitor
func NewConnectionPoolMonitor(pm *PoolManager, interval time.Duration) *ConnectionPoolMonitor {
	return &ConnectionPoolMonitor{
		pm:       pm,
		ticker:   time.NewTicker(interval),
		stopChan: make(chan bool),
	}
}

// Start starts monitoring the connection pool
func (cpm *ConnectionPoolMonitor) Start() {
	go func() {
		for {
			select {
			case <-cpm.ticker.C:
				stats := cpm.pm.GetStats()
				if stats.OpenConnections > int32(cpm.pm.cfg.MaxOpenConnections*9/10) {
					log.Printf("Warning: Connection pool approaching limit: %d/%d",
						stats.OpenConnections, cpm.pm.cfg.MaxOpenConnections)
				}

				if stats.WaitCount > 100 {
					log.Printf("Warning: High connection wait count: %d", stats.WaitCount)
				}

			case <-cpm.stopChan:
				return
			}
		}
	}()
}

// Stop stops monitoring
func (cpm *ConnectionPoolMonitor) Stop() {
	cpm.ticker.Stop()
	cpm.stopChan <- true
}

// StatementCache represents prepared statement cache configuration
type StatementCacheConfig struct {
	MaxCachedStatements int
	StaleThreshold      time.Duration
}

// DefaultStatementCacheConfig returns default statement cache settings
func DefaultStatementCacheConfig() *StatementCacheConfig {
	return &StatementCacheConfig{
		MaxCachedStatements: 100,
		StaleThreshold:      1 * time.Hour,
	}
}

// StatementCacheEntry represents a cached prepared statement
type StatementCacheEntry struct {
	Query     string
	Statement *sql.Stmt
	CreatedAt time.Time
	LastUsed  time.Time
	UseCount  int64
}

// StatementCache caches prepared statements
type StatementCache struct {
	db       *sql.DB
	cache    map[string]*StatementCacheEntry
	config   *StatementCacheConfig
	stopChan chan bool
}

// NewStatementCache creates a new statement cache
func NewStatementCache(db *sql.DB) *StatementCache {
	return &StatementCache{
		db:       db,
		cache:    make(map[string]*StatementCacheEntry),
		config:   DefaultStatementCacheConfig(),
		stopChan: make(chan bool),
	}
}

// GetStatement retrieves or creates a prepared statement
func (sc *StatementCache) GetStatement(query string) (*sql.Stmt, error) {
	if entry, exists := sc.cache[query]; exists {
		entry.LastUsed = time.Now()
		entry.UseCount++
		return entry.Statement, nil
	}

	stmt, err := sc.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}

	sc.cache[query] = &StatementCacheEntry{
		Query:     query,
		Statement: stmt,
		CreatedAt: time.Now(),
		LastUsed:  time.Now(),
		UseCount:  1,
	}

	return stmt, nil
}

// Size returns the number of cached statements
func (sc *StatementCache) Size() int {
	return len(sc.cache)
}

// Clear clears all cached statements
func (sc *StatementCache) Clear() error {
	for _, entry := range sc.cache {
		if err := entry.Statement.Close(); err != nil {
			return fmt.Errorf("failed to close statement: %w", err)
		}
	}
	sc.cache = make(map[string]*StatementCacheEntry)
	return nil
}

// Close closes the statement cache
func (sc *StatementCache) Close() error {
	return sc.Clear()
}
