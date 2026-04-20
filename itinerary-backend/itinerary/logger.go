package itinerary

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/yourusername/itinerary-backend/itinerary/common"
)

// (Moved Logger struct and methods to common package)

// NewLogger creates a new structured logger using zerolog
func NewLogger(config *Config) *common.Logger {
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

	return &common.Logger{Log: z}
}
