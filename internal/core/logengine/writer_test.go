package logengine

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

// MockWriter is a writer implementation for testing
type MockWriter struct {
	WriteFunc func(p []byte) (int, error)
}

func (m *MockWriter) Write(p []byte) (int, error) {
	if m.WriteFunc != nil {
		return m.WriteFunc(p)
	}
	return len(p), nil
}

// ErrorWriter always returns an error on write
type ErrorWriter struct {
	Err error
}

func (e *ErrorWriter) Write(p []byte) (int, error) {
	return 0, e.Err
}

func TestNewWriter(t *testing.T) {
	t.Run("Creates new writer with empty writers slice", func(t *testing.T) {
		w := NewWriter()

		// Verify type assertion
		_, ok := w.(*ionWriter)
		if !ok {
			t.Errorf("NewWriter() did not return a *ionWriter")
		}

		// Test writing to empty writer
		_, err := w.Write([]byte("test"))
		if err != nil {
			t.Errorf("Write to empty writer should not return error: got %v", err)
		}
	})
}

func TestSetWriters(t *testing.T) {
	t.Run("Sets multiple writers", func(t *testing.T) {
		w := NewWriter().(*ionWriter)
		buf1 := &bytes.Buffer{}
		buf2 := &bytes.Buffer{}

		w.SetWriters(buf1, buf2)

		if len(w.writers) != 2 {
			t.Errorf("Expected 2 writers, got %d", len(w.writers))
		}

		if w.writers[0] != buf1 || w.writers[1] != buf2 {
			t.Errorf("Writers not set correctly")
		}
	})

	t.Run("Replaces existing writers", func(t *testing.T) {
		w := NewWriter().(*ionWriter)
		buf1 := &bytes.Buffer{}

		// Set initial writers
		w.SetWriters(&bytes.Buffer{}, &bytes.Buffer{})

		// Replace with new writers
		w.SetWriters(buf1)

		if len(w.writers) != 1 {
			t.Errorf("Expected 1 writer after replacement, got %d", len(w.writers))
		}

		if w.writers[0] != buf1 {
			t.Errorf("Writers not replaced correctly")
		}
	})

	t.Run("Sets empty writers slice", func(t *testing.T) {
		w := NewWriter().(*ionWriter)

		// Set initial writers
		w.SetWriters(&bytes.Buffer{}, &bytes.Buffer{})

		// Replace with empty slice
		w.SetWriters()

		if len(w.writers) != 0 {
			t.Errorf("Expected empty writers slice, got %d writers", len(w.writers))
		}
	})
}

func TestAddWriter(t *testing.T) {
	t.Run("Adds writer to empty slice", func(t *testing.T) {
		w := NewWriter().(*ionWriter)
		buf := &bytes.Buffer{}

		w.AddWriter(buf)

		if len(w.writers) != 1 {
			t.Errorf("Expected 1 writer, got %d", len(w.writers))
		}

		if w.writers[0] != buf {
			t.Errorf("Writer not added correctly")
		}
	})

	t.Run("Adds writer to existing writers", func(t *testing.T) {
		w := NewWriter().(*ionWriter)
		buf1 := &bytes.Buffer{}
		buf2 := &bytes.Buffer{}

		w.AddWriter(buf1)
		w.AddWriter(buf2)

		if len(w.writers) != 2 {
			t.Errorf("Expected 2 writers, got %d", len(w.writers))
		}

		if w.writers[0] != buf1 || w.writers[1] != buf2 {
			t.Errorf("Writers not added correctly")
		}
	})
}

func TestWrite(t *testing.T) {
	// Capture stderr output for testing
	oldStderr := os.Stderr
	defer func() { os.Stderr = oldStderr }()

	t.Run("Writes to all writers", func(t *testing.T) {
		w := NewWriter().(*ionWriter)
		buf1 := &bytes.Buffer{}
		buf2 := &bytes.Buffer{}

		w.SetWriters(buf1, buf2)

		testData := []byte("test data")
		n, err := w.Write(testData)

		if err != nil {
			t.Errorf("Write returned error: %v", err)
		}

		if n != 0 {
			t.Errorf("Expected 0 bytes written, got %d", n)
		}

		if buf1.String() != string(testData) {
			t.Errorf("Data not written to first buffer correctly: got %q, want %q", buf1.String(), testData)
		}

		if buf2.String() != string(testData) {
			t.Errorf("Data not written to second buffer correctly: got %q, want %q", buf2.String(), testData)
		}
	})

	t.Run("Handles writer errors and continues", func(t *testing.T) {
		r, w, _ := os.Pipe()
		os.Stderr = w

		writer := NewWriter().(*ionWriter)
		buf := &bytes.Buffer{}
		errWriter := &ErrorWriter{Err: errors.New("write error")}

		writer.SetWriters(buf, errWriter)

		testData := []byte("test data")
		_, _ = writer.Write(testData)

		// Close the pipe writer to read from the pipe
		w.Close()

		// Read the stderr output
		errOutput := make([]byte, 1024)
		n, _ := r.Read(errOutput)
		errString := string(errOutput[:n])

		if !strings.Contains(errString, "Failed to write to in the 2° target") {
			t.Errorf("Expected error message for failed writer, got: %s", errString)
		}

		// Verify the successful writer still received the data
		if buf.String() != string(testData) {
			t.Errorf("Data not written to successful buffer: got %q, want %q", buf.String(), testData)
		}
	})

	t.Run("Handles nil writers", func(t *testing.T) {
		r, w, _ := os.Pipe()
		os.Stderr = w

		writer := NewWriter().(*ionWriter)
		buf := &bytes.Buffer{}

		// Set a nil writer
		writer.SetWriters(buf, nil)

		testData := []byte("test data")
		_, _ = writer.Write(testData)

		// Close the pipe writer to read from the pipe
		w.Close()

		// Read the stderr output
		errOutput := make([]byte, 1024)
		n, _ := r.Read(errOutput)
		errString := string(errOutput[:n])

		if !strings.Contains(errString, "Expected the 2° target to be not nil") {
			t.Errorf("Expected error message for nil writer, got: %s", errString)
		}

		// Verify the successful writer still received the data
		if buf.String() != string(testData) {
			t.Errorf("Data not written to successful buffer: got %q, want %q", buf.String(), testData)
		}
	})

	t.Run("Write lock prevents concurrent access", func(t *testing.T) {
		w := NewWriter().(*ionWriter)

		// Create a writer that blocks until signaled
		blockCh := make(chan struct{})

		blockingWriter := &MockWriter{
			WriteFunc: func(p []byte) (int, error) {
				// Block until signaled
				<-blockCh
				return len(p), nil
			},
		}

		var bufMutex sync.Mutex
		var buf string

		normalWriter := &MockWriter{
			WriteFunc: func(p []byte) (int, error) {
				bufMutex.Lock()
				defer bufMutex.Unlock()

				buf += string(p)
				return len(p), nil
			},
		}

		w.SetWriters(normalWriter, blockingWriter)

		go func() {
			w.Write([]byte("test1")) // blocked by second write
		}()
		time.Sleep(10 * time.Millisecond) // At least the normal writer should have written by now

		bufMutex.Lock()
		if buf != "test1" {
			t.Errorf("First write did not complete: got %q", buf)
		}
		bufMutex.Unlock()

		go func() {
			w.Write([]byte("test2"))
		}()
		time.Sleep(10 * time.Millisecond) // Same wait, but not expecting the second write to complete

		bufMutex.Lock()
		if strings.Contains(buf, "test2") {
			t.Errorf("Second write completed before first write")
		}
		bufMutex.Unlock()

		// Signal the blocking writer to continue
		close(blockCh)
		time.Sleep(10 * time.Millisecond) // Same wait, but now expecting the second write to complete

		bufMutex.Lock()
		if buf != "test1test2" {
			t.Errorf("Second write did not complete: got %q", buf)
		}
		bufMutex.Unlock()
	})
}

func TestInterface(t *testing.T) {
	t.Run("Implements IWriter interface", func(t *testing.T) {
		var _ IWriter = &ionWriter{}
	})

	t.Run("Implements io.Writer interface", func(t *testing.T) {
		var _ io.Writer = &ionWriter{}
	})
}
