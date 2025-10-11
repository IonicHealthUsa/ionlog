package loki

import (
	"context"
	"fmt"
	"io"
	"time"
)

// LokiIntegration provides integration with the ionlog logger
type LokiIntegration struct {
	client ILokiClient
	writer *LokiWriter
	ctx    context.Context
	cancel context.CancelFunc
}

// NewLokiIntegration creates a new Loki integration
func NewLokiIntegration(config LokiConfig) (*LokiIntegration, error) {
	client, err := NewLokiClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Loki client: %w", err)
	}

	writer := NewLokiWriter(client, config.Labels)

	ctx, cancel := context.WithCancel(context.Background())

	integration := &LokiIntegration{
		client: client,
		writer: writer,
		ctx:    ctx,
		cancel: cancel,
	}

	return integration, nil
}

// Writer returns the Loki writer that implements io.Writer
func (li *LokiIntegration) Writer() io.Writer {
	return li.writer
}

// Client returns the Loki client
func (li *LokiIntegration) Client() ILokiClient {
	return li.client
}

// Close closes the integration and flushes any remaining logs
func (li *LokiIntegration) Close() error {
	li.cancel()

	if err := li.writer.Close(); err != nil {
		return fmt.Errorf("failed to close Loki writer: %w", err)
	}

	if err := li.client.Close(); err != nil {
		return fmt.Errorf("failed to close Loki client: %w", err)
	}

	return nil
}

// GracefulShutdown gracefully shuts down the integration with a timeout
// It ensures all buffered logs are sent to Loki before closing
func (li *LokiIntegration) GracefulShutdown(timeout time.Duration) error {
	// Create a context with timeout for graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Flush any remaining logs with timeout (using fresh context)
	done := make(chan error, 1)
	go func() {
		// Use the shutdown context for the flush operation
		done <- li.writer.FlushBufferWithContext(shutdownCtx)
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("failed to flush logs during shutdown: %w", err)
		}
	case <-shutdownCtx.Done():
		return fmt.Errorf("shutdown timeout after %v - some logs may not have been sent", timeout)
	}

	// Signal shutdown to background processes after flush is complete
	li.cancel()

	// Close the writer
	if err := li.writer.Close(); err != nil {
		return fmt.Errorf("failed to close Loki writer: %w", err)
	}

	// Close the client
	if err := li.client.Close(); err != nil {
		return fmt.Errorf("failed to close Loki client: %w", err)
	}

	return nil
}

// PushLog directly pushes a log entry to Loki
func (li *LokiIntegration) PushLog(labels map[string]string, message string, timestamp time.Time) error {
	return li.client.PushLog(li.ctx, labels, message, timestamp)
}

// PushLogs directly pushes multiple log entries to Loki
func (li *LokiIntegration) PushLogs(logs []LogEntry) error {
	return li.client.PushLogs(li.ctx, logs)
}

// Flush flushes any buffered logs
func (li *LokiIntegration) Flush() error {
	return li.client.Flush(li.ctx)
}

// HealthCheck checks if Loki is accessible
func (li *LokiIntegration) HealthCheck() error {
	// Try to push a test log entry
	testLabels := map[string]string{
		"health_check": "true",
		"timestamp":    time.Now().Format(time.RFC3339),
	}

	return li.client.PushLog(li.ctx, testLabels, "health check", time.Now())
}

// GetStats returns statistics about the Loki integration
func (li *LokiIntegration) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"client_type":      "loki",
		"active":           li.ctx.Err() == nil,
		"context_canceled": li.ctx.Err() != nil,
	}
}
