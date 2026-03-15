package ionlog

import (
	"bytes"
	"os"
	"testing"
)

func TestSetAttributes(t *testing.T) {
	// Ensure logger is reset
	Stop()
	defer Stop()

	t.Run("should apply multiple attributes", func(t *testing.T) {
		SetAttributes(
			WithTraceMode(true),
			WithQueueSize(200),
		)

		lock.RLock()
		defer lock.RUnlock()

		if !logger.LogEngine().TraceMode() {
			t.Error("expected trace mode to be enabled")
		}
	})
}

func TestWithWriters(t *testing.T) {
	Stop()
	defer Stop()

	t.Run("should add and remove writers", func(t *testing.T) {
		buf := &bytes.Buffer{}

		SetAttributes(WithWriters(buf))

		// Verify writer is added (indirectly by checking if We can write)
		// Internal writer access is needed or we just verify it doesn't panic
		SetAttributes(WithoutWriters(buf))
	})
}

func TestWithStaticFields(t *testing.T) {
	Stop()
	defer Stop()

	t.Run("should manage static fields", func(t *testing.T) {
		fields := map[string]string{"app": "test-app"}
		SetAttributes(WithStaticFields(fields))

		// Note: No direct way to check static fields without internal access
		// but we can verify it doesn't panic and we can remove them
		SetAttributes(WithoutStaticFields("app"))
	})
}

func TestWithLogFileRotation(t *testing.T) {
	Stop()
	defer Stop()

	t.Run("should enable rotation", func(t *testing.T) {
		folder := "test-logs"
		SetAttributes(WithLogFileRotation(folder, 10*1024, Daily))

		if _, err := os.Stat(folder); err != nil {
			t.Errorf("expected folder %s to be created", folder)
		}

		os.RemoveAll(folder)
	})
}

func TestWithQueueSize(t *testing.T) {
	Stop()
	defer Stop()

	t.Run("should set queue size", func(t *testing.T) {
		SetAttributes(WithQueueSize(500))
		// Verification would require internal access to the channel capacity
	})
}

func TestWithTraceMode(t *testing.T) {
	Stop()
	defer Stop()

	t.Run("should toggle trace mode", func(t *testing.T) {
		SetAttributes(WithTraceMode(true))

		lock.RLock()
		if !logger.LogEngine().TraceMode() {
			t.Error("expected trace mode to be true")
		}
		lock.RUnlock()

		SetAttributes(WithTraceMode(false))

		lock.RLock()
		if logger.LogEngine().TraceMode() {
			t.Error("expected trace mode to be false")
		}
		lock.RUnlock()
	})
}

func TestWithCallerInfoDepth(t *testing.T) {
	Stop()
	defer Stop()

	t.Run("should set caller info depth", func(t *testing.T) {
		const depth = 5
		SetAttributes(WithCallerInfoDepth(depth))

		lock.RLock()
		if logger.LogEngine().GetCallerStackDepth() != depth {
			t.Errorf("expected depth to be %d, but got %d", depth, logger.LogEngine().GetCallerStackDepth())
		}
		lock.RUnlock()
	})
}
