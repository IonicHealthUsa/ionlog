package loki

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// LokiConfig holds configuration for Loki connection
type LokiConfig struct {
	URL       string            `json:"url"`
	Username  string            `json:"username,omitempty"`
	Password  string            `json:"password,omitempty"`
	TenantID  string            `json:"tenant_id,omitempty"`
	Labels    map[string]string `json:"labels"`
	BatchSize int               `json:"batch_size"`
	Timeout   time.Duration     `json:"timeout"`
}

// Stream represents a Loki log stream
type Stream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

// PushRequest represents a Loki push request
type PushRequest struct {
	Streams []Stream `json:"streams"`
}

// LokiClient manages the connection and writing to Loki
type LokiClient struct {
	config     LokiConfig
	httpClient *http.Client
	buffer     []Stream
	bufferSize int
}

// ILokiClient defines the interface for Loki operations
type ILokiClient interface {
	PushLog(ctx context.Context, labels map[string]string, message string, timestamp time.Time) error
	PushLogs(ctx context.Context, logs []LogEntry) error
	Flush(ctx context.Context) error
	Close() error
}

// LogEntry represents a single log entry for Loki
type LogEntry struct {
	Labels  map[string]string
	Message string
	Time    time.Time
}

// NewLokiClient creates a new Loki client
func NewLokiClient(config LokiConfig) (*LokiClient, error) {
	if config.URL == "" {
		return nil, fmt.Errorf("loki URL must be configured")
	}

	if config.BatchSize <= 0 {
		config.BatchSize = 100
	}

	if config.Timeout <= 0 {
		config.Timeout = 30 * time.Second
	}

	if config.Labels == nil {
		config.Labels = make(map[string]string)
	}

	client := &LokiClient{
		config:     config,
		httpClient: &http.Client{Timeout: config.Timeout},
		buffer:     make([]Stream, 0, config.BatchSize),
		bufferSize: config.BatchSize,
	}

	return client, nil
}

// PushLog adds a log entry to the buffer and sends if batch size is reached
func (lc *LokiClient) PushLog(ctx context.Context, labels map[string]string, message string, timestamp time.Time) error {
	// Merge with default labels
	mergedLabels := make(map[string]string)
	for k, v := range lc.config.Labels {
		mergedLabels[k] = v
	}
	for k, v := range labels {
		mergedLabels[k] = v
	}

	// Create log entry
	entry := LogEntry{
		Labels:  mergedLabels,
		Message: message,
		Time:    timestamp,
	}

	// Add to buffer
	lc.addToBuffer(entry)

	// Flush if buffer is full
	if len(lc.buffer) >= lc.bufferSize {
		return lc.Flush(ctx)
	}

	return nil
}

// PushLogs sends multiple log entries to Loki
func (lc *LokiClient) PushLogs(ctx context.Context, logs []LogEntry) error {
	if len(logs) == 0 {
		return nil
	}

	// Group logs by labels
	streams := make(map[string]Stream)

	for _, log := range logs {
		// Merge with default labels
		mergedLabels := make(map[string]string)
		for k, v := range lc.config.Labels {
			mergedLabels[k] = v
		}
		for k, v := range log.Labels {
			mergedLabels[k] = v
		}

		// Create stream key from labels
		streamKey := lc.createStreamKey(mergedLabels)

		stream, exists := streams[streamKey]
		if !exists {
			stream = Stream{
				Stream: mergedLabels,
				Values: make([][]string, 0),
			}
		}

		// Add log entry to stream
		timestamp := fmt.Sprintf("%d", log.Time.UnixNano())
		stream.Values = append(stream.Values, []string{timestamp, log.Message})
		streams[streamKey] = stream
	}

	// Convert to slice
	streamSlice := make([]Stream, 0, len(streams))
	for _, stream := range streams {
		streamSlice = append(streamSlice, stream)
	}

	return lc.sendToLoki(ctx, streamSlice)
}

// Flush sends buffered logs to Loki
func (lc *LokiClient) Flush(ctx context.Context) error {
	if len(lc.buffer) == 0 {
		return nil
	}

	streams := make([]Stream, len(lc.buffer))
	copy(streams, lc.buffer)
	lc.buffer = lc.buffer[:0] // Clear buffer

	return lc.sendToLoki(ctx, streams)
}

// Close flushes any remaining logs and closes the client
func (lc *LokiClient) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), lc.config.Timeout)
	defer cancel()

	return lc.Flush(ctx)
}

// addToBuffer adds a log entry to the internal buffer
func (lc *LokiClient) addToBuffer(entry LogEntry) {
	// Create stream key from labels
	streamKey := lc.createStreamKey(entry.Labels)

	// Find existing stream or create new one
	var stream *Stream
	for i := range lc.buffer {
		if lc.createStreamKey(lc.buffer[i].Stream) == streamKey {
			stream = &lc.buffer[i]
			break
		}
	}

	if stream == nil {
		// Create new stream
		lc.buffer = append(lc.buffer, Stream{
			Stream: entry.Labels,
			Values: make([][]string, 0),
		})
		stream = &lc.buffer[len(lc.buffer)-1]
	}

	// Add log entry to stream
	timestamp := fmt.Sprintf("%d", entry.Time.UnixNano())
	stream.Values = append(stream.Values, []string{timestamp, entry.Message})
}

// createStreamKey creates a unique key for a set of labels
func (lc *LokiClient) createStreamKey(labels map[string]string) string {
	// Simple key generation - in production, you might want a more sophisticated approach
	key := ""
	for k, v := range labels {
		key += fmt.Sprintf("%s=%s,", k, v)
	}
	return key
}

// sendToLoki sends streams to Loki
func (lc *LokiClient) sendToLoki(ctx context.Context, streams []Stream) error {
	if len(streams) == 0 {
		return nil
	}

	// Create push request
	request := PushRequest{
		Streams: streams,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal Loki request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", lc.config.URL+"/loki/api/v1/push", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	if lc.config.Username != "" && lc.config.Password != "" {
		req.SetBasicAuth(lc.config.Username, lc.config.Password)
	}
	if lc.config.TenantID != "" {
		req.Header.Set("X-Scope-OrgID", lc.config.TenantID)
	}

	// Send request
	resp, err := lc.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to Loki: %w", err)
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("loki returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
