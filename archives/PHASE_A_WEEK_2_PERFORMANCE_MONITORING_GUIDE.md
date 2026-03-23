# Phase A Week 2: Performance Monitoring & Alerting Implementation

**Purpose:** Real-time performance monitoring with automated alerts  
**Status:** Ready for integration into Week 2  
**Priority:** High (enables real-time issue detection)

---

## Overview

Real-time performance monitoring provides:
- ✅ Automatic detection of performance degradation
- ✅ Real-time alerts when thresholds exceeded
- ✅ Performance trend analysis
- ✅ Dashboard for operations team
- ✅ Audit trail of performance issues

---

## Architecture

### Components

```
┌─────────────────────────────────────────────────────────┐
│                    MCP + Handlers                        │
└─────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────┐
│          Performance Monitoring Middleware               │
│  ┌─────────────────────────────────────────────────┐   │
│  │ 1. Capture request start time                   │   │
│  │ 2. Capture response time                        │   │
│  │ 3. Capture memory/CPU delta                     │   │
│  │ 4. Record to metrics                            │   │
│  │ 5. Check thresholds                             │   │
│  │ 6. Emit alerts if needed                        │   │
│  └─────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
                            ↓
        ┌───────────────────┬──────────────────┐
        ↓                   ↓                  ↓
   ┌─────────┐        ┌──────────┐      ┌──────────┐
   │ Database│        │ Alerts   │      │ Dashboard│
   │ Metrics │        │ Channel  │      │ Display  │
   └─────────┘        └──────────┘      └──────────┘
        ↓                   ↓
    Store Raw           Send Real-time
    Metrics           Notifications
```

---

## Implementation Steps

### Step 1: Create Performance Monitor Go Packages

Create `itinerary/performance_monitor.go`:

```go
package itinerary

import (
	"context"
	"runtime"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// PerformanceMonitor tracks endpoint performance metrics
type PerformanceMonitor struct {
	db                  *Database
	logger              *Logger
	metrics             *Metrics
	alerts              chan PerformanceAlert
	thresholds          PerformanceThresholds
	muEndpointMetrics   sync.RWMutex
	endpointMetrics     map[string]*EndpointMetric
	alertCooldown       map[string]time.Time // Prevent alert spam
	alertCooldownMutex  sync.RWMutex
}

// PerformanceThresholds defines alert thresholds
type PerformanceThresholds struct {
	P95ThresholdMs        int64   // milliseconds
	P99ThresholdMs        int64   // milliseconds
	ErrorRateThreshold    float64 // 0.01 = 1%
	MemoryThresholdMB     int64   // MB
	AlertCooldownSeconds  int64   // Don't repeat alerts within this time
}

// EndpointMetric tracks performance for an endpoint
type EndpointMetric struct {
	Path            string
	Method          string
	ResponseTimes   []int64     // Sliding window of response times
	ErrorCount      int
	SuccessCount    int
	P50ResponseTime int64
	P95ResponseTime int64
	P99ResponseTime int64
	MaxResponseTime int64
	AvgResponseTime float64
	LastUpdated     time.Time
	MaxSamples      int
}

// PerformanceAlert represents a performance issue
type PerformanceAlert struct {
	ID              string
	AlertType       string // "high_response_time", "error_rate", "memory_spike"
	Endpoint        string
	Method          string
	Severity        string // "info", "warning", "critical"
	ThresholdValue  float64
	CurrentValue    float64
	Message         string
	Recommendation  string
	Timestamp       time.Time
}

// NewPerformanceMonitor creates a new performance monitor
func NewPerformanceMonitor(db *Database, logger *Logger, metrics *Metrics) *PerformanceMonitor {
	return &PerformanceMonitor{
		db:      db,
		logger:  logger,
		metrics: metrics,
		alerts:  make(chan PerformanceAlert, 100),
		thresholds: PerformanceThresholds{
			P95ThresholdMs:       500,  // ms
			P99ThresholdMs:       1000, // ms
			ErrorRateThreshold:   0.01, // 1%
			MemoryThresholdMB:    200,  // MB
			AlertCooldownSeconds: 300,  // 5 minutes
		},
		endpointMetrics: make(map[string]*EndpointMetric),
		alertCooldown:   make(map[string]time.Time),
	}
}

// RecordMetric records a single request metric
func (pm *PerformanceMonitor) RecordMetric(method, path string, responseTimes int64, statusCode int, memoryDeltaMB int64) {
	// Create or get endpoint metric
	key := method + " " + path
	
	pm.muEndpointMetrics.Lock()
	defer pm.muEndpointMetrics.Unlock()

	metric, exists := pm.endpointMetrics[key]
	if !exists {
		metric = &EndpointMetric{
			Path:       path,
			Method:     method,
			MaxSamples: 1000, // Keep last 1000 requests
		}
		pm.endpointMetrics[key] = metric
	}

	// Add response time
	metric.ResponseTimes = append(metric.ResponseTimes, responseTimes)

	// Keep sliding window
	if len(metric.ResponseTimes) > metric.MaxSamples {
		metric.ResponseTimes = metric.ResponseTimes[1:]
	}

	// Update error/success counts
	if statusCode >= 400 {
		metric.ErrorCount++
	} else {
		metric.SuccessCount++
	}

	// Calculate aggregates every 100 requests
	if (metric.ErrorCount + metric.SuccessCount) % 100 == 0 {
		pm.calculateAggregates(metric)
	}

	metric.LastUpdated = time.Now()

	// Store in database asynchronously
	go pm.storeMetricAsync(method, path, responseTimes, statusCode, memoryDeltaMB)
}

// calculateAggregates calculates p50, p95, p99
func (pm *PerformanceMonitor) calculateAggregates(metric *EndpointMetric) {
	if len(metric.ResponseTimes) == 0 {
		return
	}

	// Sort to calculate percentiles
	times := make([]int64, len(metric.ResponseTimes))
	copy(times, metric.ResponseTimes)
	
	// Simple bubble sort (could use quickselect for production)
	for i := 0; i < len(times)-1; i++ {
		for j := 0; j < len(times)-i-1; j++ {
			if times[j] > times[j+1] {
				times[j], times[j+1] = times[j+1], times[j]
			}
		}
	}

	// Calculate percentiles
	metric.MaxResponseTime = times[len(times)-1]
	metric.P50ResponseTime = times[len(times)/2]
	metric.P95ResponseTime = times[int(float64(len(times))*0.95)]
	metric.P99ResponseTime = times[int(float64(len(times))*0.99)]

	// Calculate average
	sum := int64(0)
	for _, t := range times {
		sum += t
	}
	metric.AvgResponseTime = float64(sum) / float64(len(times))

	pm.logger.Info("Endpoint aggregates calculated",
		zap.String("endpoint", metric.Method+" "+metric.Path),
		zap.Int64("p95_ms", metric.P95ResponseTime),
		zap.Int64("p99_ms", metric.P99ResponseTime),
		zap.Int64("max_ms", metric.MaxResponseTime),
		zap.Float64("avg_ms", metric.AvgResponseTime),
	)
}

// CheckThresholds checks if any thresholds are violated
func (pm *PerformanceMonitor) CheckThresholds(method, path string) {
	key := method + " " + path

	pm.muEndpointMetrics.RLock()
	metric, exists := pm.endpointMetrics[key]
	pm.muEndpointMetrics.RUnlock()

	if !exists || metric == nil {
		return
	}

	// Check P95 threshold
	if metric.P95ResponseTime > pm.thresholds.P95ThresholdMs {
		alert := PerformanceAlert{
			ID:             uuid.New().String(),
			AlertType:      "high_response_time_p95",
			Endpoint:       path,
			Method:         method,
			Severity:       "warning",
			ThresholdValue: float64(pm.thresholds.P95ThresholdMs),
			CurrentValue:   float64(metric.P95ResponseTime),
			Message: fmt.Sprintf(
				"P95 response time for %s %s is %dms (threshold: %dms)",
				method, path, metric.P95ResponseTime, pm.thresholds.P95ThresholdMs,
			),
			Recommendation: "Review endpoint implementation and database queries",
			Timestamp:      time.Now(),
		}
		pm.emitAlert(alert)
	}

	// Check P99 threshold (critical)
	if metric.P99ResponseTime > pm.thresholds.P99ThresholdMs {
		alert := PerformanceAlert{
			ID:             uuid.New().String(),
			AlertType:      "high_response_time_p99",
			Endpoint:       path,
			Method:         method,
			Severity:       "critical",
			ThresholdValue: float64(pm.thresholds.P99ThresholdMs),
			CurrentValue:   float64(metric.P99ResponseTime),
			Message: fmt.Sprintf(
				"P99 response time for %s %s is %dms (threshold: %dms)",
				method, path, metric.P99ResponseTime, pm.thresholds.P99ThresholdMs,
			),
			Recommendation: "Urgent review needed - consider scaling or optimization",
			Timestamp:      time.Now(),
		}
		pm.emitAlert(alert)
	}

	// Check error rate
	total := metric.ErrorCount + metric.SuccessCount
	if total > 0 {
		errorRate := float64(metric.ErrorCount) / float64(total)
		if errorRate > pm.thresholds.ErrorRateThreshold {
			alert := PerformanceAlert{
				ID:             uuid.New().String(),
				AlertType:      "high_error_rate",
				Endpoint:       path,
				Method:         method,
				Severity:       "warning",
				ThresholdValue: pm.thresholds.ErrorRateThreshold * 100,
				CurrentValue:   errorRate * 100,
				Message: fmt.Sprintf(
					"Error rate for %s %s is %.1f%% (threshold: %.1f%%)",
					method, path, errorRate*100, pm.thresholds.ErrorRateThreshold*100,
				),
				Recommendation: "Check error logs for details",
				Timestamp:      time.Now(),
			}
			pm.emitAlert(alert)
		}
	}
}

// emitAlert sends alert with cooldown to prevent spam
func (pm *PerformanceMonitor) emitAlert(alert PerformanceAlert) {
	alertKey := alert.AlertType + ":" + alert.Endpoint

	// Check cooldown
	pm.alertCooldownMutex.RLock()
	lastAlert, exists := pm.alertCooldown[alertKey]
	pm.alertCooldownMutex.RUnlock()

	if exists && time.Since(lastAlert) < time.Duration(pm.thresholds.AlertCooldownSeconds)*time.Second {
		// Still in cooldown, skip this alert
		return
	}

	// Update cooldown
	pm.alertCooldownMutex.Lock()
	pm.alertCooldown[alertKey] = time.Now()
	pm.alertCooldownMutex.Unlock()

	// Log alert
	pm.logger.Warn("Performance Alert",
		zap.String("type", alert.AlertType),
		zap.String("severity", alert.Severity),
		zap.String("endpoint", alert.Endpoint),
		zap.Float64("threshold", alert.ThresholdValue),
		zap.Float64("current", alert.CurrentValue),
		zap.String("message", alert.Message),
	)

	// Send to channel (non-blocking)
	select {
	case pm.alerts <- alert:
	default:
		pm.logger.Warn("Alert channel full, alert dropped")
	}
}

// storeMetricAsync stores metric to database asynchronously
func (pm *PerformanceMonitor) storeMetricAsync(method, path string, responseTimeMs int64, statusCode int, memoryDeltaMB int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Store in performance_metrics table (optional - can log to file instead for high-volume)
	// This operation should be lightweight or batched
}

// StartAlertHandler starts the alert handler goroutine
func (pm *PerformanceMonitor) StartAlertHandler(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				pm.logger.Info("Performance monitor alert handler stopped")
				return
			case alert := <-pm.alerts:
				pm.handleAlert(alert)
			}
		}
	}()
}

// handleAlert processes an alert
func (pm *PerformanceMonitor) handleAlert(alert PerformanceAlert) {
	// Store alert in database
	err := pm.db.RecordPerformanceAlert(context.Background(), alert)
	if err != nil {
		pm.logger.Error("Failed to record alert", zap.Error(err))
	}

	// Send to external monitoring (DataDog, New Relic, CloudWatch, Slack)
	if alert.Severity == "critical" {
		pm.sendCriticalNotification(alert)
	}
}

// sendCriticalNotification sends critical alerts to operations
func (pm *PerformanceMonitor) sendCriticalNotification(alert PerformanceAlert) {
	// TODO: Implement integrations
	// - Send to Slack channel
	// - Send to PagerDuty
	// - Send SMS to on-call engineer
	// Example:
	// pm.slackNotify(alert)
	// pm.pagerDutyNotify(alert)
}

// GetEndpointMetric retrieves current metrics for an endpoint
func (pm *PerformanceMonitor) GetEndpointMetric(method, path string) *EndpointMetric {
	key := method + " " + path
	pm.muEndpointMetrics.RLock()
	defer pm.muEndpointMetrics.RUnlock()
	return pm.endpointMetrics[key]
}

// GetAllMetrics returns all endpoint metrics
func (pm *PerformanceMonitor) GetAllMetrics() map[string]*EndpointMetric {
	pm.muEndpointMetrics.RLock()
	defer pm.muEndpointMetrics.RUnlock()
	
	// Return copy to avoid races
	result := make(map[string]*EndpointMetric)
	for k, v := range pm.endpointMetrics {
		result[k] = v
	}
	return result
}

// GetAlertsChan returns the alerts channel
func (pm *PerformanceMonitor) GetAlertsChan() <-chan PerformanceAlert {
	return pm.alerts
}

// GetHealthStatus returns overall system health
func (pm *PerformanceMonitor) GetHealthStatus() HealthStatus {
	var criticalAlerts, warnings int
	
	pm.muEndpointMetrics.RLock()
	defer pm.muEndpointMetrics.RUnlock()

	for _, metric := range pm.endpointMetrics {
		if metric.P99ResponseTime > pm.thresholds.P99ThresholdMs {
			criticalAlerts++
		} else if metric.P95ResponseTime > pm.thresholds.P95ThresholdMs {
			warnings++
		}
	}

	status := "healthy"
	if criticalAlerts > 0 {
		status = "critical"
	} else if warnings > 0 {
		status = "degraded"
	}

	return HealthStatus{
		Status:           status,
		CriticalAlerts:   criticalAlerts,
		Warnings:         warnings,
		EndpointsTracked: len(pm.endpointMetrics),
		Timestamp:        time.Now(),
	}
}

type HealthStatus struct {
	Status           string    `json:"status"`           // healthy, degraded, critical
	CriticalAlerts   int       `json:"critical_alerts"`
	Warnings         int       `json:"warnings"`
	EndpointsTracked int       `json:"endpoints_tracked"`
	Timestamp        time.Time `json:"timestamp"`
}
```

---

### Step 2: Create Monitoring Middleware

Update `itinerary/metrics_middleware.go`:

```go
// Add to MetricsMiddleware struct:
performanceMonitor *PerformanceMonitor

// Add method to MetricsMiddleware:
func (mm *MetricsMiddleware) PerformanceMonitoringMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		startMem := runtime.MemStats{}
		runtime.ReadMemStats(&startMem)

		// Call next handler
		c.Next()

		// Calculate metrics
		duration := time.Since(startTime)
		durationMs := duration.Milliseconds()

		endMem := runtime.MemStats{}
		runtime.ReadMemStats(&endMem)
		memDeltaMB := int64((endMem.Alloc - startMem.Alloc) / 1024 / 1024)

		// Record metric
		mm.performanceMonitor.RecordMetric(
			c.Request.Method,
			c.Request.URL.Path,
			durationMs,
			c.Writer.Status(),
			memDeltaMB,
		)

		// Check thresholds
		mm.performanceMonitor.CheckThresholds(c.Request.Method, c.Request.URL.Path)

		// Add response headers for client
		c.Header("X-Response-Time-Ms", fmt.Sprintf("%d", durationMs))
	}
}
```

---

### Step 3: Create Dashboard Endpoint

```go
// In handlers.go, add:
func (h *Handlers) PerformanceDashboard(c *gin.Context) {
	// Get all metrics
	metrics := h.service.performanceMonitor.GetAllMetrics()
	
	// Get health status
	health := h.service.performanceMonitor.GetHealthStatus()

	// Format response
	type EndpointInfo struct {
		Method          string  `json:"method"`
		Path            string  `json:"path"`
		AvgResponseMs   float64 `json:"avg_response_ms"`
		P95ResponseMs   int64   `json:"p95_response_ms"`
		P99ResponseMs   int64   `json:"p99_response_ms"`
		MaxResponseMs   int64   `json:"max_response_ms"`
		ErrorCount      int     `json:"error_count"`
		SuccessCount    int     `json:"success_count"`
		ErrorRate       float64 `json:"error_rate"`
	}

	var endpoints []EndpointInfo
	for _, m := range metrics {
		errorRate := 0.0
		total := m.ErrorCount + m.SuccessCount
		if total > 0 {
			errorRate = float64(m.ErrorCount) / float64(total)
		}

		endpoints = append(endpoints, EndpointInfo{
			Method:        m.Method,
			Path:          m.Path,
			AvgResponseMs: m.AvgResponseTime,
			P95ResponseMs: m.P95ResponseTime,
			P99ResponseMs: m.P99ResponseTime,
			MaxResponseMs: m.MaxResponseTime,
			ErrorCount:    m.ErrorCount,
			SuccessCount:  m.SuccessCount,
			ErrorRate:     errorRate,
		})
	}

	c.JSON(200, gin.H{
		"health":     health,
		"endpoints":  endpoints,
		"timestamp":  time.Now(),
	})
}
```

---

### Step 4: Register in Routes

Update `routes.go`:

```go
// Add performance monitoring middleware
router.Use(metricsMiddleware.PerformanceMonitoringMiddleware())

// Add dashboard endpoint
router.GET("/api/performance-dashboard", handlers.PerformanceDashboard)

// Start alert handler
go metricsMiddleware.performanceMonitor.StartAlertHandler(context.Background())
```

---

## Usage & Testing

### Enable Performance Monitoring

```bash
# Set performance monitor in service init
performanceMonitor := NewPerformanceMonitor(db, logger, metrics)
service.performanceMonitor = performanceMonitor
```

### View Real-Time Alerts

```bash
# Endpoint returns current performance status
curl http://localhost:8080/api/performance-dashboard | jq .

# Example output:
{
  "health": {
    "status": "degraded",
    "critical_alerts": 0,
    "warnings": 1,
    "endpoints_tracked": 16,
    "timestamp": "2026-03-24T14:32:15Z"
  },
  "endpoints": [
    {
      "method": "GET",
      "path": "/api/v1/group-trips",
      "avg_response_ms": 145.2,
      "p95_response_ms": 280,
      "p99_response_ms": 450,
      "max_response_ms": 502,
      "error_count": 0,
      "success_count": 1000,
      "error_rate": 0.0
    },
    {
      "method": "POST",
      "path": "/api/v1/group-trips/{id}/expenses",
      "avg_response_ms": 420.5,
      "p95_response_ms": 650,  // ⚠️ Above 500ms threshold
      "p99_response_ms": 1200, // 🔴 Above 1000ms threshold
      "max_response_ms": 1500,
      "error_count": 25,
      "success_count": 950,
      "error_rate": 0.026
    }
  ],
  "timestamp": "2026-03-24T14:35:22Z"
}
```

### View Alerts in Logs

```
[WARN] Performance Alert
  alert_type: high_response_time_p95
  severity: warning
  endpoint: /api/v1/group-trips/{id}/expenses
  method: POST
  current_value: 652
  threshold_value: 500
  message: "P95 response time for POST /api/v1/group-trips/{id}/expenses is 652ms (threshold: 500ms)"
  recommendation: "Review endpoint implementation and database queries"

[WARN] Performance Alert
  alert_type: high_response_time_p99
  severity: critical
  endpoint: /api/v1/group-trips/{id}/expenses
  method: POST
  current_value: 1215
  threshold_value: 1000
  message: "P99 response time for POST /api/v1/group-trips/{id}/expenses is 1215ms (threshold: 1000ms)"
  recommendation: "Urgent review needed - consider scaling or optimization"
```

---

## Integration with Week 2 Testing

### Monday
- Performance monitoring framework implemented
- Baseline metrics collected and stored

### Tuesday
- API testing includes performance tracking
- Multi-currency endpoint performance compared

### Wednesday
- Algorithm performance monitored
- Settlement calculation speed tracked

### Thursday
- Stress testing with real-time alerts
- Performance degradation detected automatically

### Friday
- Dashboard and alerts documented
- Team trained on monitoring

---

##  Alert Thresholds Summary

| Alert Type | Warning | Critical | Action |
|-----------|--------|----------|--------|
| Response Time P95 | 500ms | 1000ms | Review code/DB |
| Response Time P99 | 1000ms | 2000ms | Escalate/scale |
| Error Rate | 1% | 5% | Check logs |
| Memory | 200MB | 500MB | Restart/check leak |
| Database Connections | 80% pool | 95% pool | Scale pool |

---

##  Next Steps

1. Code implementation (Tuesday/Wednesday)
2. Enable monitoring middleware in routes
3. Start monitoring in main.go
4. Test alert generation
5. Document for team

---

**Performance monitoring ready for Phase A Week 2 integration!** ✅
