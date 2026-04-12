package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

// IndexManager manages database indexes for query optimization
type IndexManager struct {
	db *sql.DB
}

// Index represents a database index
type Index struct {
	Name      string
	TableName string
	Columns   []string
	IsUnique  bool
	IsPartial bool
	Condition string
}

// NewIndexManager creates a new index manager
func NewIndexManager(db *sql.DB) *IndexManager {
	return &IndexManager{db: db}
}

// CreateIndex creates a new database index
func (im *IndexManager) CreateIndex(ctx context.Context, index *Index) error {
	if index.Name == "" || index.TableName == "" || len(index.Columns) == 0 {
		return fmt.Errorf("invalid index configuration: name, table, and columns required")
	}

	query := fmt.Sprintf("CREATE INDEX IF NOT EXISTS %s ON %s (", index.Name, index.TableName)
	for i, col := range index.Columns {
		if i > 0 {
			query += ", "
		}
		query += col
	}
	query += ")"

	if index.IsPartial && index.Condition != "" {
		query += " WHERE " + index.Condition
	}

	_, err := im.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	log.Printf("Created index: %s on %s (%v)", index.Name, index.TableName, index.Columns)
	return nil
}

// CreateCompositeIndex creates a multi-column index
func (im *IndexManager) CreateCompositeIndex(ctx context.Context, name string, tableName string, columns ...string) error {
	return im.CreateIndex(ctx, &Index{
		Name:      name,
		TableName: tableName,
		Columns:   columns,
	})
}

// CreateUniqueIndex creates a unique index
func (im *IndexManager) CreateUniqueIndex(ctx context.Context, name string, tableName string, columns ...string) error {
	return im.CreateIndex(ctx, &Index{
		Name:      name,
		TableName: tableName,
		Columns:   columns,
		IsUnique:  true,
	})
}

// CreatePartialIndex creates a partial index with a condition
func (im *IndexManager) CreatePartialIndex(ctx context.Context, name string, tableName string, condition string, columns ...string) error {
	return im.CreateIndex(ctx, &Index{
		Name:      name,
		TableName: tableName,
		Columns:   columns,
		IsPartial: true,
		Condition: condition,
	})
}

// DropIndex removes an index
func (im *IndexManager) DropIndex(ctx context.Context, indexName string) error {
	query := fmt.Sprintf("DROP INDEX IF EXISTS %s", indexName)
	_, err := im.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to drop index: %w", err)
	}

	log.Printf("Dropped index: %s", indexName)
	return nil
}

// ListIndexes retrieves all indexes for a table
func (im *IndexManager) ListIndexes(ctx context.Context, tableName string) ([]string, error) {
	query := `
		SELECT indexname 
		FROM pg_indexes 
		WHERE tablename = $1 
		AND indexname NOT LIKE 'pg_%'
	`

	rows, err := im.db.QueryContext(ctx, query, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to list indexes: %w", err)
	}
	defer rows.Close()

	var indexes []string
	for rows.Next() {
		var indexName string
		if err := rows.Scan(&indexName); err != nil {
			return nil, fmt.Errorf("failed to scan index: %w", err)
		}
		indexes = append(indexes, indexName)
	}

	return indexes, rows.Err()
}

// IndexStatistics represents index usage statistics
type IndexStatistics struct {
	IndexName     string
	TableName     string
	IndexSize     int64
	IndexScans    int64
	IndexSegScans int64
	IndexTuples   int64
}

// GetIndexStatistics retrieves statistics for an index
func (im *IndexManager) GetIndexStatistics(ctx context.Context, tableName string) ([]IndexStatistics, error) {
	query := `
		SELECT 
			i.indexname,
			i.tablename,
			pg_relation_size(i.indexrelid) AS index_size,
			COALESCE(s.idx_scan, 0) AS index_scans,
			COALESCE(s.idx_tup_read, 0) AS index_tuples
		FROM pg_indexes i
		LEFT JOIN pg_stat_user_indexes s ON i.indexname = s.indexrelname
		WHERE i.tablename = $1
		AND i.indexname NOT LIKE 'pg_%'
		ORDER BY pg_relation_size(i.indexrelid) DESC
	`

	rows, err := im.db.QueryContext(ctx, query, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get index statistics: %w", err)
	}
	defer rows.Close()

	var stats []IndexStatistics
	for rows.Next() {
		var stat IndexStatistics
		if err := rows.Scan(&stat.IndexName, &stat.TableName, &stat.IndexSize, &stat.IndexScans, &stat.IndexTuples); err != nil {
			return nil, fmt.Errorf("failed to scan statistics: %w", err)
		}
		stats = append(stats, stat)
	}

	return stats, rows.Err()
}

// UnusedIndexes identifies indexes that are not being used
func (im *IndexManager) UnusedIndexes(ctx context.Context) ([]string, error) {
	query := `
		SELECT i.indexname
		FROM pg_indexes i
		LEFT JOIN pg_stat_user_indexes s ON i.indexname = s.indexrelname
		WHERE (s.idx_scan = 0 OR s.idx_scan IS NULL)
		AND i.indexname NOT LIKE 'pg_%'
		AND i.indexname NOT LIKE '%_pkey'
	`

	rows, err := im.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get unused indexes: %w", err)
	}
	defer rows.Close()

	var indexes []string
	for rows.Next() {
		var indexName string
		if err := rows.Scan(&indexName); err != nil {
			return nil, fmt.Errorf("failed to scan index: %w", err)
		}
		indexes = append(indexes, indexName)
	}

	return indexes, rows.Err()
}

// IndexRecommendations represents index recommendations
type IndexRecommendation struct {
	TableName     string
	Columns       []string
	Benefit       string
	ExistingUsage int64
}

// AnalyzeSlowQueries suggests indexes for slow queries
func (im *IndexManager) AnalyzeSlowQueries(ctx context.Context) ([]IndexRecommendation, error) {
	// Query pg_stat_statements for slow queries
	query := `
		SELECT 
			t.tablename,
			a.attname,
			s.calls,
			s.mean_exec_time
		FROM pg_stat_statements s
		JOIN pg_class c ON s.query LIKE '%' || c.relname || '%'
		JOIN pg_tables t ON c.relname = t.tablename
		JOIN pg_attribute a ON a.attrelid = c.oid
		WHERE s.mean_exec_time > 1000 -- queries taking > 1 second
		AND s.calls > 10
		GROUP BY t.tablename, a.attname, s.calls, s.mean_exec_time
		ORDER BY s.mean_exec_time DESC
		LIMIT 10
	`

	rows, err := im.db.QueryContext(ctx, query)
	if err != nil {
		// pg_stat_statements might not be available
		return nil, fmt.Errorf("failed to analyze slow queries: %w", err)
	}
	defer rows.Close()

	var recommendations []IndexRecommendation
	for rows.Next() {
		var tableName, columnName string
		var calls int64
		var meanTime float64

		if err := rows.Scan(&tableName, &columnName, &calls, &meanTime); err != nil {
			continue
		}

		recommendations = append(recommendations, IndexRecommendation{
			TableName:     tableName,
			Columns:       []string{columnName},
			Benefit:       fmt.Sprintf("Average time: %.2fms", meanTime),
			ExistingUsage: calls,
		})
	}

	return recommendations, rows.Err()
}

// InitializeOptimalIndexes creates a set of optimized indexes for the schema
func (im *IndexManager) InitializeOptimalIndexes(ctx context.Context) error {
	// Define optimal indexes for Itinerary schema
	optimalIndexes := []Index{
		// Users table indexes
		{Name: "idx_users_email", TableName: "users", Columns: []string{"email"}, IsUnique: true},
		{Name: "idx_users_created", TableName: "users", Columns: []string{"created_at"}},

		// Itineraries table indexes
		{Name: "idx_itineraries_user_id", TableName: "itineraries", Columns: []string{"user_id"}},
		{Name: "idx_itineraries_status", TableName: "itineraries", Columns: []string{"status"}},
		{Name: "idx_itineraries_user_status", TableName: "itineraries", Columns: []string{"user_id", "status"}},

		// Destinations table indexes
		{Name: "idx_destinations_itinerary_id", TableName: "destinations", Columns: []string{"itinerary_id"}},
		{Name: "idx_destinations_location", TableName: "destinations", Columns: []string{"location"}},

		// Activities table indexes
		{Name: "idx_activities_destination_id", TableName: "activities", Columns: []string{"destination_id"}},
		{Name: "idx_activities_date", TableName: "activities", Columns: []string{"date"}},

		// Flights table indexes
		{Name: "idx_flights_itinerary_id", TableName: "flights", Columns: []string{"itinerary_id"}},
		{Name: "idx_flights_departure_date", TableName: "flights", Columns: []string{"departure_date"}},

		// Hotels table indexes
		{Name: "idx_hotels_destination_id", TableName: "hotels", Columns: []string{"destination_id"}},
		{Name: "idx_hotels_check_in", TableName: "hotels", Columns: []string{"check_in_date"}},
	}

	for _, idx := range optimalIndexes {
		if err := im.CreateIndex(ctx, &idx); err != nil {
			log.Printf("Warning: Failed to create index %s: %v", idx.Name, err)
		}
	}

	return nil
}
