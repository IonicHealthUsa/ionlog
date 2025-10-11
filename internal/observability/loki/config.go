package loki

import (
	"os"
	"strconv"
	"time"
)

// DefaultLokiConfig returns a default Loki configuration
func DefaultLokiConfig() LokiConfig {
	return LokiConfig{
		URL:       getEnvOrDefault("LOKI_URL", "http://localhost:3100"),
		Username:  getEnvOrDefault("LOKI_USERNAME", ""),
		Password:  getEnvOrDefault("LOKI_PASSWORD", ""),
		TenantID:  getEnvOrDefault("LOKI_TENANT_ID", ""),
		Labels:    getDefaultLabels(),
		BatchSize: getEnvIntOrDefault("LOKI_BATCH_SIZE", 100),
		Timeout:   getEnvDurationOrDefault("LOKI_TIMEOUT", 30*time.Second),
	}
}

// getDefaultLabels returns default labels for Loki
func getDefaultLabels() map[string]string {
	return map[string]string{
		"service":     getEnvOrDefault("SERVICE_NAME", "ionlog"),
		"environment": getEnvOrDefault("ENVIRONMENT", "development"),
		"version":     getEnvOrDefault("VERSION", "1.0.0"),
		"component":   getEnvOrDefault("COMPONENT", "logger"),
	}
}

// getEnvOrDefault gets an environment variable or returns a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvIntOrDefault gets an environment variable as int or returns a default value
func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvDurationOrDefault gets an environment variable as duration or returns a default value
func getEnvDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// WithLoki creates a Loki configuration builder
func WithLoki(url string, labels map[string]string) LokiConfig {
	config := DefaultLokiConfig()
	config.URL = url

	// Merge with default labels
	for k, v := range labels {
		config.Labels[k] = v
	}

	return config
}

// WithLokiAuth adds authentication to Loki configuration
func WithLokiAuth(config LokiConfig, username, password string) LokiConfig {
	config.Username = username
	config.Password = password
	return config
}

// WithLokiTenant adds tenant ID to Loki configuration
func WithLokiTenant(config LokiConfig, tenantID string) LokiConfig {
	config.TenantID = tenantID
	return config
}

// WithLokiBatchSize sets the batch size for Loki configuration
func WithLokiBatchSize(config LokiConfig, batchSize int) LokiConfig {
	config.BatchSize = batchSize
	return config
}

// WithLokiTimeout sets the timeout for Loki configuration
func WithLokiTimeout(config LokiConfig, timeout time.Duration) LokiConfig {
	config.Timeout = timeout
	return config
}
