package handlers

import (
	"context"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthCheckResponse represents the health check response
type HealthCheckResponse struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Uptime    string                 `json:"uptime"`
	Services  map[string]interface{} `json:"services"`
	Version   string                 `json:"version"`
}

// HealthCheckHandler handles /health endpoint
func (h *Handlers) HealthCheckHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check database connectivity
	dbStatus := "healthy"
	var dbLatency time.Duration

	if h.service.db != nil && h.service.db.conn != nil {
		start := time.Now()
		if err := h.service.db.conn.PingContext(ctx); err != nil {
			dbStatus = "unhealthy"
			h.logger.Error("Database health check failed: " + err.Error())
		}
		dbLatency = time.Since(start)
	} else {
		dbStatus = "unavailable"
	}

	// Prepare response
	response := HealthCheckResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Uptime:    "0h0m0s", // Could track actual uptime
		Version:   "1.0.0",
		Services: map[string]interface{}{
			"database": map[string]interface{}{
				"status":  dbStatus,
				"latency": dbLatency.String(),
			},
			"cache": map[string]interface{}{
				"status": "operational",
			},
		},
	}

	// If any service is unhealthy, mark overall status as degraded
	if dbStatus != "healthy" {
		response.Status = "degraded"
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	c.JSON(http.StatusOK, response)
}

// ReadinessHandler handles /ready endpoint for Kubernetes readiness probe
func (h *Handlers) ReadinessHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Check if database is ready
	if h.service.db == nil || h.service.db.conn == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"ready":  false,
			"reason": "database not initialized",
		})
		return
	}

	if err := h.service.db.conn.PingContext(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"ready":  false,
			"reason": "database not responding",
			"error":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ready":     true,
		"timestamp": time.Now(),
	})
}

// LivenessHandler handles /live endpoint for Kubernetes liveness probe
func (h *Handlers) LivenessHandler(c *gin.Context) {
	// Liveness probe just checks if the application is running
	// It should be simple and not make external calls
	c.JSON(http.StatusOK, gin.H{
		"alive":     true,
		"timestamp": time.Now(),
	})
}

// MetricsHandler exposes Prometheus-style metrics
func (h *Handlers) MetricsHandler(c *gin.Context) {
	if h.metrics == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "metrics not initialized",
		})
		return
	}

	// Return metrics in Prometheus format
	h.metrics.mu.RLock()
	totalRequests := int64(0)
	totalErrors := int64(0)
	for _, v := range h.metrics.HTTPRequestsTotal {
		totalRequests += v
	}
	for _, v := range h.metrics.HTTPErrorsTotal {
		totalErrors += v
	}
	h.metrics.mu.RUnlock()

	metrics := map[string]interface{}{
		"http_requests_total":    totalRequests,
		"http_errors_total":      totalErrors,
		"database_queries_total": h.metrics.DatabaseQueriesTotal,
		"destinations_created":   atomic.LoadInt64(&h.metrics.DestinationsCreated),
		"itineraries_created":    atomic.LoadInt64(&h.metrics.ItinerariesCreated),
		"comments_created":       atomic.LoadInt64(&h.metrics.CommentsCreated),
		"likes_total":            atomic.LoadInt64(&h.metrics.LikesTotal),
	}

	c.JSON(http.StatusOK, metrics)
}

// StatusHandler returns detailed service status
func (h *Handlers) StatusHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	status := map[string]interface{}{
		"service":     "itinerary-backend",
		"version":     "1.0.0",
		"status":      "running",
		"timestamp":   time.Now().UTC(),
		"environment": gin.Mode(),
		"diagnostics": map[string]interface{}{
			"database": h.checkDatabaseStatus(ctx),
			"cache":    h.checkCacheStatus(ctx),
		},
	}

	c.JSON(http.StatusOK, status)
}

// checkDatabaseStatus checks database health
func (h *Handlers) checkDatabaseStatus(ctx context.Context) map[string]interface{} {
	result := map[string]interface{}{
		"status":  "unknown",
		"details": map[string]interface{}{},
	}

	if h.service.db == nil || h.service.db.conn == nil {
		result["status"] = "not initialized"
		return result
	}

	// Check connection
	if err := h.service.db.conn.PingContext(ctx); err != nil {
		result["status"] = "unreachable"
		result["error"] = err.Error()
		return result
	}

	result["status"] = "healthy"

	// Get connection pool stats
	stats := h.service.db.conn.Stats()
	result["details"] = map[string]interface{}{
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration.String(),
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
		"max_open_connections": stats.MaxOpenConnections,
	}

	return result
}

// checkCacheStatus checks cache health
func (h *Handlers) checkCacheStatus(ctx context.Context) map[string]interface{} {
	result := map[string]interface{}{
		"status": "not configured",
	}

	// Cache implementation would check Redis or in-memory cache here
	// For now, return placeholder
	result["status"] = "operational"

	return result
}

// RegisterHealthRoutes registers all health check routes
func RegisterHealthRoutes(router *gin.Engine, handlers *Handlers) {
	// Health checks (no auth required)
	router.GET("/health", handlers.HealthCheckHandler) // Kubernetes liveness probe
	router.GET("/ready", handlers.ReadinessHandler)    // Kubernetes readiness probe
	router.GET("/live", handlers.LivenessHandler)      // Kubernetes startup probe
	router.GET("/status", handlers.StatusHandler)      // Detailed status
	router.GET("/metrics", handlers.MetricsHandler)    // Prometheus metrics

	// Health check aliases for compatibility
	router.GET("/healthz", handlers.HealthCheckHandler)
	router.GET("/readyz", handlers.ReadinessHandler)
	router.GET("/livez", handlers.LivenessHandler)
}
