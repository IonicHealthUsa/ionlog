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

	t.Run("Adds two writers to empty slice", func(t *testing.T) {
		w := NewWriter().(*ionWriter)
		buf1 := &bytes.Buffer{}
		buf2 := &bytes.Buffer{}

		w.AddWriter(buf1, buf2)

		if len(w.writers) != 2 {
			t.Errorf("Expected 2 writer, got %d", len(w.writers))
		}

		if w.writers[0] != buf1 || w.writers[1] != buf2 {
			t.Errorf("Writers not added correctly")
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

	t.Run("Adds two writer to existing writers", func(t *testing.T) {
		w := NewWriter().(*ionWriter)
		buf1 := &bytes.Buffer{}
		buf2 := &bytes.Buffer{}
		buf3 := &bytes.Buffer{}

		w.AddWriter(buf1)
		w.AddWriter(buf2, buf3)

		if len(w.writers) != 3 {
			t.Errorf("Expected 3 writers, but got %d", len(w.writers))
		}

		if w.writers[0] != buf1 || w.writers[1] != buf2 || w.writers[2] != buf3 {
			t.Error("Writers not added correctly")
		}
	})

	t.Run("should timeout when mutex is lock", func(t *testing.T) {
		w := NewWriter().(*ionWriter)
		buf := &bytes.Buffer{}

		w.writeLock.Lock()
		go func() {
			w.AddWriter(buf)
		}()
		time.Sleep(10 * time.Millisecond)

		if len(w.writers) != 0 {
			t.Errorf("expected empty writers slice, but got %q writers", len(w.writers))
		}

		w.writeLock.Unlock()
		time.Sleep(10 * time.Millisecond)
		w.writeLock.Lock()
		if len(w.writers) != 1 {
			t.Errorf("expected one writer slice, but got %q writers", len(w.writers))
		}

		if w.writers[0] != buf {
			t.Errorf("writers not replaced correctly")
		}
		w.writeLock.Unlock()
	})
}

func TestDeleteWriter(t *testing.T) {
	t.Run("should delete a writer on the list", func(t *testing.T) {
		w := NewWriter().(*ionWriter)
		buf1 := &bytes.Buffer{}
		buf2 := &bytes.Buffer{}

		w.writers = append(w.writers, buf1, buf2)

		if len(w.writers) != 2 {
			t.Errorf("expected the size of writers to be 2, but got %q", len(w.writers))
		}

		if w.writers[0] != buf1 || w.writers[1] != buf2 {
			t.Errorf("Writers not added correctly")
			return
		}

		w.DeleteWriter(buf2)

		if len(w.writers) != 1 {
			t.Errorf("expected the size of writers to be 1, but got %q", len(w.writers))
		}

		w.DeleteWriter(buf1)

		if len(w.writers) != 0 {
			t.Errorf("expected the size of writers to be 0, but got %q", len(w.writers))
		}
	})

	t.Run("should delete all writer on the list", func(t *testing.T) {
		w := NewWriter().(*ionWriter)
		buf1 := &bytes.Buffer{}
		buf2 := &bytes.Buffer{}

		w.writers = append(w.writers, buf1, buf2)

		if len(w.writers) != 2 {
			t.Errorf("expected the size of writers to be 2, but got %q", len(w.writers))
		}

		if w.writers[0] != buf1 || w.writers[1] != buf2 {
			t.Errorf("Writers not added correctly")
			return
		}

		w.DeleteWriter(buf1, buf2)

		if len(w.writers) != 0 {
			t.Errorf("expected the size of writers to be 0, but got %d", len(w.writers))
		}
	})

	t.Run("should not delete anyone writers on the list", func(t *testing.T) {
		w := NewWriter().(*ionWriter)
		buf1 := &bytes.Buffer{}
		buf2 := &bytes.Buffer{}
		buf3 := &bytes.Buffer{}

		w.writers = append(w.writers, buf1, buf2)

		if len(w.writers) != 2 {
			t.Errorf("expected the size of writers to be 2, but got %q", len(w.writers))
		}

		if w.writers[0] != buf1 || w.writers[1] != buf2 {
			t.Errorf("Writers not added correctly")
			return
		}

		w.DeleteWriter(buf3)

		if len(w.writers) != 2 {
			t.Errorf("expected the size of writers to be 0, but got %d", len(w.writers))
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

		w.AddWriter(buf1, buf2)

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

		writer.AddWriter(buf, errWriter)

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
		writer.AddWriter(buf, nil)

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

		w.AddWriter(normalWriter, blockingWriter)

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
