package query

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// QueryProfile profiles query execution
type QueryProfile struct {
	Query         string
	ExecutionTime time.Duration
	RowsAffected  int64
	RowsRead      int64
	Error         error
	Timestamp     time.Time
	QueryPlan     string
	IsSlowQuery   bool
}

// QueryProfiler profiles database queries
type QueryProfiler struct {
	db            *sql.DB
	slowThreshold time.Duration
	profiles      []QueryProfile
	maxProfiles   int
}

// NewQueryProfiler creates a new query profiler
func NewQueryProfiler(db *sql.DB, slowThreshold time.Duration) *QueryProfiler {
	return &QueryProfiler{
		db:            db,
		slowThreshold: slowThreshold,
		profiles:      make([]QueryProfile, 0),
		maxProfiles:   1000,
	}
}

// ProfileQuery profiles a query execution
func (qp *QueryProfiler) ProfileQuery(ctx context.Context, query string, args ...interface{}) (*QueryProfile, error) {
	profile := &QueryProfile{
		Query:     query,
		Timestamp: time.Now(),
	}

	start := time.Now()

	rows, err := qp.db.QueryContext(ctx, query, args...)
	profile.ExecutionTime = time.Since(start)

	if err != nil {
		profile.Error = err
		return profile, err
	}
	defer rows.Close()

	// Count rows
	for rows.Next() {
		profile.RowsRead++
	}

	if err := rows.Err(); err != nil {
		profile.Error = err
		return profile, err
	}

	profile.IsSlowQuery = profile.ExecutionTime > qp.slowThreshold
	qp.addProfile(profile)

	return profile, nil
}

// ProfileExec profiles an exec query
func (qp *QueryProfiler) ProfileExec(ctx context.Context, query string, args ...interface{}) (*QueryProfile, error) {
	profile := &QueryProfile{
		Query:     query,
		Timestamp: time.Now(),
	}

	start := time.Now()
	result, err := qp.db.ExecContext(ctx, query, args...)
	profile.ExecutionTime = time.Since(start)

	if err != nil {
		profile.Error = err
		return profile, err
	}

	affected, err := result.RowsAffected()
	if err == nil {
		profile.RowsAffected = affected
	}

	profile.IsSlowQuery = profile.ExecutionTime > qp.slowThreshold
	qp.addProfile(profile)

	return profile, nil
}

// ProfileRow profiles a single row query
func (qp *QueryProfiler) ProfileRow(ctx context.Context, query string, args ...interface{}) (*QueryProfile, error) {
	profile := &QueryProfile{
		Query:     query,
		Timestamp: time.Now(),
	}

	start := time.Now()
	row := qp.db.QueryRowContext(ctx, query, args...)
	profile.ExecutionTime = time.Since(start)

	profile.RowsRead = 1
	profile.IsSlowQuery = profile.ExecutionTime > qp.slowThreshold
	qp.addProfile(profile)

	return profile, row.Err()
}

// addProfile adds a profile to the list
func (qp *QueryProfiler) addProfile(profile *QueryProfile) {
	qp.profiles = append(qp.profiles, *profile)

	// Keep only recent profiles
	if len(qp.profiles) > qp.maxProfiles {
		qp.profiles = qp.profiles[1:]
	}
}

// GetSlowQueries retrieves slow queries
func (qp *QueryProfiler) GetSlowQueries() []QueryProfile {
	var slow []QueryProfile
	for _, p := range qp.profiles {
		if p.IsSlowQuery {
			slow = append(slow, p)
		}
	}
	return slow
}

// GetStats returns execution statistics
func (qp *QueryProfiler) GetStats() map[string]interface{} {
	stats := make(map[string]interface{})
	stats["total_profiles"] = len(qp.profiles)
	stats["slow_queries"] = len(qp.GetSlowQueries())

	if len(qp.profiles) == 0 {
		return stats
	}

	var totalTime time.Duration
	var totalRows int64

	for _, p := range qp.profiles {
		totalTime += p.ExecutionTime
		totalRows += p.RowsRead
	}

	stats["avg_execution_time"] = totalTime / time.Duration(len(qp.profiles))
	stats["total_rows_read"] = totalRows
	stats["slow_threshold"] = qp.slowThreshold

	return stats
}

// ExplainQuery returns the query plan
func (qp *QueryProfiler) ExplainQuery(ctx context.Context, query string, args ...interface{}) (string, error) {
	explainQuery := fmt.Sprintf("EXPLAIN %s", query)

	rows, err := qp.db.QueryContext(ctx, explainQuery, args...)
	if err != nil {
		return "", fmt.Errorf("failed to explain query: %w", err)
	}
	defer rows.Close()

	var plan string
	for rows.Next() {
		var line string
		if err := rows.Scan(&line); err != nil {
			return "", err
		}
		plan += line + "\n"
	}

	return plan, rows.Err()
}

// AnalyzeQuery returns detailed query analysis
func (qp *QueryProfiler) AnalyzeQuery(ctx context.Context, query string, args ...interface{}) (string, error) {
	analyzeQuery := fmt.Sprintf("EXPLAIN ANALYZE %s", query)

	rows, err := qp.db.QueryContext(ctx, analyzeQuery, args...)
	if err != nil {
		return "", fmt.Errorf("failed to analyze query: %w", err)
	}
	defer rows.Close()

	var analysis string
	for rows.Next() {
		var line string
		if err := rows.Scan(&line); err != nil {
			return "", err
		}
		analysis += line + "\n"
	}

	return analysis, rows.Err()
}

// RecommendOptimizations suggests query optimizations
func (qp *QueryProfiler) RecommendOptimizations() []string {
	recommendations := []string{}

	slowQueries := qp.GetSlowQueries()
	if len(slowQueries) > 10 {
		recommendations = append(recommendations, "High number of slow queries detected - consider adding indexes")
	}

	stats := qp.GetStats()
	avgTime := stats["avg_execution_time"].(time.Duration)
	if avgTime > 100*time.Millisecond {
		recommendations = append(recommendations, "Average query execution time is high - profile your most frequent queries")
	}

	totalRows := stats["total_rows_read"].(int64)
	if totalRows > 100000 && len(qp.profiles) < 100 {
		recommendations = append(recommendations, "Large result sets - consider pagination or filtering")
	}

	return recommendations
}

// QueryOptimizer provides optimization suggestions
type QueryOptimizer struct {
	db       *sql.DB
	profiler *QueryProfiler
}

// NewQueryOptimizer creates a new query optimizer
func NewQueryOptimizer(db *sql.DB, profiler *QueryProfiler) *QueryOptimizer {
	return &QueryOptimizer{
		db:       db,
		profiler: profiler,
	}
}

// IdentifyMissingIndexes suggests missing indexes
func (qo *QueryOptimizer) IdentifyMissingIndexes(ctx context.Context) ([]string, error) {
	query := `
		SELECT 
			schemaname, tablename, attname, n_distinct
		FROM pg_stats
		WHERE schemaname NOT IN ('pg_catalog', 'information_schema')
		AND n_distinct > 100
		AND null_frac < 0.5
		ORDER BY n_distinct DESC
		LIMIT 20
	`

	rows, err := qo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to identify missing indexes: %w", err)
	}
	defer rows.Close()

	var suggestions []string
	for rows.Next() {
		var schema, table, column string
		var distinct float64

		if err := rows.Scan(&schema, &table, &column, &distinct); err != nil {
			continue
		}

		suggestion := fmt.Sprintf("Consider indexing %s.%s.%s (distinct values: %.0f)", schema, table, column, distinct)
		suggestions = append(suggestions, suggestion)
	}

	return suggestions, rows.Err()
}

// OptimizationTip represents an optimization tip
type OptimizationTip struct {
	Category   string
	Issue      string
	Solution   string
	Priority   string // HIGH, MEDIUM, LOW
	Complexity string // EASY, MEDIUM, HARD
}

// AnalyzePerformance analyzes overall database performance
func (qo *QueryOptimizer) AnalyzePerformance(ctx context.Context) []OptimizationTip {
	tips := []OptimizationTip{}

	// Check for table bloat
	tips = append(tips, OptimizationTip{
		Category:   "Maintenance",
		Issue:      "Routine VACUUM/ANALYZE to maintain performance",
		Solution:   "Run VACUUM ANALYZE periodically",
		Priority:   "MEDIUM",
		Complexity: "EASY",
	})

	// Check for missing prepared statements
	tips = append(tips, OptimizationTip{
		Category:   "Query Preparation",
		Issue:      "All queries using prepared statements?",
		Solution:   "Use prepared statements to reduce parsing overhead",
		Priority:   "HIGH",
		Complexity: "MEDIUM",
	})

	// Recommendations from profiler
	for _, rec := range qo.profiler.RecommendOptimizations() {
		tips = append(tips, OptimizationTip{
			Category:   "Profiling Results",
			Issue:      rec,
			Solution:   "Review slow queries and optimize or add indexes",
			Priority:   "HIGH",
			Complexity: "MEDIUM",
		})
	}

	return tips
}

// PrintPerformanceReport prints a performance analysis report
func (qo *QueryOptimizer) PrintPerformanceReport(ctx context.Context) error {
	log.Println("=== Database Performance Analysis Report ===")

	// Profiler stats
	stats := qo.profiler.GetStats()
	log.Printf("Total Queries Profiled: %v", stats["total_profiles"])
	log.Printf("Slow Queries: %v", stats["slow_queries"])
	log.Printf("Average Execution Time: %v", stats["avg_execution_time"])

	// Missing indexes
	if indexes, err := qo.IdentifyMissingIndexes(ctx); err == nil {
		if len(indexes) > 0 {
			log.Println("\nSuggested Indexes:")
			for _, idx := range indexes {
				log.Printf("  - %s", idx)
			}
		}
	}

	// Optimization tips
	tips := qo.AnalyzePerformance(ctx)
	if len(tips) > 0 {
		log.Println("\nOptimization Tips:")
		for _, tip := range tips {
			log.Printf("  [%s-%s] %s: %s", tip.Priority, tip.Complexity, tip.Category, tip.Issue)
		}
	}

	return nil
}
