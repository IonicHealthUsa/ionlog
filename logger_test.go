package ionlog

import (
	"sync"
	"testing"

	"github.com/IonicHealthUsa/ionlog/internal/service"
)

func TestStart(t *testing.T) {
	// Ensure logger is stopped before and after test
	Stop()
	defer Stop()

	t.Run("should start the logger service", func(t *testing.T) {
		Start()

		lock.RLock()
		status := logger.Status()
		lock.RUnlock()

		if status != service.Running {
			t.Errorf("expected status to be Running, but got %v", status)
		}
	})

	t.Run("should allow multiple concurrent calls to Start", func(t *testing.T) {
		const concurrentCalls = 100
		var wg sync.WaitGroup
		wg.Add(concurrentCalls)

		for i := 0; i < concurrentCalls; i++ {
			go func() {
				defer wg.Done()
				Start()
			}()
		}

		wg.Wait()

		lock.RLock()
		status := logger.Status()
		lock.RUnlock()

		if status != service.Running {
			t.Errorf("expected status to be Running, but got %v", status)
		}
	})

	t.Run("should handle Start after Stop", func(t *testing.T) {
		Start()
		Stop()
		Start()

		lock.RLock()
		status := logger.Status()
		lock.RUnlock()

		if status != service.Running {
			t.Errorf("expected status to be Running after restart, but got %v", status)
		}
	})
}

func TestStartConcurrent_Stress(t *testing.T) {
	// Stress test with many concurrent calls to Start and Stop
	const iterations = 10
	const concurrentCalls = 50

	for i := 0; i < iterations; i++ {
		var wg sync.WaitGroup
		wg.Add(concurrentCalls)

		for j := 0; j < concurrentCalls; j++ {
			go func() {
				defer wg.Done()
				Start()
			}()
		}

		// Wait for all Start calls
		wg.Wait()

		// Stop and wait for reset
		Stop()

		// Final Start to ensure recovery
		Start()

		lock.RLock()
		status := logger.Status()
		lock.RUnlock()

		if status != service.Running {
			t.Errorf("iteration %d: expected status to be Running, but got %v", i, status)
		}
	}
}
