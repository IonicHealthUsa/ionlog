package logbuilder

import (
	"reflect"
	"strconv"
	"testing"
	"time"
)

var fakeMessage = "We shall not cease from exploration and the end of all our exploring will be to arrive where we started and know the place for the first time."

type CallerInfo struct {
	File        string
	PackageName string
	Function    string
	Line        int
}

func BenchmarkStaticFields(b *testing.B) {
	l2 := NewLogBuilder()

	b.Run("log builder compile function", func(b *testing.B) {
		for range b.N {
			var callerInfo CallerInfo

			l2.AddFields(
				"time", time.Now().Format(time.RFC3339),
				"level", "INFO",
				"msg", fakeMessage,
				"file", callerInfo.File,
				"package", callerInfo.PackageName,
				"function", callerInfo.Function,
				"line", strconv.Itoa(callerInfo.Line),
			)
			_ = l2.Compile()
		}
	})
}

func TestNewLogBuilder(t *testing.T) {
	t.Run("Creates instance with initialized buffer", func(t *testing.T) {
		lb := NewLogBuilder()

		// Check type assertion
		_lb, ok := lb.(*logBuilder)
		if !ok {
			t.Errorf("NewLogBuilder() did not return a *logBuilder")
		}

		b := make([]byte, bufsize)
		b[0] = byte('{')

		if cap(_lb.buf) != bufsize {
			t.Errorf("expected the size of buffer to be %v, but got %v", bufsize, cap(_lb.buf))
		}
		if !reflect.DeepEqual(_lb.buf, b) {
			t.Errorf("expected the byte to be %q, but got %q", b, _lb.buf)
		}
		if _lb.p != 1 {
			t.Errorf("expected the pointer of buffer to be %v, but got %v", 1, _lb.p)
		}

		// Compile to check initial state
		result := lb.Compile()
		expected := []byte("{}\n")

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Initial buffer state incorrect, got: %s, want: %s", result, expected)
		}
	})
}

func TestWriteString(t *testing.T) {
	t.Run("should transform the string to byte", func(t *testing.T) {
		lb := NewLogBuilder()
		_lb, ok := lb.(*logBuilder)
		if !ok {
			t.Errorf("NewLogBuilder() did not return a *logBuilder")
		}

		str := "hello world"
		_lb.writeString(str)
		b := []byte("{" + str + "\n")

		if reflect.DeepEqual(_lb.buf[:_lb.p], b) {
			t.Errorf("expected the buffer to be %q, but got %q", b, _lb.buf[:_lb.p])
		}
	})
}

func TestResetBuff(t *testing.T) {
	t.Run("should reset the buffer pointer and overwrite the buffer", func(t *testing.T) {
		lb := NewLogBuilder()
		_lb, ok := lb.(*logBuilder)
		if !ok {
			t.Errorf("NewLogBuilder() did not return a *logBuilder")
		}

		b := []byte("hello world")
		copy(_lb.buf, b)
		_lb.p = uint(len(b))

		_lb.resetBuff()

		if _lb.p != 1 {
			t.Errorf("expected the buffer pointer to be %q, but got %q", 1, _lb.p)
		}

		if _lb.buf[0] != byte('{') {
			t.Errorf("expected the byte to be %q, but got %q", byte('{'), _lb.buf[0])
		}
	})
}

func TestAddFields(t *testing.T) {
	t.Run("Adds single field", func(t *testing.T) {
		lb := NewLogBuilder()
		lb.AddFields("key", "value")

		result := lb.Compile()
		expected := []byte("{\"key\":\"value\"}\n")

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("AddFields single field incorrect, got: %s, want: %s", result, expected)
		}
	})

	t.Run("Adds multiple fields", func(t *testing.T) {
		lb := NewLogBuilder()
		lb.AddFields("key1", "value1", "key2", "value2")

		result := lb.Compile()
		expected := []byte("{\"key1\":\"value1\",\"key2\":\"value2\"}\n")

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("AddFields multiple fields incorrect, got: %s, want: %s", result, expected)
		}
	})

	t.Run("Handles odd number of arguments", func(t *testing.T) {
		lb := NewLogBuilder()
		lb.AddFields("key1", "value1", "key2")

		result := lb.Compile()
		expected := []byte("{}\n")

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("AddFields with odd arguments incorrect, got: %s, want: %s", result, expected)
		}
	})

	t.Run("Handles empty arguments", func(t *testing.T) {
		lb := NewLogBuilder()
		lb.AddFields()

		result := lb.Compile()
		expected := []byte("{}\n")

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("AddFields with empty arguments incorrect, got: %s, want: %s", result, expected)
		}
	})

	t.Run("Multiple AddFields calls", func(t *testing.T) {
		lb := NewLogBuilder()
		lb.AddFields("key1", "value1")
		lb.AddFields("key2", "value2")

		result := lb.Compile()
		expected := []byte("{\"key1\":\"value1\",\"key2\":\"value2\"}\n")

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Multiple AddFields calls incorrect, got: %s, want: %s", result, expected)
		}
	})
}

func TestCompile(t *testing.T) {
	t.Run("Returns correct JSON format", func(t *testing.T) {
		lb := NewLogBuilder()
		lb.AddFields("key", "value")

		result := lb.Compile()
		expected := []byte("{\"key\":\"value\"}\n")

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Compile format incorrect, got: %s, want: %s", result, expected)
		}
	})

	t.Run("Resets buffer after compile", func(t *testing.T) {
		lb := NewLogBuilder()
		lb.AddFields("key1", "value1")
		lb.Compile()

		// Add new fields after compile
		lb.AddFields("key2", "value2")

		result := lb.Compile()
		expected := []byte("{\"key2\":\"value2\"}\n")

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Buffer reset after compile incorrect, got: %s, want: %s", result, expected)
		}
	})

	t.Run("Empty JSON when no fields added", func(t *testing.T) {
		lb := NewLogBuilder()

		result := lb.Compile()
		expected := []byte("{}\n")

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Empty JSON format incorrect, got: %s, want: %s", result, expected)
		}
	})
}

func TestLogBuilderSequence(t *testing.T) {
	t.Run("Full sequence of operations", func(t *testing.T) {
		lb := NewLogBuilder()

		// First set of fields
		lb.AddFields("request_id", "abc123", "user_id", "user456")
		result1 := lb.Compile()
		expected1 := []byte("{\"request_id\":\"abc123\",\"user_id\":\"user456\"}\n")

		if !reflect.DeepEqual(result1, expected1) {
			t.Errorf("First sequence incorrect, got: %s, want: %s", result1, expected1)
		}

		// Second set of fields
		lb.AddFields("event", "login", "status", "success")
		result2 := lb.Compile()
		expected2 := []byte("{\"event\":\"login\",\"status\":\"success\"}\n")

		if !reflect.DeepEqual(result2, expected2) {
			t.Errorf("Second sequence incorrect, got: %s, want: %s", result2, expected2)
		}
	})
}

func TestSpecialCharacters(t *testing.T) {
	t.Run("Handles special characters in keys and values", func(t *testing.T) {
		lb := NewLogBuilder()
		lb.AddFields("key-with-dash", "value",
			"key_with_underscore", "value:with:colons",
			"emoji", "😀🔥")

		result := lb.Compile()
		expected := []byte("{\"key-with-dash\":\"value\",\"key_with_underscore\":\"value:with:colons\",\"emoji\":\"😀🔥\"}\n")

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Special characters handling incorrect, got: %s, want: %s", result, expected)
		}
	})
}

func TestBufLimits(t *testing.T) {
	t.Run("Buffer expands when full", func(t *testing.T) {
		lb := NewLogBuilder()
		_lb, ok := lb.(*logBuilder)
		if !ok {
			t.Errorf("NewLogBuilder() did not return a *logBuilder")
		}

		// Fill buffer to max
		for len(_lb.buf) < maxBufsize {
			lb.AddFields("key", "value")
		}

		if len(_lb.buf) != maxBufsize {
			t.Error("Buffer size is not the expected size")
		}

		// add one more field
		lb.AddFields("key", "value")

		if len(_lb.buf) != maxBufsize {
			t.Error("Buffer size is not the expected size")
		}
	})
}

func TestWriteByte(t *testing.T) {
	t.Run("should write the byte on buffer", func(t *testing.T) {
		lb := NewLogBuilder()
		_lb, ok := lb.(*logBuilder)
		if !ok {
			t.Errorf("NewLogBuilder() did not return a *logBuilder")
		}

		b := byte('r')
		_lb.writeByte(b)

		if _lb.buf[_lb.p-1] != b {
			t.Errorf("exepcted the value of %q on buffer to be %q, bug got %q", _lb.p-1, b, _lb.buf[_lb.p-1])
		}
	})
}
