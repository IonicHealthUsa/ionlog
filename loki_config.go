package ionlog

import (
	"os"
	"sync"
	"time"

	"github.com/IonicHealthUsa/ionlog/internal/observability/loki"
)

var (
	lokiShutdownTimeout = getDefaultLokiShutdownTimeout() // Configurable timeout
	lokiShutdownMutex   sync.RWMutex

	lokiIntegration *loki.LokiIntegration
	lokiMutex       sync.RWMutex
)

// LokiConfig represents the configuration for Loki integration
// This is a public struct that users can use without importing internal packages
type LokiConfig struct {
	URL       string
	Labels    map[string]string
	Username  string
	Password  string
	TenantID  string
	BatchSize int
	Timeout   time.Duration
}

// WithLoki creates a basic Loki configuration with URL and labels
func WithLoki(url string, labels map[string]string) LokiConfig {
	return LokiConfig{
		URL:       url,
		Labels:    labels,
		BatchSize: 1000,
		Timeout:   30 * time.Second,
	}
}

// WithLokiAuth adds authentication to a Loki configuration
func WithLokiAuth(config LokiConfig, username, password string) LokiConfig {
	config.Username = username
	config.Password = password
	return config
}

// WithLokiTenant adds tenant ID to a Loki configuration
func WithLokiTenant(config LokiConfig, tenantID string) LokiConfig {
	config.TenantID = tenantID
	return config
}

// WithLokiBatchSize sets the batch size for a Loki configuration
func WithLokiBatchSize(config LokiConfig, batchSize int) LokiConfig {
	config.BatchSize = batchSize
	return config
}

// WithLokiTimeout sets the timeout for a Loki configuration
func WithLokiTimeout(config LokiConfig, timeout time.Duration) LokiConfig {
	config.Timeout = timeout
	return config
}

// toInternalConfig converts the public LokiConfig to the internal loki.LokiConfig
func (lc LokiConfig) toInternalConfig() loki.LokiConfig {
	internalConfig := loki.LokiConfig{
		URL:       lc.URL,
		Labels:    lc.Labels,
		Username:  lc.Username,
		Password:  lc.Password,
		TenantID:  lc.TenantID,
		BatchSize: lc.BatchSize,
		Timeout:   lc.Timeout,
	}
	return internalConfig
}

// getDefaultLokiShutdownTimeout returns the default timeout from environment or fallback
func getDefaultLokiShutdownTimeout() time.Duration {
	// Check environment variable first
	if envTimeout := os.Getenv("IONLOG_LOKI_SHUTDOWN_TIMEOUT"); envTimeout != "" {
		if timeout, err := time.ParseDuration(envTimeout); err == nil {
			return timeout
		}
	}

	// Fallback to default
	return 10 * time.Second
}

func getLokiIntegration() *loki.LokiIntegration {
	lokiMutex.RLock()
	defer lokiMutex.RUnlock()
	return lokiIntegration
}

// GetLokiShutdownTimeout returns the current Loki shutdown timeout
// This respects the following priority:
// 1. Programmatically set timeout (via WithLokiShutdownTimeout)
// 2. Environment variable (IONLOG_LOKI_SHUTDOWN_TIMEOUT)
// 3. Default timeout (10 seconds)
func GetLokiShutdownTimeout() time.Duration {
	lokiShutdownMutex.RLock()
	defer lokiShutdownMutex.RUnlock()
	return lokiShutdownTimeout
}

// ResetLokiShutdownTimeout resets the timeout to use environment variable or default
// This is useful for testing or when you want to revert to environment-based configuration
func ResetLokiShutdownTimeout() {
	lokiShutdownMutex.Lock()
	lokiShutdownTimeout = getDefaultLokiShutdownTimeout()
	lokiShutdownMutex.Unlock()
}

// setLokiShutdownTimeout sets the timeout for graceful Loki shutdown
// This is an internal function used by WithLokiShutdownTimeout
func setLokiShutdownTimeout(timeout time.Duration) {
	lokiShutdownMutex.Lock()
	// Ensure minimum timeout
	if timeout < 1*time.Second {
		timeout = 1 * time.Second
	}
	lokiShutdownTimeout = timeout
	lokiShutdownMutex.Unlock()
}

// gracefulShutdownLokiInternal gracefully shuts down the Loki integration with a timeout
// This is an internal function used by ionlog.Stop() for transparent shutdown
func gracefulShutdownLokiInternal(timeout time.Duration) error {
	lokiMutex.Lock()
	defer lokiMutex.Unlock()

	if lokiIntegration == nil {
		return nil // No Loki integration to shutdown
	}

	// Ensure minimum timeout to prevent immediate failures
	if timeout < 1*time.Second {
		timeout = 1 * time.Second
	}

	_ = lokiIntegration.GracefulShutdown(timeout)
	lokiIntegration = nil // Clear the reference

	// Return nil to make shutdown completely transparent
	// Any errors are handled internally and don't affect the application
	return nil
}
