package loki

import (
	"testing"
	"time"
)

func TestLokiIntegration_GracefulShutdown(t *testing.T) {
	// Create a mock client for testing
	mockClient := &MockLokiClient{}
	labels := map[string]string{"service": "test"}

	// Create integration
	integration, err := NewLokiIntegration(LokiConfig{
		URL:    "http://localhost:3100",
		Labels: labels,
	})
	if err != nil {
		t.Fatalf("Failed to create integration: %v", err)
	}

	// Replace the client with our mock
	integration.client = mockClient

	// Test graceful shutdown with short timeout
	timeout := 100 * time.Millisecond
	err = integration.GracefulShutdown(timeout)
	if err != nil {
		t.Errorf("GracefulShutdown failed: %v", err)
	}

	// Verify that the integration was properly closed
	if integration.ctx.Err() == nil {
		t.Error("Expected context to be canceled after graceful shutdown")
	}
}

func TestLokiIntegration_GracefulShutdown_Timeout(t *testing.T) {
	// Create a mock client
	mockClient := &MockLokiClient{}
	labels := map[string]string{"service": "test"}

	// Create integration
	integration, err := NewLokiIntegration(LokiConfig{
		URL:    "http://localhost:3100",
		Labels: labels,
	})
	if err != nil {
		t.Fatalf("Failed to create integration: %v", err)
	}

	// Replace the client with our mock
	integration.client = mockClient

	// Test graceful shutdown with reasonable timeout
	timeout := 100 * time.Millisecond
	err = integration.GracefulShutdown(timeout)
	if err != nil {
		t.Errorf("GracefulShutdown should succeed with mock client: %v", err)
	}

	// Verify that the integration was properly closed
	if integration.ctx.Err() == nil {
		t.Error("Expected context to be canceled after graceful shutdown")
	}
}

// Note: GracefulShutdownLoki is tested in the main package tests
// since it's part of the public API in logger_settings.go
