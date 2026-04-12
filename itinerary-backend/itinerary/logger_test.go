package itinerary

import (
	"testing"
)

// TestLoggerDebug verifies debug logging
func TestLoggerDebug(t *testing.T) {
	logger := &Logger{}

	// Should not panic
	logger.Debug("test message", "key", "value")
}

// TestLoggerInfo verifies info logging
func TestLoggerInfo(t *testing.T) {
	logger := &Logger{}

	// Should not panic
	logger.Info("info message", "level", "info")
}

// TestLoggerError verifies error logging
func TestLoggerError(t *testing.T) {
	logger := &Logger{}

	// Should not panic
	logger.Error("error message", "error", "test error")
}

// TestLoggerWarn verifies warning logging
func TestLoggerWarn(t *testing.T) {
	logger := &Logger{}

	// Should not panic
	logger.Warn("warning message", "warning", "test warning")
}

// TestLoggerMultipleFields verifies logging with multiple fields
func TestLoggerMultipleFields(t *testing.T) {
	logger := &Logger{}

	// Should handle multiple key-value pairs
	logger.Info("complex message",
		"user_id", "user-001",
		"trip_id", "trip-001",
		"action", "created",
		"timestamp", "2026-03-23T10:00:00Z",
	)
}

// TestLoggerEmptyMessage verifies logging empty messages
func TestLoggerEmptyMessage(t *testing.T) {
	logger := &Logger{}

	// Should handle empty message
	logger.Info("")
	logger.Debug("")
	logger.Error("")
}

// TestLoggerSpecialCharacters verifies logging special characters
func TestLoggerSpecialCharacters(t *testing.T) {
	logger := &Logger{}

	logger.Debug("message with special chars", "value", "test@#$%^&*()")
	logger.Error("unicode message", "emoji", "🎉🌍✈️")
}
