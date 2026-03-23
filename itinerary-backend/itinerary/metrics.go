package itinerary

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// Metrics holds application metrics in Prometheus style
type Metrics struct {
	mu sync.RWMutex

	// HTTP Request Metrics
	HTTPRequestsTotal     map[string]int64     // {method_path}
	HTTPRequestDuration   map[string][]float64 // {method_path} - latencies in ms
	HTTPErrorsTotal       map[string]int64     // {method_status_code}
	HTTPResponseSizeTotal map[string]int64     // {method_path} - response sizes in bytes
	HTTPRequestSizeTotal  map[string]int64     // {method_path} - request sizes in bytes

	// Business Logic Metrics
	DestinationsCreated int64
	ItinerariesCreated  int64
	CommentsCreated     int64
	LikesTotal          int64
	SearchQueries       int64

	// Database Metrics
	DatabaseQueriesTotal  map[string]int64     // {operation}
	DatabaseQueryErrors   map[string]int64     // {operation}
	DatabaseQueryDuration map[string][]float64 // {operation} - durations in ms

	// System Metrics
	StartTime          time.Time
	CurrentConnections int64
	PeakConnections    int64
	GoroutineCount     int

	// Error Metrics
	PanicRecoveries     int64
	ValidationErrors    int64
	AuthorizationErrors int64
}

// NewMetrics creates a new metrics instance
func NewMetrics() *Metrics {
	return &Metrics{
		HTTPRequestsTotal:     make(map[string]int64),
		HTTPErrorsTotal:       make(map[string]int64),
		HTTPResponseSizeTotal: make(map[string]int64),
		HTTPRequestSizeTotal:  make(map[string]int64),
		HTTPRequestDuration:   make(map[string][]float64),

		DatabaseQueriesTotal:  make(map[string]int64),
		DatabaseQueryErrors:   make(map[string]int64),
		DatabaseQueryDuration: make(map[string][]float64),

		StartTime: time.Now(),
	}
}

// RecordHTTPRequest records an HTTP request metric
func (m *Metrics) RecordHTTPRequest(method, path string, statusCode int, latencyMs float64, reqSize, respSize int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := fmt.Sprintf("%s_%s", method, path)
	m.HTTPRequestsTotal[key]++
	m.HTTPRequestDuration[key] = append(m.HTTPRequestDuration[key], latencyMs)
	m.HTTPResponseSizeTotal[key] += respSize
	m.HTTPRequestSizeTotal[key] += reqSize

	if statusCode >= 400 {
		errorKey := fmt.Sprintf("%s_%s_%d", method, path, statusCode)
		m.HTTPErrorsTotal[errorKey]++
	}
}

// RecordDatabaseQuery records a database query metric
func (m *Metrics) RecordDatabaseQuery(operation string, latencyMs float64, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.DatabaseQueriesTotal[operation]++
	m.DatabaseQueryDuration[operation] = append(m.DatabaseQueryDuration[operation], latencyMs)

	if err != nil {
		m.DatabaseQueryErrors[operation]++
	}
}

// RecordDestinationCreated increments destination creation counter
func (m *Metrics) RecordDestinationCreated() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.DestinationsCreated++
}

// RecordItineraryCreated increments itinerary creation counter
func (m *Metrics) RecordItineraryCreated() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ItinerariesCreated++
}

// RecordCommentCreated increments comment creation counter
func (m *Metrics) RecordCommentCreated() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.CommentsCreated++
}

// RecordLike increments like counter
func (m *Metrics) RecordLike() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.LikesTotal++
}

// RecordSearchQuery increments search query counter
func (m *Metrics) RecordSearchQuery() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.SearchQueries++
}

// RecordValidationError increments validation error counter
func (m *Metrics) RecordValidationError() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ValidationErrors++
}

// RecordPanicRecovery increments panic recovery counter
func (m *Metrics) RecordPanicRecovery() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.PanicRecoveries++
}

// UpdateConnections updates the current connection count
func (m *Metrics) UpdateConnections(count int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.CurrentConnections = count
	if count > m.PeakConnections {
		m.PeakConnections = count
	}
}

// GetMetricsSnapshot returns a snapshot of current metrics
func (m *Metrics) GetMetricsSnapshot() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Calculate averages
	httpLatencies := make(map[string]float64)
	for key, durations := range m.HTTPRequestDuration {
		if len(durations) > 0 {
			sum := float64(0)
			for _, d := range durations {
				sum += d
			}
			httpLatencies[key] = sum / float64(len(durations))
		}
	}

	dbLatencies := make(map[string]float64)
	for key, durations := range m.DatabaseQueryDuration {
		if len(durations) > 0 {
			sum := float64(0)
			for _, d := range durations {
				sum += d
			}
			dbLatencies[key] = sum / float64(len(durations))
		}
	}

	uptime := time.Since(m.StartTime)

	return map[string]interface{}{
		"http": map[string]interface{}{
			"requests_total":      m.HTTPRequestsTotal,
			"errors_total":        m.HTTPErrorsTotal,
			"average_latency_ms":  httpLatencies,
			"response_size_total": m.HTTPResponseSizeTotal,
			"request_size_total":  m.HTTPRequestSizeTotal,
		},
		"business": map[string]interface{}{
			"destinations_created": m.DestinationsCreated,
			"itineraries_created":  m.ItinerariesCreated,
			"comments_created":     m.CommentsCreated,
			"likes_total":          m.LikesTotal,
			"search_queries":       m.SearchQueries,
		},
		"database": map[string]interface{}{
			"queries_total":      m.DatabaseQueriesTotal,
			"query_errors":       m.DatabaseQueryErrors,
			"average_latency_ms": dbLatencies,
		},
		"errors": map[string]interface{}{
			"panic_recoveries":     m.PanicRecoveries,
			"validation_errors":    m.ValidationErrors,
			"authorization_errors": m.AuthorizationErrors,
		},
		"system": map[string]interface{}{
			"uptime_seconds":        uptime.Seconds(),
			"current_connections":   m.CurrentConnections,
			"peak_connections":      m.PeakConnections,
			"goroutines":            runtime.NumGoroutine(),
			"memory_alloc_mb":       getAllocMemoryMB(),
			"memory_total_alloc_mb": getTotalAllocMemoryMB(),
		},
	}
}

// getAllocMemoryMB returns current memory allocation in MB
func getAllocMemoryMB() float64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return float64(m.Alloc) / 1024 / 1024
}

// getTotalAllocMemoryMB returns total memory allocated in MB
func getTotalAllocMemoryMB() float64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return float64(m.TotalAlloc) / 1024 / 1024
}
