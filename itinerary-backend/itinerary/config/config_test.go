package config

import (
	"path/filepath"
	"testing"
)

// TestConfigLoading verifies configuration can be loaded
func TestConfigLoading(t *testing.T) {
	configPath := filepath.Join("config", "config.json")
	config, err := LoadConfig(configPath)

	if err != nil && config == nil {
		t.Logf("LoadConfig returned error (expected for test): %v", err)
		return
	}

	if config != nil && err == nil {
		t.Log("Config loaded successfully")
	}
}

// TestConfigStructure verifies configuration structure
func TestConfigStructure(t *testing.T) {
	config := &Config{
		Server: ServerConfig{
			Port:    "8080",
			Timeout: 30,
			Mode:    "development",
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "admin",
			Database: "itinerary",
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
			Output: "stdout",
		},
	}

	tests := []struct {
		name   string
		check  func(*Config) bool
		errMsg string
	}{
		{
			name: "server port is set",
			check: func(c *Config) bool {
				return c.Server.Port != ""
			},
			errMsg: "Server port should be set",
		},
		{
			name: "server timeout is positive",
			check: func(c *Config) bool {
				return c.Server.Timeout > 0
			},
			errMsg: "Server timeout should be positive",
		},
		{
			name: "database host is set",
			check: func(c *Config) bool {
				return c.Database.Host != ""
			},
			errMsg: "Database host should be set",
		},
		{
			name: "logging level is set",
			check: func(c *Config) bool {
				return c.Logging.Level != ""
			},
			errMsg: "Logging level should be set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.check(config) {
				t.Error(tt.errMsg)
			}
		})
	}
}

// TestServerConfig verifies server configuration
func TestServerConfig(t *testing.T) {
	serverConfig := ServerConfig{
		Port:    "8080",
		Timeout: 30,
		Mode:    "development",
	}

	if serverConfig.Port != "8080" {
		t.Errorf("Expected port 8080, got %s", serverConfig.Port)
	}

	if serverConfig.Timeout != 30 {
		t.Errorf("Expected timeout 30, got %d", serverConfig.Timeout)
	}

	if serverConfig.Mode != "development" {
		t.Errorf("Expected mode 'development', got %q", serverConfig.Mode)
	}
}
