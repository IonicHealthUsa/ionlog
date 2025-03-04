package usecases

import (
	"testing"

	"github.com/IonicHealthUsa/ionlog/internal/infrastructure/memory"
)

func TestLogOnce(t *testing.T) {
	t.Run("First Log", func(t *testing.T) {
		r := memory.NewRecordMemory()
		if !LogOnce(r, "pkg", "function", "file", "msg") {
			t.Errorf("LogOnce() failed")
		}
	})

	t.Run("Two logs check same logs msg", func(t *testing.T) {
		r := memory.NewRecordMemory()

		if !LogOnce(r, "pkg", "function", "file", "msg") {
			t.Errorf("LogOnce() failed")
		}

		if LogOnce(r, "pkg", "function", "file", "msg") {
			t.Errorf("LogOnce() failed")
		}
	})

	t.Run("Two logs but new msg", func(t *testing.T) {
		r := memory.NewRecordMemory()

		if !LogOnce(r, "pkg", "function", "file", "msg") {
			t.Errorf("LogOnce() failed")
		}

		if !LogOnce(r, "pkg", "function", "file", "New Msg") {
			t.Errorf("LogOnce() failed")
		}
	})
}
