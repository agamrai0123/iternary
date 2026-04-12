package itinerary

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// Logger wraps zerolog logger with application-specific methods
type Logger struct {
	log zerolog.Logger
}

// NewLogger creates a new structured logger using zerolog
func NewLogger(config *Config) *Logger {
	// Create log directory if it doesn't exist
	logDir := "log"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.Mkdir(logDir, 0755)
	}

	// Set log level
	levelStr := config.Logging.Level
	level, err := zerolog.ParseLevel(levelStr)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// Create log file
	logFile := fmt.Sprintf("%s/itinerary-%s.log", logDir, time.Now().Format("2006-01-02"))
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("Cannot open log file: %v", err))
	}

	// Create logger with timestamp context
	z := zerolog.New(file).With().Timestamp().Logger()

	return &Logger{log: z}
}

// Info logs an info message with optional fields
func (l *Logger) Info(message string, fields ...interface{}) {
	ctx := l.log.Info()
	for i := 0; i < len(fields)-1; i += 2 {
		key := fmt.Sprintf("%v", fields[i])
		ctx = ctx.Interface(key, fields[i+1])
	}
	ctx.Msg(message)
}

// Error logs an error message with optional fields
func (l *Logger) Error(message string, fields ...interface{}) {
	ctx := l.log.Error()
	for i := 0; i < len(fields)-1; i += 2 {
		key := fmt.Sprintf("%v", fields[i])
		ctx = ctx.Interface(key, fields[i+1])
	}
	ctx.Msg(message)
}

// Debug logs a debug message with optional fields
func (l *Logger) Debug(message string, fields ...interface{}) {
	ctx := l.log.Debug()
	for i := 0; i < len(fields)-1; i += 2 {
		key := fmt.Sprintf("%v", fields[i])
		ctx = ctx.Interface(key, fields[i+1])
	}
	ctx.Msg(message)
}

// Warn logs a warning message with optional fields
func (l *Logger) Warn(message string, fields ...interface{}) {
	ctx := l.log.Warn()
	for i := 0; i < len(fields)-1; i += 2 {
		key := fmt.Sprintf("%v", fields[i])
		ctx = ctx.Interface(key, fields[i+1])
	}
	ctx.Msg(message)
}

// RequestLogger is middleware for logging HTTP requests with structured fields
func (l *Logger) RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(startTime)
		latencyMs := float64(latency.Milliseconds())

		// Get status code
		statusCode := c.Writer.Status()

		// Determine log level based on status code
		logLevel := zerolog.InfoLevel
		if statusCode >= http.StatusBadRequest && statusCode < http.StatusInternalServerError {
			logLevel = zerolog.WarnLevel
		} else if statusCode >= http.StatusInternalServerError {
			logLevel = zerolog.ErrorLevel
		}

		// Log request details
		logEvent := l.log.WithLevel(logLevel).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", statusCode).
			Float64("latency_ms", latencyMs).
			Str("client_ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent())

		// Add error if present
		if len(c.Errors) > 0 {
			logEvent = logEvent.Interface("errors", c.Errors.Errors())
		}

		logEvent.Msg("http_request")
	}
}

// ErrorLogger is middleware for centralized error handling
func (l *Logger) ErrorLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Handle errors if any occurred
		if len(c.Errors) > 0 {
			for _, err := range c.Errors.ByType(gin.ErrorTypeBind) {
				l.Error("bind_error", "error", err.Error(), "type", err.Type)
			}
			for _, err := range c.Errors.ByType(gin.ErrorTypePublic) {
				l.Error("public_error", "error", err.Error())
			}
			for _, err := range c.Errors.ByType(gin.ErrorTypePrivate) {
				l.Error("private_error", "error", err.Error())
			}
		}
	}
}

// RecoveryLogger is middleware for panic recovery with logging
func (l *Logger) RecoveryLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				l.Error("panic_recovered",
					"error", fmt.Sprintf("%v", err),
					"method", c.Request.Method,
					"path", c.Request.URL.Path,
					"client_ip", c.ClientIP(),
				)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
			}
		}()
		c.Next()
	}
}

// GetZerolog returns the underlying zerolog logger for direct access
func (l *Logger) GetZerolog() zerolog.Logger {
	return l.log
}
