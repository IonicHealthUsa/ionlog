package logcore

import (
	"testing"
)

func TestNewLogger(t *testing.T) {
	_logger := newLogger()

	if _logger == nil {
		t.Errorf("Expected logger to be not nil")
	}
}
