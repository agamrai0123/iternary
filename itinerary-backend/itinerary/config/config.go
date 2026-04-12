package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	Logging  LoggingConfig  `json:"logging"`
	API      APIConfig      `json:"api"`
}

// ServerConfig represents server settings
type ServerConfig struct {
	Port    string `json:"port"`
	Timeout int    `json:"timeout"`
	Mode    string `json:"mode"`
}

// DatabaseConfig represents database settings
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Database string `json:"database"`
	Password string `json:"password"` // Load from env in production
}

// LoggingConfig represents logging settings
type LoggingConfig struct {
	Level  string `json:"level"`
	Format string `json:"format"`
	Output string `json:"output"`
}

// APIConfig represents API settings
type APIConfig struct {
	PageSize             int        `json:"page_size"`
	MaxItemsPerItinerary int        `json:"max_items_per_itinerary"`
	CORS                 CORSConfig `json:"cors"`
}

// CORSConfig represents CORS settings
type CORSConfig struct {
	AllowedOrigins []string `json:"allowed_origins"`
	AllowedMethods []string `json:"allowed_methods"`
	AllowedHeaders []string `json:"allowed_headers"`
}

// LoadConfig loads configuration from JSON file and environment variables
func LoadConfig(configPath string) (*Config, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Override with environment variables (production best practice)
	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
		config.Database.Password = dbPassword
	}
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		config.Database.Host = dbHost
	}

	return &config, nil
}
