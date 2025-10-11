package loki

import (
	"os"
	"testing"
	"time"
)

func TestDefaultLokiConfig(t *testing.T) {
	config := DefaultLokiConfig()

	if config.URL == "" {
		t.Error("Default URL should not be empty")
	}

	if config.BatchSize <= 0 {
		t.Error("Default batch size should be positive")
	}

	if config.Timeout <= 0 {
		t.Error("Default timeout should be positive")
	}

	if config.Labels == nil {
		t.Error("Default labels should not be nil")
	}
}

func TestWithLoki(t *testing.T) {
	url := "http://test:3100"
	labels := map[string]string{"service": "test", "env": "dev"}

	config := WithLoki(url, labels)

	if config.URL != url {
		t.Errorf("Expected URL %s, got %s", url, config.URL)
	}

	if config.Labels["service"] != "test" {
		t.Errorf("Expected service label 'test', got %s", config.Labels["service"])
	}

	if config.Labels["env"] != "dev" {
		t.Errorf("Expected env label 'dev', got %s", config.Labels["env"])
	}
}

func TestWithLokiAuth(t *testing.T) {
	config := DefaultLokiConfig()
	username := "testuser"
	password := "testpass"

	config = WithLokiAuth(config, username, password)

	if config.Username != username {
		t.Errorf("Expected username %s, got %s", username, config.Username)
	}

	if config.Password != password {
		t.Errorf("Expected password %s, got %s", password, config.Password)
	}
}

func TestWithLokiTenant(t *testing.T) {
	config := DefaultLokiConfig()
	tenantID := "tenant123"

	config = WithLokiTenant(config, tenantID)

	if config.TenantID != tenantID {
		t.Errorf("Expected tenant ID %s, got %s", tenantID, config.TenantID)
	}
}

func TestWithLokiBatchSize(t *testing.T) {
	config := DefaultLokiConfig()
	batchSize := 200

	config = WithLokiBatchSize(config, batchSize)

	if config.BatchSize != batchSize {
		t.Errorf("Expected batch size %d, got %d", batchSize, config.BatchSize)
	}
}

func TestWithLokiTimeout(t *testing.T) {
	config := DefaultLokiConfig()
	timeout := 60 * time.Second

	config = WithLokiTimeout(config, timeout)

	if config.Timeout != timeout {
		t.Errorf("Expected timeout %v, got %v", timeout, config.Timeout)
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	// Test with existing environment variable
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	value := getEnvOrDefault("TEST_VAR", "default")
	if value != "test_value" {
		t.Errorf("Expected 'test_value', got %s", value)
	}

	// Test with non-existing environment variable
	value = getEnvOrDefault("NON_EXISTING_VAR", "default")
	if value != "default" {
		t.Errorf("Expected 'default', got %s", value)
	}
}

func TestGetEnvIntOrDefault(t *testing.T) {
	// Test with valid integer environment variable
	os.Setenv("TEST_INT", "42")
	defer os.Unsetenv("TEST_INT")

	value := getEnvIntOrDefault("TEST_INT", 0)
	if value != 42 {
		t.Errorf("Expected 42, got %d", value)
	}

	// Test with invalid integer environment variable
	os.Setenv("TEST_INVALID_INT", "not_a_number")
	defer os.Unsetenv("TEST_INVALID_INT")

	value = getEnvIntOrDefault("TEST_INVALID_INT", 10)
	if value != 10 {
		t.Errorf("Expected 10, got %d", value)
	}

	// Test with non-existing environment variable
	value = getEnvIntOrDefault("NON_EXISTING_INT", 5)
	if value != 5 {
		t.Errorf("Expected 5, got %d", value)
	}
}

func TestGetEnvDurationOrDefault(t *testing.T) {
	// Test with valid duration environment variable
	os.Setenv("TEST_DURATION", "30s")
	defer os.Unsetenv("TEST_DURATION")

	value := getEnvDurationOrDefault("TEST_DURATION", time.Second)
	if value != 30*time.Second {
		t.Errorf("Expected 30s, got %v", value)
	}

	// Test with invalid duration environment variable
	os.Setenv("TEST_INVALID_DURATION", "not_a_duration")
	defer os.Unsetenv("TEST_INVALID_DURATION")

	value = getEnvDurationOrDefault("TEST_INVALID_DURATION", 5*time.Second)
	if value != 5*time.Second {
		t.Errorf("Expected 5s, got %v", value)
	}

	// Test with non-existing environment variable
	value = getEnvDurationOrDefault("NON_EXISTING_DURATION", 10*time.Second)
	if value != 10*time.Second {
		t.Errorf("Expected 10s, got %v", value)
	}
}
