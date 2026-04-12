# Query Optimization System Documentation

## Table of Contents
- [Overview](#overview)
- [Components](#components)
- [Installation & Setup](#installation--setup)
- [Quick Start](#quick-start)
- [Advanced Usage](#advanced-usage)
- [Monitoring](#monitoring)
- [Best Practices](#best-practices)
- [Performance Benchmarks](#performance-benchmarks)

## Overview

The Query Optimization System provides comprehensive tools for optimizing database performance:

- **Index Management**: Automatic index creation, analysis, and optimization
- **Connection Pooling**: Intelligent connection pool configuration
- **Query Profiling**: Track query execution times and identify slowdowns
- **Query Optimization**: Recommendations and pattern analysis
- **Performance Monitoring**: Real-time health checks and reporting

## Components

### 1. Index Manager (`indexes.go`)

Manages database indexes for optimal query performance.

**Features:**
- Create, drop, and manage indexes
- Composite and partial indexes
- Index statistics and analysis
- Unused index detection
- Slow query analysis

**Key Methods:**
```go
CreateIndex(ctx, index)           // Create an index
CreateCompositeIndex(ctx, name, table, cols...) // Multi-column index
CreatePartialIndex(ctx, name, table, condition, cols...) // Filtered index
ListIndexes(ctx, tableName)       // List all indexes
GetIndexStatistics(ctx, tableName) // Index usage statistics
UnusedIndexes(ctx)                // Find unused indexes
```

**Example Usage:**
```go
indexMgr := database.NewIndexManager(db)

// Create index
err := indexMgr.CreateIndex(ctx, &database.Index{
    Name:      "idx_users_email",
    TableName: "users",
    Columns:   []string{"email"},
    IsUnique:  true,
})

// Get statistics
stats, err := indexMgr.GetIndexStatistics(ctx, "users")
for _, stat := range stats {
    fmt.Printf("Index: %s, Scans: %d, Size: %d bytes\n", 
        stat.IndexName, stat.IndexScans, stat.IndexSize)
}
```

### 2. Connection Pool Manager (`pool.go`)

Optimizes database connection pooling.

**Features:**
- Pre-configured profiles (default, high-throughput, low-latency)
- Connection pool statistics
- Health monitoring
- Auto-cleanup of old connections

**Pool Configurations:**
```go
DefaultConfig()           // 25 max, 5 idle
HighThroughputConfig()   // 100 max, 20 idle
LowLatencyConfig()       // 50 max, 15 idle
```

**Example Usage:**
```go
poolMgr := database.NewPoolManager(db)

// Configure for high throughput
poolMgr.OptimizeForScenario("high-throughput")

// Monitor health
stats := poolMgr.GetStats()
fmt.Printf("Connections in use: %d\n", stats.InUse)
fmt.Printf("Wait count: %d\n", stats.WaitCount)

// Periodic monitoring
monitor := database.NewConnectionPoolMonitor(poolMgr, 30*time.Second)
monitor.Start()
defer monitor.Stop()
```

### 3. Query Profiler (`query_profiler.go`)

Profiles and analyzes query performance.

**Features:**
- Query execution time tracking
- Slow query detection
- Query plan analysis (EXPLAIN)
- Execution statistics

**Example Usage:**
```go
profiler := query.NewQueryProfiler(db, 100*time.Millisecond)

// Profile a query
profile, err := profiler.ProfileQuery(ctx, 
    "SELECT * FROM users WHERE id = $1", userID)

fmt.Printf("Execution time: %v\n", profile.ExecutionTime)
fmt.Printf("Rows read: %d\n", profile.RowsRead)
fmt.Printf("Is slow: %v\n", profile.IsSlowQuery)

// Get EXPLAIN plan
plan, err := profiler.ExplainQuery(ctx, query, args...)
fmt.Println(plan)

// Get detailed analysis
analysis, err := profiler.AnalyzeQuery(ctx, query, args...)
fmt.Println(analysis)
```

### 4. Query Optimizer (`query_optimizer.go`)

Provides optimization recommendations and analysis.

**Components:**
- Batch insert operations
- Pagination support
- JOIN optimization
- Query pattern analysis

**Example Usage:**
```go
optimizer := query.NewOptimizedQueryBuilder(db)

// Batch insert (more efficient than individual INSERTs)
err := optimizer.BatchInsert(ctx, "users",
    []string{"name", "email"},
    [][]interface{}{
        {"Alice", "alice@example.com"},
        {"Bob", "bob@example.com"},
        {"Charlie", "charlie@example.com"},
    })

// Paginated query
rows, err := optimizer.PaginatedQuery(ctx,
    "SELECT * FROM users ORDER BY id",
    pageSize=20, pageNum=2)

// Analyze query patterns
recs := query.AnalyzeQueryPattern("SELECT * FROM table")
for _, rec := range recs.HighPriority {
    fmt.Println("High Priority:", rec)
}
```

### 5. Optimization Module (`optimization_module.go`)

Main module that brings everything together.

**Features:**
- Unified initialization
- Comprehensive reporting
- Health checks
- Statistics collection

**Example Usage:**
```go
module := database.NewQueryOptimizationModule(db)

// Initialize all components
err := module.Initialize(ctx, "production")

// Run health check
err := module.HealthCheck()

// Get statistics
stats, err := module.GetStats(ctx)

// Print comprehensive report
module.PrintReport(ctx)

// Verify optimization setup
checklist := module.VerifyOptimization(ctx)
```

## Installation & Setup

### Step 1: Import Packages

```go
import (
    "yourapp/database"
    "yourapp/database/query"
)
```

### Step 2: Initialize Database Connection

```go
db, err := sql.Open("postgres", dsn)
if err != nil {
    panic(err)
}
```

### Step 3: Create and Initialize Module

```go
module := database.NewQueryOptimizationModule(db)
err := module.Initialize(ctx, "default")
```

## Quick Start

### Basic Setup

```go
package main

import (
    "context"
    "database/sql"
    "time"
    "yourapp/database"
)

func main() {
    // Connect to database
    db, _ := sql.Open("postgres", "connection_string")
    defer db.Close()

    ctx := context.Background()

    // Create optimization module
    module := database.NewQueryOptimizationModule(db)
    
    // Initialize
    module.Initialize(ctx, "default")
    
    // Check health
    module.HealthCheck()
    
    // Print report
    module.PrintReport(ctx)
}
```

### With Caching

```go
// Use cache + optimization together
cache := GetCacheInstance()
module := database.NewQueryOptimizationModule(db)

// Check cache and profile query
cacheKey := "users:1"
if cached, err := cache.Get(cacheKey); err == nil {
    return cached
}

// If not cached, profile and cache the result
profile, err := module.Profiler.ProfileQuery(ctx, query)
cache.Set(cacheKey, result, 1*time.Hour)
```

## Advanced Usage

### Custom Index Strategy

```go
indexMgr := database.NewIndexManager(db)

// Create composite index for common queries
indexMgr.CreateCompositeIndex(ctx, "idx_itinerary_user_status",
    "itineraries",
    "user_id", "status")

// Create partial index for active records only
indexMgr.CreatePartialIndex(ctx, "idx_active_users",
    "users",
    "is_active = true",
    "id")
```

### Performance Monitoring Loop

```go
// Monitor performance every 5 minutes
ticker := time.NewTicker(5 * time.Minute)
defer ticker.Stop()

for range ticker.C {
    stats, _ := module.GetStats(ctx)
    fmt.Printf("Slow queries: %d\n", stats.SlowQueryCount)
    fmt.Printf("Avg time: %v\n", stats.AvgQueryTime)
    fmt.Printf("Unused indexes: %d\n", stats.UnusedIndexCount)
}
```

### Optimized Batch Operations

```go
builder := query.NewOptimizedQueryBuilder(db)

// Batch insert 1000 records
values := make([][]interface{}, 1000)
for i := 0; i < 1000; i++ {
    values[i] = []interface{}{
        fmt.Sprintf("User %d", i),
        fmt.Sprintf("user%d@example.com", i),
    }
}

builder.BatchInsert(ctx, "users",
    []string{"name", "email"},
    values)

// More efficient than 1000 separate INSERTs
```

## Monitoring

### Real-time Health Checks

```go
// Start monitoring
monitor := database.NewConnectionPoolMonitor(poolMgr, 30*time.Second)
monitor.Start()
defer monitor.Stop()

// Alerts will be logged automatically
```

### Query Profiling Dashboard

```go
// Get profiling statistics
stats := module.Profiler.GetStats()
fmt.Printf("Total Queries: %v\n", stats["total_profiles"])
fmt.Printf("Slow Queries: %v\n", stats["slow_queries"])
fmt.Printf("Avg Time: %v\n", stats["avg_execution_time"])

// Identify slow queries
slowQueries := module.Profiler.GetSlowQueries()
for _, q := range slowQueries {
    fmt.Printf("%s took %v\n", q.Query, q.ExecutionTime)
}
```

### Index Usage Analysis

```go
// Find unused indexes
unused, _ := module.Indexes.UnusedIndexes(ctx)
fmt.Println("Unused indexes:", unused)

// Get index statistics
stats, _ := module.Indexes.GetIndexStatistics(ctx, "users")
for _, stat := range stats {
    fmt.Printf("Index %s: %d scans\n", stat.IndexName, stat.IndexScans)
}
```

## Best Practices

### 1. Index Design

✅ **DO:**
- Index foreign keys
- Index columns used in WHERE clauses
- Create composite indexes for common queries
- Use partial indexes for filtered queries

❌ **DON'T:**
- Create too many indexes (4-6 per table is typical)
- Index low-cardinality columns
- Forget to drop unused indexes

### 2. Query Patterns

✅ **DO:**
- Use prepared statements
- Paginate large result sets
- Use JOINs instead of subqueries
- Batch operations when possible

❌ **DON'T:**
- Use SELECT * (select needed columns)
- Forget WHERE clauses
- Use aggregations without GROUP BY properly
- Calculate in WHERE clause

### 3. Connection Pooling

✅ **DO:**
- Use appropriate pool sizes for your workload
- Monitor pool statistics
- Set connection max lifetime
- Use idle timeout

❌ **DON'T:**
- Create new connections for every query
- Use unlimited connection pools
- Leave connections idle indefinitely

### 4. Performance Monitoring

✅ **DO:**
- Profile slow queries
- Monitor query execution time
- Review query plans (EXPLAIN)
- Regularly analyze performance

❌ **DON'T:**
- Ignore slow queries
- Skip index analysis
- Forget to VACUUM/ANALYZE
- Deploy without performance testing

## Performance Benchmarks

Expected improvements with optimization:

| Scenario | Before | After | Improvement |
|----------|--------|-------|-------------|
| Simple SELECT | 2-5ms | 0.5-1ms | 4-5x faster |
| Join query | 10-20ms | 2-5ms | 3-4x faster |
| Batch insert (100 rows) | 50-100ms | 5-10ms | 5-10x faster |
| COUNT query | 100ms | 10-20ms | 5-10x faster |
| Pagination | 15-30ms | 2-5ms | 5-10x faster |

## Troubleshooting

### High Connection Wait Times
```go
// Solution: Increase pool size
poolMgr.Configure(&database.ConnectionPoolConfig{
    MaxOpenConnections: 50,
    MaxIdleConnections: 15,
})
```

### Slow Queries Not Improving
```go
// Solution: Analyze slow queries
plan, _ := module.Profiler.ExplainQuery(ctx, slowQuery)
fmt.Println(plan) // Look for sequential scans
```

### Unused Indexes Accumulating
```go
// Solution: Drop unused indexes
unusedIndexes, _ := module.Indexes.UnusedIndexes(ctx)
for _, idx := range unusedIndexes {
    module.Indexes.DropIndex(ctx, idx)
}
```

## Summary

The Query Optimization System provides:
- ✅ Comprehensive index management
- ✅ Connection pool optimization
- ✅ Query profiling and analysis
- ✅ Performance recommendations
- ✅ Real-time monitoring
- ✅ Production-ready features

For more examples, see [optimization_examples.go](./examples.go)
