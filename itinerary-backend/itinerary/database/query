package query

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// OptimizedQueryBuilder helps build optimized database queries
type OptimizedQueryBuilder struct {
	db *sql.DB
}

// NewOptimizedQueryBuilder creates a new optimized query builder
func NewOptimizedQueryBuilder(db *sql.DB) *OptimizedQueryBuilder {
	return &OptimizedQueryBuilder{db: db}
}

// BatchInsert efficiently inserts multiple rows
func (oqb *OptimizedQueryBuilder) BatchInsert(ctx context.Context, table string, columns []string, values [][]interface{}) error {
	if len(values) == 0 {
		return fmt.Errorf("no values to insert")
	}

	// Build insert query with multiple value sets
	columnStr := strings.Join(columns, ", ")
	placeholders := make([]string, 0, len(values))
	args := make([]interface{}, 0)

	paramIdx := 1
	for _, row := range values {
		rowPlaceholders := make([]string, len(row))
		for i, val := range row {
			rowPlaceholders[i] = fmt.Sprintf("$%d", paramIdx)
			args = append(args, val)
			paramIdx++
		}
		placeholders = append(placeholders, "("+strings.Join(rowPlaceholders, ", ")+")")
	}

	query := fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES %s",
		table,
		columnStr,
		strings.Join(placeholders, ", "),
	)

	_, err := oqb.db.ExecContext(ctx, query, args...)
	return err
}

// PaginatedQuery executes a paginated query efficiently
func (oqb *OptimizedQueryBuilder) PaginatedQuery(ctx context.Context, baseQuery string, pageSize, pageNum int, args ...interface{}) (*sql.Rows, error) {
	if pageSize <= 0 || pageNum <= 0 {
		return nil, fmt.Errorf("invalid page size or number")
	}

	offset := (pageNum - 1) * pageSize

	// Append LIMIT and OFFSET
	query := fmt.Sprintf("%s LIMIT %d OFFSET %d", baseQuery, pageSize, offset)

	return oqb.db.QueryContext(ctx, query, args...)
}

// BulkUpdate efficiently updates multiple rows
func (oqb *OptimizedQueryBuilder) BulkUpdate(ctx context.Context, table string, updates map[string]interface{}, whereClause string) error {
	if len(updates) == 0 {
		return fmt.Errorf("no updates provided")
	}

	setClauses := make([]string, 0, len(updates))
	args := make([]interface{}, 0)
	paramIdx := 1

	for col, val := range updates {
		setClauses = append(setClauses, fmt.Sprintf("%s = $%d", col, paramIdx))
		args = append(args, val)
		paramIdx++
	}

	query := fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s",
		table,
		strings.Join(setClauses, ", "),
		whereClause,
	)

	_, err := oqb.db.ExecContext(ctx, query, args...)
	return err
}

// JoinQuery builds optimized join queries
func (oqb *OptimizedQueryBuilder) JoinQuery(ctx context.Context, mainTable string, joinTables map[string]string, selectColumns []string, whereClause string, args ...interface{}) (*sql.Rows, error) {
	// mainTable: "users u"
	// joinTables: {"itineraries i": "ON u.id = i.user_id"}
	// selectColumns: []string{"u.id", "u.name", "i.title", "i.status"}

	joinClauses := make([]string, 0)
	for table, condition := range joinTables {
		joinClauses = append(joinClauses, fmt.Sprintf("JOIN %s %s", table, condition))
	}

	query := fmt.Sprintf(
		"SELECT %s FROM %s %s WHERE %s",
		strings.Join(selectColumns, ", "),
		mainTable,
		strings.Join(joinClauses, " "),
		whereClause,
	)

	return oqb.db.QueryContext(ctx, query, args...)
}

// CachedCount returns a count query with caching recommendation
func (oqb *OptimizedQueryBuilder) CachedCount(ctx context.Context, table string, whereClause string, args ...interface{}) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", table, whereClause)

	var count int64
	err := oqb.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

// QueryOptimizationPattern represents common query patterns
type QueryOptimizationPattern struct {
	Name        string
	Description string
	Query       string
	Performance string
	ClickClause string // Clause to use for optimization
}

// CommonPatterns returns common optimization patterns
func CommonPatterns() []QueryOptimizationPattern {
	return []QueryOptimizationPattern{
		{
			Name:        "Select Only Required Columns",
			Description: "Avoid SELECT *",
			Query:       "SELECT col1, col2 FROM table WHERE ...",
			Performance: "Reduces data transfer and memory usage",
			ClickClause: "SELECT specific_columns INSTEAD OF *",
		},
		{
			Name:        "Use LIMIT for Large Results",
			Description: "Paginate large result sets",
			Query:       "SELECT * FROM table LIMIT 100 OFFSET 200",
			Performance: "Reduces memory and network overhead",
			ClickClause: "LIMIT + OFFSET FOR PAGINATION",
		},
		{
			Name:        "Index on WHERE Clauses",
			Description: "Index columns used in WHERE clauses",
			Query:       "SELECT * FROM table WHERE indexed_col = ?",
			Performance: "O(log n) instead of O(n) scan",
			ClickClause: "CREATE INDEX idx_col ON table(column)",
		},
		{
			Name:        "Use JOIN Instead of Subqueries",
			Description: "Joins are usually more efficient",
			Query:       "SELECT * FROM t1 JOIN t2 ON t1.id = t2.t1_id",
			Performance: "Better query planner optimization",
			ClickClause: "JOIN INSTEAD OF SUBQUERY",
		},
		{
			Name:        "Batch Operations",
			Description: "Combine multiple operations",
			Query:       "INSERT INTO table VALUES (...), (...), (...)",
			Performance: "Reduces round trips and overhead",
			ClickClause: "BATCH INSERT/UPDATE/DELETE",
		},
	}
}

// QueryCache stores query results
type QueryCache struct {
	cache map[string][]interface{}
}

// NewQueryCache creates a new query cache
func NewQueryCache() *QueryCache {
	return &QueryCache{
		cache: make(map[string][]interface{}),
	}
}

// CacheKey generates a cache key for a query
func CacheKey(query string, args ...interface{}) string {
	return fmt.Sprintf("%s:%v", query, args)
}

// Set caches a query result
func (qc *QueryCache) Set(key string, result []interface{}) {
	qc.cache[key] = result
}

// Get retrieves a cached result
func (qc *QueryCache) Get(key string) ([]interface{}, bool) {
	result, exists := qc.cache[key]
	return result, exists
}

// Clear clears the cache
func (qc *QueryCache) Clear() {
	qc.cache = make(map[string][]interface{})
}

// QueryPlan represents a database query plan
type QueryPlan struct {
	Query         string
	Plan          string
	Cost          float64
	EstimatedRows int64
	ActualRows    int64
	ExecutionTime float64
}

// QueryPlanner plans query execution
type QueryPlanner struct {
	db *sql.DB
}

// NewQueryPlanner creates a new query planner
func NewQueryPlanner(db *sql.DB) *QueryPlanner {
	return &QueryPlanner{db: db}
}

// PlanQuery explains a query plan
func (qp *QueryPlanner) PlanQuery(ctx context.Context, query string, args ...interface{}) (*QueryPlan, error) {
	explainQuery := fmt.Sprintf("EXPLAIN (FORMAT JSON) %s", query)

	rows, err := qp.db.QueryContext(ctx, explainQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to plan query: %w", err)
	}
	defer rows.Close()

	var jsonPlan string
	if rows.Next() {
		if err := rows.Scan(&jsonPlan); err != nil {
			return nil, err
		}
	}

	return &QueryPlan{
		Query: query,
		Plan:  jsonPlan,
	}, nil
}

// Recommendations represents optimization recommendations
type Recommendations struct {
	HighPriority   []string // Critical optimizations
	MediumPriority []string // Important optimizations
	LowPriority    []string // Nice-to-have optimizations
}

// AnalyzeQueryPattern analyzes a query for optimization opportunities
func AnalyzeQueryPattern(query string) *Recommendations {
	recs := &Recommendations{
		HighPriority:   []string{},
		MediumPriority: []string{},
		LowPriority:    []string{},
	}

	query = strings.ToUpper(query)

	// Check for SELECT *
	if strings.Contains(query, "SELECT *") {
		recs.HighPriority = append(recs.HighPriority, "Avoid SELECT * - specify only needed columns")
	}

	// Check for missing WHERE clause (except for specific cases)
	if !strings.Contains(query, "WHERE") && !strings.Contains(query, "JOIN") {
		recs.MediumPriority = append(recs.MediumPriority, "Consider adding WHERE clause to limit result set")
	}

	// Check for potential N+1 query pattern
	if strings.Contains(query, "SELECT") && strings.Count(query, "SELECT") > 1 {
		recs.MediumPriority = append(recs.MediumPriority, "Potential N+1 query pattern - use JOIN instead")
	}

	// Check for LIKE with wildcard at start
	if strings.Contains(query, "LIKE '%") {
		recs.HighPriority = append(recs.HighPriority, "LIKE '%pattern' cannot use indexes efficiently - consider full-text search")
	}

	// Check for complex calculations in WHERE
	if strings.Contains(query, "WHERE") && (strings.Contains(query, "+") || strings.Contains(query, "*")) {
		recs.MediumPriority = append(recs.MediumPriority, "Calculations in WHERE clause prevent index usage")
	}

	// Check for DISTINCT on large tables
	if strings.Contains(query, "DISTINCT") {
		recs.LowPriority = append(recs.LowPriority, "DISTINCT requires extra processing - ensure it's necessary")
	}

	return recs
}
