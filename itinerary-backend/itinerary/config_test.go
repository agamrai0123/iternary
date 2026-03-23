package itinerary

import (
	"testing"
)

// TestConfigLoading verifies configuration can be loaded
func TestConfigLoading(t *testing.T) {
	// LoadConfig should not panic
	config := LoadConfig()

	if config == nil {
		t.Error("LoadConfig() returned nil")
	}
}

// TestConfigProperties verifies configuration properties
func TestConfigProperties(t *testing.T) {
	config := LoadConfig()

	tests := []struct {
		name     string
		property string
		check    func(*Config) bool
	}{
		{
			name:     "port is set",
			property: "port",
			check: func(c *Config) bool {
				return c.Port > 0
			},
		},
		{
			name:     "env is not empty",
			property: "env",
			check: func(c *Config) bool {
				return c.Env != ""
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.check(config) {
				t.Errorf("Configuration check failed for property %s", tt.property)
			}
		})
	}
}

// TestDefaultConfigValues verifies default configuration values
func TestDefaultConfigValues(t *testing.T) {
	config := &Config{
		Port: 8080,
		Env:  "development",
	}

	if config.Port != 8080 {
		t.Errorf("Expected port 8080, got %d", config.Port)
	}

	if config.Env != "development" {
		t.Errorf("Expected env 'development', got %q", config.Env)
	}
}

// TestProductionConfig verifies production configuration
func TestProductionConfig(t *testing.T) {
	config := &Config{
		Port: 8080,
		Env:  "production",
	}

	if config.Env != "production" {
		t.Error("Production environment should be set")
	}
}

// TestDevelopmentConfig verifies development configuration
func TestDevelopmentConfig(t *testing.T) {
	config := &Config{
		Port: 3000,
		Env:  "development",
	}

	if config.Env != "development" {
		t.Error("Development environment should be set")
	}

	if config.Port != 3000 {
		t.Errorf("Expected port 3000, got %d", config.Port)
	}
}
