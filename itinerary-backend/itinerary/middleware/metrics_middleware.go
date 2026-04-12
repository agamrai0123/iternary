package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// MetricsMiddleware is middleware for recording HTTP metrics
type MetricsMiddleware struct {
	metrics *Metrics
	logger  *Logger
}

// NewMetricsMiddleware creates a new metrics middleware
func NewMetricsMiddleware(metrics *Metrics, logger *Logger) *MetricsMiddleware {
	return &MetricsMiddleware{
		metrics: metrics,
		logger:  logger,
	}
}

// MetricsHandler is the gin middleware handler
func (mm *MetricsMiddleware) MetricsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Record request size
		contentLength := int64(0)
		if c.Request.ContentLength > 0 {
			contentLength = c.Request.ContentLength
		}

		// Create a response writer wrapper to capture response size
		responseWriter := &responseWriterWrapper{
			ResponseWriter: c.Writer,
			size:           0,
		}
		c.Writer = responseWriter

		// Process request
		c.Next()

		// Calculate metrics
		latency := time.Since(startTime).Seconds() * 1000 // Convert to ms
		statusCode := c.Writer.Status()

		// Record metrics
		path := c.Request.URL.Path
		method := c.Request.Method
		mm.metrics.RecordHTTPRequest(method, path, statusCode, latency, contentLength, int64(responseWriter.size))

		// Create log fields
		logFields := []interface{}{
			"method", method,
			"path", path,
			"status", statusCode,
			"latency_ms", latency,
			"request_size_bytes", contentLength,
			"response_size_bytes", responseWriter.size,
			"client_ip", c.ClientIP(),
		}

		// Log at appropriate level
		if statusCode >= http.StatusInternalServerError {
			mm.logger.Error("http_request_error", logFields...)
		} else if statusCode >= http.StatusBadRequest {
			mm.logger.Warn("http_request_invalid", logFields...)
		} else {
			mm.logger.Info("http_request_success", logFields...)
		}
	}
}

// responseWriterWrapper wraps gin.ResponseWriter to capture response size
type responseWriterWrapper struct {
	gin.ResponseWriter
	size int
}

func (w *responseWriterWrapper) Write(b []byte) (int, error) {
	w.size += len(b)
	return w.ResponseWriter.Write(b)
}

func (w *responseWriterWrapper) WriteString(s string) (int, error) {
	w.size += len(s)
	return w.ResponseWriter.WriteString(s)
}

// ErrorHandlerMiddleware handles errors with proper logging and metrics
func (mm *MetricsMiddleware) ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check for errors
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				if validationErr, ok := err.Err.(*APIError); ok {
					mm.logger.Error("api_error",
						"code", validationErr.Code,
						"message", validationErr.Message,
						"details", validationErr.Details,
						"status", validationErr.StatusCode,
						"path", c.Request.URL.Path,
						"method", c.Request.Method,
					)

					// Record metrics based on error type
					switch validationErr.Code {
					case ErrValidationError:
						mm.metrics.RecordValidationError()
					case ErrUnauthorized:
						mm.metrics.mu.Lock()
						mm.metrics.AuthorizationErrors++
						mm.metrics.mu.Unlock()
					}
				} else {
					mm.logger.Error("request_error",
						"error", err.Error(),
						"type", err.Type,
						"path", c.Request.URL.Path,
					)
				}
			}
		}
	}
}

// DatabaseMetricsMiddleware wraps database operations with metrics
func (mm *MetricsMiddleware) RecordDatabaseQuery(operation string, fn func() error) error {
	startTime := time.Now()
	err := fn()
	latency := time.Since(startTime).Seconds() * 1000 // Convert to ms

	mm.metrics.RecordDatabaseQuery(operation, latency, err)

	if err != nil {
		mm.logger.Error("database_query_error",
			"operation", operation,
			"latency_ms", latency,
			"error", err.Error(),
		)
	} else {
		mm.logger.Debug("database_query_success",
			"operation", operation,
			"latency_ms", latency,
		)
	}

	return err
}

// PanicRecoveryMiddleware recovers from panics and logs them
func (mm *MetricsMiddleware) PanicRecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				mm.metrics.RecordPanicRecovery()

				mm.logger.Error("panic_recovered",
					"error", err,
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
					"client_ip", c.ClientIP(),
				)

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": map[string]interface{}{
						"code":    ErrInternalServer,
						"message": "Internal server error",
					},
				})
			}
		}()
		c.Next()
	}
}

// MetricsEndpoint returns metrics in Prometheus-compatible format
func (mm *MetricsMiddleware) MetricsEndpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		snapshot := mm.metrics.GetMetricsSnapshot()
		c.JSON(http.StatusOK, snapshot)
	}
}

// HealthCheckMiddleware for detailed health check
func (mm *MetricsMiddleware) HealthCheckEndpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":     "healthy",
			"timestamp":  time.Now(),
			"uptime_sec": time.Since(mm.metrics.StartTime).Seconds(),
		})
	}
}
