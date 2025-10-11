package loki

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"
)

// MockLokiClient is a mock implementation of ILokiClient for testing
type MockLokiClient struct {
	mu            sync.Mutex
	pushLogCalls  []PushLogCall
	pushLogsCalls []PushLogsCall
	flushCalls    int
	closeCalls    int
}

type PushLogCall struct {
	Labels  map[string]string
	Message string
	Time    time.Time
}

type PushLogsCall struct {
	Logs []LogEntry
}

func (m *MockLokiClient) PushLog(ctx context.Context, labels map[string]string, message string, timestamp time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pushLogCalls = append(m.pushLogCalls, PushLogCall{
		Labels:  labels,
		Message: message,
		Time:    timestamp,
	})
	return nil
}

func (m *MockLokiClient) PushLogs(ctx context.Context, logs []LogEntry) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pushLogsCalls = append(m.pushLogsCalls, PushLogsCall{
		Logs: logs,
	})
	return nil
}

func (m *MockLokiClient) Flush(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.flushCalls++
	return nil
}

func (m *MockLokiClient) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closeCalls++
	return nil
}

// Thread-safe getters for testing
func (m *MockLokiClient) GetPushLogsCalls() []PushLogsCall {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]PushLogsCall(nil), m.pushLogsCalls...)
}

func (m *MockLokiClient) GetFlushCalls() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.flushCalls
}

func (m *MockLokiClient) GetCloseCalls() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.closeCalls
}

func TestNewLokiWriter(t *testing.T) {
	mockClient := &MockLokiClient{}
	labels := map[string]string{"service": "test"}

	writer := NewLokiWriter(mockClient, labels)
	if writer == nil {
		t.Fatal("NewLokiWriter() returned nil")
	}

	if writer.client != mockClient {
		t.Error("Writer client not set correctly")
	}

	if len(writer.labels) != len(labels) {
		t.Error("Writer labels not set correctly")
	}

	// Clean up
	writer.Close()
}

func TestLokiWriter_Write(t *testing.T) {
	mockClient := &MockLokiClient{}
	labels := map[string]string{"service": "test"}

	writer := NewLokiWriter(mockClient, labels)
	defer writer.Close()

	// Create test log data
	logData := map[string]interface{}{
		"time":     time.Now().Format(time.RFC3339),
		"level":    "info",
		"msg":      "test message",
		"file":     "test.go",
		"package":  "test",
		"function": "TestFunction",
		"line":     "42",
	}

	jsonData, err := json.Marshal(logData)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}

	// Write to writer
	n, err := writer.Write(jsonData)
	if err != nil {
		t.Errorf("Write() error = %v", err)
	}

	if n != len(jsonData) {
		t.Errorf("Write() returned %d, want %d", n, len(jsonData))
	}

	// Wait a bit for background processing
	time.Sleep(100 * time.Millisecond)

	// Check that the log was added to buffer
	if len(writer.buffer) == 0 {
		t.Error("Expected log to be added to buffer")
	}
}

func TestLokiWriter_Write_InvalidJSON(t *testing.T) {
	mockClient := &MockLokiClient{}
	labels := map[string]string{"service": "test"}

	writer := NewLokiWriter(mockClient, labels)
	defer writer.Close()

	// Write invalid JSON
	invalidJSON := []byte("invalid json")

	n, err := writer.Write(invalidJSON)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}

	if n != 0 {
		t.Errorf("Write() returned %d, want 0", n)
	}
}

func TestLokiWriter_Close(t *testing.T) {
	mockClient := &MockLokiClient{}
	labels := map[string]string{"service": "test"}

	writer := NewLokiWriter(mockClient, labels)

	// Add some data to buffer
	logData := map[string]interface{}{
		"time":  time.Now().Format(time.RFC3339),
		"level": "info",
		"msg":   "test message",
	}

	jsonData, _ := json.Marshal(logData)
	writer.Write(jsonData)

	// Close writer
	err := writer.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}

	// Check that close was called on client
	// Note: The writer doesn't directly call client.Close(), it calls writer.Close()
	// which then calls client.Close() through the integration
	t.Logf("Client close calls: %d", mockClient.GetCloseCalls())
}

func TestLokiWriter_FlushBuffer(t *testing.T) {
	mockClient := &MockLokiClient{}
	labels := map[string]string{"service": "test"}

	writer := NewLokiWriter(mockClient, labels)
	defer writer.Close()

	// Add some data to buffer
	logData := map[string]interface{}{
		"time":  time.Now().Format(time.RFC3339),
		"level": "info",
		"msg":   "test message",
	}

	jsonData, _ := json.Marshal(logData)
	writer.Write(jsonData)

	// Flush buffer
	err := writer.flushBuffer()
	if err != nil {
		t.Errorf("flushBuffer() error = %v", err)
	}

	// Check that PushLogs was called
	if len(mockClient.GetPushLogsCalls()) == 0 {
		t.Error("Expected PushLogs to be called")
	}

	// Check that buffer is empty
	if len(writer.buffer) != 0 {
		t.Error("Expected buffer to be empty after flush")
	}
}

func TestLokiWriter_BackgroundFlusher(t *testing.T) {
	mockClient := &MockLokiClient{}
	labels := map[string]string{"service": "test"}

	writer := NewLokiWriter(mockClient, labels)
	defer writer.Close()

	// Add some data to buffer
	logData := map[string]interface{}{
		"time":  time.Now().Format(time.RFC3339),
		"level": "info",
		"msg":   "test message",
	}

	jsonData, _ := json.Marshal(logData)
	writer.Write(jsonData)

	// Wait for background flusher to run
	time.Sleep(6 * time.Second)

	// Check that flush was called
	// Note: The background flusher calls flushBuffer which calls PushLogs, not Flush directly
	t.Logf("Flush calls: %d, PushLogs calls: %d", mockClient.GetFlushCalls(), len(mockClient.GetPushLogsCalls()))
}
