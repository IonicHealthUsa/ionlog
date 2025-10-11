package loki

import (
	"testing"
	"time"
)

func TestNewLokiIntegration(t *testing.T) {
	config := LokiConfig{
		URL:       "http://localhost:3100",
		BatchSize: 100,
		Timeout:   30 * time.Second,
		Labels:    map[string]string{"service": "test"},
	}

	integration, err := NewLokiIntegration(config)
	if err != nil {
		t.Fatalf("NewLokiIntegration() error = %v", err)
	}

	if integration == nil {
		t.Fatal("NewLokiIntegration() returned nil")
	}

	if integration.client == nil {
		t.Error("Integration client is nil")
	}

	if integration.writer == nil {
		t.Error("Integration writer is nil")
	}

	// Clean up
	integration.Close()
}

func TestLokiIntegration_Writer(t *testing.T) {
	config := LokiConfig{
		URL:       "http://localhost:3100",
		BatchSize: 100,
		Timeout:   30 * time.Second,
		Labels:    map[string]string{"service": "test"},
	}

	integration, err := NewLokiIntegration(config)
	if err != nil {
		t.Fatalf("NewLokiIntegration() error = %v", err)
	}
	defer integration.Close()

	writer := integration.Writer()
	if writer == nil {
		t.Error("Writer() returned nil")
	}
}

func TestLokiIntegration_Client(t *testing.T) {
	config := LokiConfig{
		URL:       "http://localhost:3100",
		BatchSize: 100,
		Timeout:   30 * time.Second,
		Labels:    map[string]string{"service": "test"},
	}

	integration, err := NewLokiIntegration(config)
	if err != nil {
		t.Fatalf("NewLokiIntegration() error = %v", err)
	}
	defer integration.Close()

	client := integration.Client()
	if client == nil {
		t.Error("Client() returned nil")
	}
}

func TestLokiIntegration_PushLog(t *testing.T) {
	config := LokiConfig{
		URL:       "http://localhost:3100",
		BatchSize: 100,
		Timeout:   30 * time.Second,
		Labels:    map[string]string{"service": "test"},
	}

	integration, err := NewLokiIntegration(config)
	if err != nil {
		t.Fatalf("NewLokiIntegration() error = %v", err)
	}
	defer integration.Close()

	labels := map[string]string{"level": "info"}
	message := "test message"
	timestamp := time.Now()

	// This will fail in tests since we don't have a real Loki instance
	err = integration.PushLog(labels, message, timestamp)
	// Note: PushLog adds to buffer, so it won't error immediately
	if err != nil {
		t.Logf("PushLog error (expected): %v", err)
	}
}

func TestLokiIntegration_PushLogs(t *testing.T) {
	config := LokiConfig{
		URL:       "http://localhost:3100",
		BatchSize: 100,
		Timeout:   30 * time.Second,
		Labels:    map[string]string{"service": "test"},
	}

	integration, err := NewLokiIntegration(config)
	if err != nil {
		t.Fatalf("NewLokiIntegration() error = %v", err)
	}
	defer integration.Close()

	logs := []LogEntry{
		{
			Labels:  map[string]string{"level": "info"},
			Message: "test message 1",
			Time:    time.Now(),
		},
		{
			Labels:  map[string]string{"level": "error"},
			Message: "test message 2",
			Time:    time.Now(),
		},
	}

	// PushLogs buffers logs and doesn't immediately fail
	// The error will occur during the actual HTTP request
	err = integration.PushLogs(logs)
	if err != nil {
		// Log the expected error but don't fail the test
		t.Logf("Expected error when pushing to non-existent Loki instance: %v", err)
	}
}

func TestLokiIntegration_Flush(t *testing.T) {
	config := LokiConfig{
		URL:       "http://localhost:3100",
		BatchSize: 100,
		Timeout:   30 * time.Second,
		Labels:    map[string]string{"service": "test"},
	}

	integration, err := NewLokiIntegration(config)
	if err != nil {
		t.Fatalf("NewLokiIntegration() error = %v", err)
	}
	defer integration.Close()

	// Flush should not error even with empty buffer
	err = integration.Flush()
	if err != nil {
		t.Errorf("Flush() error = %v", err)
	}
}

func TestLokiIntegration_HealthCheck(t *testing.T) {
	config := LokiConfig{
		URL:       "http://localhost:3100",
		BatchSize: 100,
		Timeout:   30 * time.Second,
		Labels:    map[string]string{"service": "test"},
	}

	integration, err := NewLokiIntegration(config)
	if err != nil {
		t.Fatalf("NewLokiIntegration() error = %v", err)
	}
	defer integration.Close()

	// Health check will fail in tests since we don't have a real Loki instance
	err = integration.HealthCheck()
	// Note: HealthCheck adds to buffer, so it won't error immediately
	if err != nil {
		t.Logf("HealthCheck error (expected): %v", err)
	}
}

func TestLokiIntegration_GetStats(t *testing.T) {
	config := LokiConfig{
		URL:       "http://localhost:3100",
		BatchSize: 100,
		Timeout:   30 * time.Second,
		Labels:    map[string]string{"service": "test"},
	}

	integration, err := NewLokiIntegration(config)
	if err != nil {
		t.Fatalf("NewLokiIntegration() error = %v", err)
	}
	defer integration.Close()

	stats := integration.GetStats()
	if stats == nil {
		t.Error("GetStats() returned nil")
	}

	// Check expected stats
	if stats["client_type"] != "loki" {
		t.Errorf("Expected client_type to be 'loki', got %v", stats["client_type"])
	}

	if stats["active"] != true {
		t.Errorf("Expected active to be true, got %v", stats["active"])
	}
}

func TestLokiIntegration_Close(t *testing.T) {
	config := LokiConfig{
		URL:       "http://localhost:3100",
		BatchSize: 100,
		Timeout:   30 * time.Second,
		Labels:    map[string]string{"service": "test"},
	}

	integration, err := NewLokiIntegration(config)
	if err != nil {
		t.Fatalf("NewLokiIntegration() error = %v", err)
	}

	err = integration.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}

	// Check that context is canceled
	select {
	case <-integration.ctx.Done():
		// Expected
	default:
		t.Error("Expected context to be canceled after Close()")
	}
}
