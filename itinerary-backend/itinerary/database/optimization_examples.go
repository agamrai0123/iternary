package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"yourapp/database/query"
)

// ExampleInitializeOptimization demonstrates initialization
func ExampleInitializeOptimization(db *sql.DB) error {
	ctx := context.Background()

	// Create optimization module
	module := NewQueryOptimizationModule(db)

	// Initialize with production settings
	err := module.Initialize(ctx, "production")
	if err != nil {
		return fmt.Errorf("failed to initialize optimization: %w", err)
	}

	// Verify health
	if err := module.HealthCheck(); err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}

	fmt.Println("Optimization initialized successfully")
	return nil
}

// ExampleProfileQueryExecution demonstrates query profiling
func ExampleProfileQueryExecution(db *sql.DB) error {
	ctx := context.Background()

	// Create profiler
	profiler := query.NewQueryProfiler(db, 100*time.Millisecond)

	// Profile a query
	profile, err := profiler.ProfileQuery(ctx,
		"SELECT id, name, email FROM users WHERE status = $1 LIMIT 100",
		"active")

	if err != nil {
		return err
	}

	fmt.Printf("Query: %s\n", profile.Query)
	fmt.Printf("Execution time: %v\n", profile.ExecutionTime)
	fmt.Printf("Rows read: %d\n", profile.RowsRead)
	fmt.Printf("Is slow query: %v\n", profile.IsSlowQuery)

	return nil
}

// ExampleCreateOptimalIndexes demonstrates index creation
func ExampleCreateOptimalIndexes(db *sql.DB) error {
	ctx := context.Background()

	indexMgr := NewIndexManager(db)

	// Create single-column index
	err := indexMgr.CreateIndex(ctx, &Index{
		Name:      "idx_users_email",
		TableName: "users",
		Columns:   []string{"email"},
		IsUnique:  true,
	})
	if err != nil {
		return err
	}

	// Create composite index
	err = indexMgr.CreateCompositeIndex(ctx,
		"idx_itineraries_user_status",
		"itineraries",
		"user_id", "status")
	if err != nil {
		return err
	}

	// Create partial index
	err = indexMgr.CreatePartialIndex(ctx,
		"idx_active_destinations",
		"destinations",
		"is_active = true",
		"id")
	if err != nil {
		return err
	}

	fmt.Println("Optimal indexes created")
	return nil
}

// ExampleConfigureConnectionPool demonstrates pool configuration
func ExampleConfigureConnectionPool(db *sql.DB) {
	poolMgr := NewPoolManager(db)

	// For production with high traffic
	poolMgr.OptimizeForScenario("high-throughput")

	// Get statistics
	stats := poolMgr.GetStats()
	fmt.Printf("Open Connections: %d\n", stats.OpenConnections)
	fmt.Printf("In Use: %d\n", stats.InUse)
	fmt.Printf("Idle: %d\n", stats.Idle)
	fmt.Printf("Wait Duration: %v\n", stats.WaitDuration)

	// Monitor pool periodically
	monitor := NewConnectionPoolMonitor(poolMgr, 30*time.Second)
	monitor.Start()
	defer monitor.Stop()
}

// ExampleBatchInsert demonstrates efficient batch operations
func ExampleBatchInsert(db *sql.DB) error {
	ctx := context.Background()
	builder := query.NewOptimizedQueryBuilder(db)

	// Prepare batch data
	values := [][]interface{}{
		{"Alice", "alice@example.com", "2024-04-01"},
		{"Bob", "bob@example.com", "2024-04-02"},
		{"Charlie", "charlie@example.com", "2024-04-03"},
	}

	// Insert all at once (more efficient than separate INSERTs)
	err := builder.BatchInsert(ctx, "users",
		[]string{"name", "email", "created_at"},
		values)

	if err != nil {
		return fmt.Errorf("batch insert failed: %w", err)
	}

	fmt.Println("Batch insert completed")
	return nil
}

// ExamplePaginatedQuery demonstrates pagination
func ExamplePaginatedQuery(db *sql.DB) error {
	ctx := context.Background()
	builder := query.NewOptimizedQueryBuilder(db)

	// Get page 2 with 20 items per page
	rows, err := builder.PaginatedQuery(ctx,
		"SELECT id, name, email FROM users ORDER BY id",
		pageSize=20,
		pageNum=2)

	if err != nil {
		return err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var id int
		var name, email string
		if err := rows.Scan(&id, &name, &email); err != nil {
			continue
		}
		count++
	}

	fmt.Printf("Retrieved %d items from page 2\n", count)
	return nil
}

// ExampleAnalyzeSlowQueries demonstrates slow query analysis
func ExampleAnalyzeSlowQueries(db *sql.DB) error {
	ctx := context.Background()

	profiler := query.NewQueryProfiler(db, 50*time.Millisecond)
	optimizer := query.NewQueryOptimizer(db, profiler)

	// Run some queries to profile
	for i := 0; i < 100; i++ {
		_, _ = profiler.ProfileQuery(ctx,
			"SELECT COUNT(*) FROM itineraries WHERE user_id = $1",
			i%10)
	}

	// Analyze slow queries
	// tips := optimizer.AnalyzePerformance(ctx)
	// for _, tip := range tips {
	//     fmt.Printf("[%s] %s: %s\n", tip.Priority, tip.Category, tip.Issue)
	// }

	return nil
}

// ExampleQueryOptimizationPatterns demonstrates query optimization patterns
func ExampleQueryOptimizationPatterns() {
	patterns := query.CommonPatterns()

	fmt.Println("Common Query Optimization Patterns:\n")
	for i, pattern := range patterns {
		fmt.Printf("%d. %s\n", i+1, pattern.Name)
		fmt.Printf("   Description: %s\n", pattern.Description)
		fmt.Printf("   Example: %s\n", pattern.Query)
		fmt.Printf("   Benefit: %s\n\n", pattern.Performance)
	}
}

// ExampleAnalyzeQueryPattern demonstrates query pattern analysis
func ExampleAnalyzeQueryPattern() {
	// Analyze a BAD query
	badQuery := "SELECT * FROM users WHERE LOWER(name) LIKE '%john%'"
	recs := query.AnalyzeQueryPattern(badQuery)

	fmt.Println("Query Analysis for: " + badQuery)
	fmt.Println("\nHigh Priority Issues:")
	for _, rec := range recs.HighPriority {
		fmt.Println("  - " + rec)
	}

	fmt.Println("\nMedium Priority Issues:")
	for _, rec := range recs.MediumPriority {
		fmt.Println("  - " + rec)
	}
}

// ExamplePerformanceMonitoring demonstrates continuous monitoring
func ExamplePerformanceMonitoring(db *sql.DB) {
	module := NewQueryOptimizationModule(db)
	ctx := context.Background()

	// Start monitoring
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			stats, _ := module.GetStats(ctx)
			fmt.Printf("[Monitor] Slow queries: %d, Avg time: %v\n",
				stats.SlowQueryCount, stats.AvgQueryTime)

			// Alert if too many slow queries
			if stats.SlowQueryCount > 10 {
				fmt.Println("⚠️ Warning: High number of slow queries detected!")
			}
		}
	}()

	// Run for a while
	time.Sleep(30 * time.Second)
}

// ExampleComprehensiveReport demonstrates full reporting
func ExampleComprehensiveReport(db *sql.DB) error {
	ctx := context.Background()
	module := NewQueryOptimizationModule(db)

	// Initialize
	err := module.Initialize(ctx, "production")
	if err != nil {
		return err
	}

	// Run some queries to generate profiling data
	builder := query.NewOptimizedQueryBuilder(db)
	for i := 0; i < 50; i++ {
		_, _ = builder.PaginatedQuery(ctx,
			"SELECT id, name FROM users ORDER BY id",
			20, 1)
	}

	// Print comprehensive report
	return module.PrintReport(ctx)
}

// ExampleOptimizationVerification demonstrates verification
func ExampleOptimizationVerification(db *sql.DB) error {
	ctx := context.Background()
	module := NewQueryOptimizationModule(db)

	err := module.Initialize(ctx, "production")
	if err != nil {
		return err
	}

	// Verify optimization setup
	checklist := module.VerifyOptimization(ctx)

	fmt.Println("Optimization Verification:")
	fmt.Printf("✓ Connection Pooling: %v\n", checklist.HasConnectionPooling)
	fmt.Printf("✓ Indexes: %v\n", checklist.HasIndexes)
	fmt.Printf("✓ Pagination: %v\n", checklist.UsesPagination)
	fmt.Printf("✓ Result Caching: %v\n", checklist.CachesResults)
	fmt.Printf("✓ Monitoring: %v\n", checklist.MonitoringEnabled)
	fmt.Printf("✓ Prepared Statements: %v\n", checklist.UsesQuerryPreparedStatements)

	return nil
}

// ExampleMultiTierOptimization demonstrates combining multiple optimizations
func ExampleMultiTierOptimization(db *sql.DB) error {
	ctx := context.Background()

	// 1. Configure connection pool
	poolMgr := NewPoolManager(db)
	poolMgr.OptimizeForScenario("high-throughput")

	// 2. Create indexes
	indexMgr := NewIndexManager(db)
	indexMgr.InitializeOptimalIndexes(ctx)

	// 3. Profile queries
	profiler := query.NewQueryProfiler(db, 100*time.Millisecond)

	// 4. Build optimized queries
	builder := query.NewOptimizedQueryBuilder(db)

	// Execute with all optimizations
	rows, err := builder.PaginatedQuery(ctx,
		"SELECT id, name, email FROM users WHERE status = $1 ORDER BY created_at DESC",
		pageSize=20,
		pageNum=1,
		"active")

	if err != nil {
		return err
	}
	defer rows.Close()

	// Profile the execution
	profile, _ := profiler.ProfileQuery(ctx,
		"SELECT id, name, email FROM users WHERE status = 'active'")

	fmt.Printf("Rows: %d, Time: %v\n", profile.RowsRead, profile.ExecutionTime)
	return nil
}

// ExampleIndexOptimization demonstrates index-focused optimization
func ExampleIndexOptimization(db *sql.DB) error {
	ctx := context.Background()
	indexMgr := NewIndexManager(db)

	// Find unused indexes
	unused, err := indexMgr.UnusedIndexes(ctx)
	if err == nil && len(unused) > 0 {
		fmt.Println("Unused indexes found:")
		for _, idx := range unused {
			fmt.Printf("  - %s (consider dropping)\n", idx)
		}
	}

	// Get index statistics
	stats, err := indexMgr.GetIndexStatistics(ctx, "users")
	if err == nil {
		fmt.Println("\nIndex Statistics for 'users' table:")
		for _, stat := range stats {
			fmt.Printf("  - %s: %d scans, %d bytes\n",
				stat.IndexName, stat.IndexScans, stat.IndexSize)
		}
	}

	return nil
}

// ExampleQueryPlanAnalysis demonstrates query plan analysis
func ExampleQueryPlanAnalysis(db *sql.DB) error {
	ctx := context.Background()
	profiler := query.NewQueryProfiler(db, 100*time.Millisecond)

	query := "SELECT * FROM users WHERE email = $1"

	// Get query plan
	plan, err := profiler.ExplainQuery(ctx, query, "test@example.com")
	if err != nil {
		return err
	}

	fmt.Println("Query Plan:")
	fmt.Println(plan)

	// Get detailed analysis
	analysis, err := profiler.AnalyzeQuery(ctx, query, "test@example.com")
	if err != nil {
		return err
	}

	fmt.Println("\nDetailed Analysis:")
	fmt.Println(analysis)

	return nil
}
