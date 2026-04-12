package middleware

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements per-IP rate limiting
type RateLimiter struct {
	requestCounts map[string][]time.Time
	maxRequests   int
	timeWindow    time.Duration
	mu            sync.RWMutex
	logger        Logger
	cleanupTicker *time.Ticker
}

// Logger interface for dependency injection
type Logger interface {
	Error(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Debug(msg string, fields ...interface{})
}

// NewRateLimiter creates a new rate limiter
// maxRequests: max requests allowed in timeWindow (e.g., 5)
// timeWindow: time duration for rate limit (e.g., 1 minute)
func NewRateLimiter(maxRequests int, timeWindow time.Duration, logger Logger) *RateLimiter {
	rl := &RateLimiter{
		requestCounts: make(map[string][]time.Time),
		maxRequests:   maxRequests,
		timeWindow:    timeWindow,
		logger:        logger,
		cleanupTicker: time.NewTicker(5 * time.Minute),
	}

	// Cleanup old entries periodically to prevent memory leak
	go rl.cleanupLoop()

	return rl
}

// Middleware returns the Gin middleware function
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !rl.Allow(c.ClientIP()) {
			rl.logger.Warn("rate_limit_exceeded", "ip", rl.maskIP(c.ClientIP()))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// AllowLogin returns rate limiter middleware specifically for login endpoints
// Uses stricter limits: 5 attempts per minute per IP
func (rl *RateLimiter) AllowLogin() gin.HandlerFunc {
	strictRL := &RateLimiter{
		requestCounts: rl.requestCounts,
		maxRequests:   5,
		timeWindow:    1 * time.Minute,
		logger:        rl.logger,
		mu:            rl.mu,
	}

	return func(c *gin.Context) {
		if !strictRL.Allow(c.ClientIP()) {
			strictRL.logger.Warn("login_rate_limit_exceeded", "ip", rl.maskIP(c.ClientIP()))
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many login attempts. Please try again in 1 minute.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Allow checks if a request from the given IP is allowed
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.timeWindow)

	// Get or create request history for this IP
	if timestamps, exists := rl.requestCounts[ip]; exists {
		// Remove old requests outside the window
		validRequests := make([]time.Time, 0, len(timestamps))
		for _, t := range timestamps {
			if t.After(windowStart) {
				validRequests = append(validRequests, t)
			}
		}

		// Check if at limit
		if len(validRequests) >= rl.maxRequests {
			rl.requestCounts[ip] = validRequests
			return false
		}

		// Add current request
		validRequests = append(validRequests, now)
		rl.requestCounts[ip] = validRequests
	} else {
		// First request from this IP
		rl.requestCounts[ip] = []time.Time{now}
	}

	return true
}

// Reset clears the rate limiter state (useful for testing)
func (rl *RateLimiter) Reset() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.requestCounts = make(map[string][]time.Time)
}

// cleanupLoop periodically removes stale entries to prevent memory leaks
func (rl *RateLimiter) cleanupLoop() {
	defer rl.cleanupTicker.Stop()

	for range rl.cleanupTicker.C {
		rl.cleanup()
	}
}

// cleanup removes expired entries from the rate limiter
func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.timeWindow)
	cleanedCount := 0

	for ip, timestamps := range rl.requestCounts {
		validRequests := make([]time.Time, 0, len(timestamps))
		for _, t := range timestamps {
			if t.After(windowStart) {
				validRequests = append(validRequests, t)
			}
		}

		if len(validRequests) == 0 {
			delete(rl.requestCounts, ip)
			cleanedCount++
		} else if len(validRequests) < len(timestamps) {
			rl.requestCounts[ip] = validRequests
		}
	}

	if cleanedCount > 0 {
		rl.logger.Debug("rate_limiter_cleanup", "removed_ips", cleanedCount, "remaining_ips", len(rl.requestCounts))
	}
}

// maskIP masks the last octet of IPv4 for logging privacy
func (rl *RateLimiter) maskIP(ip string) string {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return "unknown"
	}

	// IPv4
	if ipv4 := parsedIP.To4(); ipv4 != nil {
		return fmt.Sprintf("%d.%d.%d.0/24", ipv4[0], ipv4[1], ipv4[2])
	}

	// IPv6 - mask last 64 bits
	if ipv6 := parsedIP.To16(); ipv6 != nil {
		return fmt.Sprintf("%x:%x:%x:%x:0:0:0:0/64", ipv6[0:2], ipv6[2:4], ipv6[4:6], ipv6[6:8])
	}

	return "unknown"
}

// GetStats returns current rate limiter statistics (for monitoring)
func (rl *RateLimiter) GetStats() map[string]interface{} {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	return map[string]interface{}{
		"tracked_ips":   len(rl.requestCounts),
		"max_requests":  rl.maxRequests,
		"time_window":   rl.timeWindow.String(),
		"total_tracked": len(rl.requestCounts),
	}
}
