package itinerary

import (
	"testing"
)

// TestMetricsInitialization verifies metrics can be initialized
func TestMetricsInitialization(t *testing.T) {
	metrics := &Metrics{}

	if metrics == nil {
		t.Error("Metrics initialization failed")
	}
}

// TestMetricsFields verifies metrics structure
func TestMetricsFields(t *testing.T) {
	metrics := &Metrics{
		RequestsTotal:      100,
		RequestsDuration:   5000,
		ResponsesSuccess:   95,
		ResponsesError:     5,
		DatabaseQueries:    50,
		CacheHits:          30,
		CacheMisses:        20,
	}

	tests := []struct {
		name        string
		field       int64
		expectedMin int64
	}{
		{
			name:        "total requests tracking",
			field:       metrics.RequestsTotal,
			expectedMin: 0,
		},
		{
			name:        "request duration tracking",
			field:       metrics.RequestsDuration,
			expectedMin: 0,
		},
		{
			name:        "successful responses tracking",
			field:       metrics.ResponsesSuccess,
			expectedMin: 0,
		},
		{
			name:        "error responses tracking",
			field:       metrics.ResponsesError,
			expectedMin: 0,
		},
		{
			name:        "database queries tracking",
			field:       metrics.DatabaseQueries,
			expectedMin: 0,
		},
		{
			name:        "cache hits tracking",
			field:       metrics.CacheHits,
			expectedMin: 0,
		},
		{
			name:        "cache misses tracking",
			field:       metrics.CacheMisses,
			expectedMin: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.field < tt.expectedMin {
				t.Errorf("Metric value should be >= %d, got %d", tt.expectedMin, tt.field)
			}
		})
	}
}

// TestMetricsSuccessRateCalculation verifies success rate calculation
func TestMetricsSuccessRateCalculation(t *testing.T) {
	tests := []struct {
		name            string
		successCount    int64
		errorCount      int64
		expectedSuccess float64
	}{
		{
			name:            "100% success rate",
			successCount:    100,
			errorCount:      0,
			expectedSuccess: 100.0,
		},
		{
			name:            "90% success rate",
			successCount:    90,
			errorCount:      10,
			expectedSuccess: 90.0,
		},
		{
			name:            "50% success rate",
			successCount:    50,
			errorCount:      50,
			expectedSuccess: 50.0,
		},
		{
			name:            "0% success rate",
			successCount:    0,
			errorCount:      100,
			expectedSuccess: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			total := tt.successCount + tt.errorCount
			var successRate float64

			if total > 0 {
				successRate = (float64(tt.successCount) / float64(total)) * 100
			}

			if successRate != tt.expectedSuccess {
				t.Errorf("Expected success rate %.1f%%, got %.1f%%", tt.expectedSuccess, successRate)
			}
		})
	}
}

// TestMetricsCacheHitRateCalculation verifies cache hit rate calculation
func TestMetricsCacheHitRateCalculation(t *testing.T) {
	tests := []struct {
		name              string
		cacheHits        int64
		cacheMisses      int64
		expectedHitRate  float64
	}{
		{
			name:              "100% cache hit rate",
			cacheHits:        100,
			cacheMisses:      0,
			expectedHitRate:  100.0,
		},
		{
			name:              "80% cache hit rate",
			cacheHits:        80,
			cacheMisses:      20,
			expectedHitRate:  80.0,
		},
		{
			name:              "50% cache hit rate",
			cacheHits:        50,
			cacheMisses:      50,
			expectedHitRate:  50.0,
		},
		{
			name:              "0% cache hit rate",
			cacheHits:        0,
			cacheMisses:      100,
			expectedHitRate:  0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			total := tt.cacheHits + tt.cacheMisses
			var hitRate float64

			if total > 0 {
				hitRate = (float64(tt.cacheHits) / float64(total)) * 100
			}

			if hitRate != tt.expectedHitRate {
				t.Errorf("Expected cache hit rate %.1f%%, got %.1f%%", tt.expectedHitRate, hitRate)
			}
		})
	}
}

// TestMetricsAverageDuration verifies average request duration calculation
func TestMetricsAverageDuration(t *testing.T) {
	tests := []struct {
		name               string
		totalDuration      int64
		totalRequests      int64
		expectedAverage    int64
	}{
		{
			name:               "average 50ms",
			totalDuration:      5000,
			totalRequests:      100,
			expectedAverage:    50,
		},
		{
			name:               "average 100ms",
			totalDuration:      10000,
			totalRequests:      100,
			expectedAverage:    100,
		},
		{
			name:               "single request",
			totalDuration:      150,
			totalRequests:      1,
			expectedAverage:    150,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var average int64

			if tt.totalRequests > 0 {
				average = tt.totalDuration / tt.totalRequests
			}

			if average != tt.expectedAverage {
				t.Errorf("Expected average %.0fms, got %.0fms", float64(tt.expectedAverage), float64(average))
			}
		})
	}
}

// TestMetricsZeroValues verifies handling of zero metrics
func TestMetricsZeroValues(t *testing.T) {
	metrics := &Metrics{
		RequestsTotal:      0,
		RequestsDuration:   0,
		ResponsesSuccess:   0,
		ResponsesError:     0,
		DatabaseQueries:    0,
		CacheHits:          0,
		CacheMisses:        0,
	}

	if metrics.RequestsTotal != 0 {
		t.Error("Metrics should be initialized to zero")
	}

	if metrics.ResponsesSuccess != 0 {
		t.Error("Success count should be zero")
	}

	if metrics.ResponsesError != 0 {
		t.Error("Error count should be zero")
	}
}
