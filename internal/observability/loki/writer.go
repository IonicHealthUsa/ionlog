package loki

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// LokiWriter implements io.Writer interface for Loki integration
type LokiWriter struct {
	client    ILokiClient
	labels    map[string]string
	buffer    []LogEntry
	bufferMux sync.Mutex
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewLokiWriter creates a new Loki writer
func NewLokiWriter(client ILokiClient, labels map[string]string) *LokiWriter {
	ctx, cancel := context.WithCancel(context.Background())

	writer := &LokiWriter{
		client: client,
		labels: labels,
		buffer: make([]LogEntry, 0),
		ctx:    ctx,
		cancel: cancel,
	}

	// Start background flusher
	go writer.backgroundFlusher()

	return writer
}

// Write implements io.Writer interface
func (lw *LokiWriter) Write(p []byte) (int, error) {
	// Parse the log entry from JSON
	var logData map[string]interface{}
	if err := json.Unmarshal(p, &logData); err != nil {
		return 0, fmt.Errorf("failed to parse log data: %w", err)
	}

	// Extract log information
	message, ok := logData["msg"].(string)
	if !ok {
		message = string(p) // Fallback to raw data
	}

	// Extract timestamp
	var timestamp time.Time
	if timeStr, ok := logData["time"].(string); ok {
		if parsedTime, err := time.Parse(time.RFC3339, timeStr); err == nil {
			timestamp = parsedTime
		} else {
			timestamp = time.Now()
		}
	} else {
		timestamp = time.Now()
	}

	// Create labels from log data
	labels := make(map[string]string)
	for k, v := range lw.labels {
		labels[k] = v
	}

	// Add log-specific labels
	if level, ok := logData["level"].(string); ok {
		labels["level"] = level
	}
	if file, ok := logData["file"].(string); ok {
		labels["file"] = file
	}
	if pkg, ok := logData["package"].(string); ok {
		labels["package"] = pkg
	}
	if function, ok := logData["function"].(string); ok {
		labels["function"] = function
	}
	if line, ok := logData["line"].(string); ok {
		labels["line"] = line
	}

	// Create log entry
	entry := LogEntry{
		Labels:  labels,
		Message: message,
		Time:    timestamp,
	}

	// Add to buffer
	lw.bufferMux.Lock()
	lw.buffer = append(lw.buffer, entry)
	bufferLen := len(lw.buffer)
	lw.bufferMux.Unlock()

	// Send immediately if buffer is getting large
	if bufferLen >= 10 {
		go lw.flushBuffer()
	}

	return len(p), nil
}

// Close closes the writer and flushes any remaining logs
func (lw *LokiWriter) Close() error {
	lw.cancel()
	return lw.flushBuffer()
}

// FlushBuffer sends buffered logs to Loki (public method)
func (lw *LokiWriter) FlushBuffer() error {
	return lw.flushBuffer()
}

// FlushBufferWithContext sends buffered logs to Loki with a custom context
func (lw *LokiWriter) FlushBufferWithContext(ctx context.Context) error {
	lw.bufferMux.Lock()
	if len(lw.buffer) == 0 {
		lw.bufferMux.Unlock()
		return nil
	}

	// Copy buffer and clear it
	logs := make([]LogEntry, len(lw.buffer))
	copy(logs, lw.buffer)
	lw.buffer = lw.buffer[:0]
	lw.bufferMux.Unlock()

	// Send to Loki with custom context
	return lw.client.PushLogs(ctx, logs)
}

// flushBuffer sends buffered logs to Loki
func (lw *LokiWriter) flushBuffer() error {
	lw.bufferMux.Lock()
	if len(lw.buffer) == 0 {
		lw.bufferMux.Unlock()
		return nil
	}

	// Copy buffer and clear it
	logs := make([]LogEntry, len(lw.buffer))
	copy(logs, lw.buffer)
	lw.buffer = lw.buffer[:0]
	lw.bufferMux.Unlock()

	// Send to Loki
	return lw.client.PushLogs(lw.ctx, logs)
}

// backgroundFlusher periodically flushes the buffer
func (lw *LokiWriter) backgroundFlusher() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-lw.ctx.Done():
			return
		case <-ticker.C:
			lw.flushBuffer()
		}
	}
}
