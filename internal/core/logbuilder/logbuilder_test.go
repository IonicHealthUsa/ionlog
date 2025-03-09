package logbuilder

import (
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

	b.Run("log builder String function", func(b *testing.B) {
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

// func TestNewLogBuilder(t *testing.T) {
// 	t.Run("should create a new log builder with initialized fields map", func(t *testing.T) {
// 		lb := NewLogBuilder()
//
// 		if lb == nil {
// 			t.Fatal("expected non-nil logBuilder")
// 		}
//
// 		if lb.fields == nil {
// 			t.Fatal("expected initialized fields map")
// 		}
//
// 		if len(lb.fields) != 0 {
// 			t.Errorf("expected empty fields map, got map with %d entries", len(lb.fields))
// 		}
//
// 		if lb.base != "" {
// 			t.Errorf("expected empty base string, got %q", lb.base)
// 		}
// 	})
// }
//
// func TestAddField(t *testing.T) {
// 	t.Run("should add a field to the logBuilder", func(t *testing.T) {
// 		lb := NewLogBuilder()
// 		lb.AddField("key", "value")
//
// 		if len(lb.fields) != 1 {
// 			t.Fatalf("expected 1 field, got %d", len(lb.fields))
// 		}
//
// 		if lb.fields["key"] != "value" {
// 			t.Errorf("expected value %q for key %q, got %q", "value", "key", lb.fields["key"])
// 		}
// 	})
//
// 	t.Run("should replace existing field with same key", func(t *testing.T) {
// 		lb := NewLogBuilder()
// 		lb.AddField("key", "value1")
// 		lb.AddField("key", "value2")
//
// 		if len(lb.fields) != 1 {
// 			t.Fatalf("expected 1 field, got %d", len(lb.fields))
// 		}
//
// 		if lb.fields["key"] != "value2" {
// 			t.Errorf("expected value %q for key %q, got %q", "value2", "key", lb.fields["key"])
// 		}
// 	})
// }
//
// func TestString(t *testing.T) {
// 	t.Run("should return empty JSON object for empty logBuilder", func(t *testing.T) {
// 		lb := NewLogBuilder()
// 		result := lb.String()
// 		expected := "{}\n"
//
// 		if result != expected {
// 			t.Errorf("expected %q, got %q", expected, result)
// 		}
// 	})
//
// 	t.Run("should properly format single field as JSON", func(t *testing.T) {
// 		lb := NewLogBuilder()
// 		lb.AddField("key", "value")
// 		result := lb.String()
// 		expected := "{\"key\":\"value\"}\n"
//
// 		if result != expected {
// 			t.Errorf("expected %q, got %q", expected, result)
// 		}
// 	})
//
// 	t.Run("should properly format multiple fields as JSON", func(t *testing.T) {
// 		lb := NewLogBuilder()
// 		lb.AddField("key1", "value1")
// 		lb.AddField("key2", "value2")
//
// 		result := lb.String()
//
// 		// Since map iteration order is not guaranteed in Go, we need to check if the result
// 		// contains all expected parts without relying on the exact order
// 		if !strings.Contains(result, "\"key1\":\"value1\"") {
// 			t.Errorf("expected result to contain %q", "\"key1\":\"value1\"")
// 		}
//
// 		if !strings.Contains(result, "\"key2\":\"value2\"") {
// 			t.Errorf("expected result to contain %q", "\"key2\":\"value2\"")
// 		}
//
// 		if !strings.HasPrefix(result, "{") {
// 			t.Errorf("expected result to start with %q", "{")
// 		}
//
// 		if !strings.HasSuffix(result, "}\n") {
// 			t.Errorf("expected result to end with %q", "}\n")
// 		}
//
// 		// Check that we have exactly one comma (for two fields)
// 		commaCount := strings.Count(result, ",")
// 		if commaCount != 1 {
// 			t.Errorf("expected result to contain exactly 1 comma, got %d", commaCount)
// 		}
// 	})
//
// 	t.Run("should handle special characters in values", func(t *testing.T) {
// 		lb := NewLogBuilder()
// 		lb.AddField("key", "value with spaces")
// 		result := lb.String()
// 		expected := "{\"key\":\"value with spaces\"}\n"
//
// 		if result != expected {
// 			t.Errorf("expected %q, got %q", expected, result)
// 		}
// 	})
//
// 	t.Run("should handle special characters in keys", func(t *testing.T) {
// 		lb := NewLogBuilder()
// 		lb.AddField("key-with-hyphens", "value")
// 		result := lb.String()
// 		expected := "{\"key-with-hyphens\":\"value\"}\n"
//
// 		if result != expected {
// 			t.Errorf("expected %q, got %q", expected, result)
// 		}
// 	})
// }
//
// func TestCombinedOperations(t *testing.T) {
// 	t.Run("should maintain consistent state through multiple operations", func(t *testing.T) {
// 		lb := NewLogBuilder()
//
// 		// Add initial fields
// 		lb.AddField("key1", "value1")
// 		lb.AddField("key2", "value2")
//
// 		// Verify initial string output
// 		initialResult := lb.String()
// 		if !strings.Contains(initialResult, "\"key1\":\"value1\"") {
// 			t.Errorf("expected initial result to contain %q", "\"key1\":\"value1\"")
// 		}
// 		if !strings.Contains(initialResult, "\"key2\":\"value2\"") {
// 			t.Errorf("expected initial result to contain %q", "\"key2\":\"value2\"")
// 		}
//
// 		// Check comma count (should be 1 for 2 fields)
// 		fisrtCommaCount := strings.Count(initialResult, ",")
// 		if fisrtCommaCount != 1 {
// 			t.Errorf("expected updated result to contain exactly 1 commas, got %d", fisrtCommaCount)
// 		}
//
// 		// Modify an existing field
// 		lb.AddField("key1", "updated-value")
//
// 		// Add a new field
// 		lb.AddField("key3", "value3")
//
// 		// Verify updated string output
// 		updatedResult := lb.String()
// 		if !strings.Contains(updatedResult, "\"key1\":\"updated-value\"") {
// 			t.Errorf("expected updated result to contain %q", "\"key1\":\"updated-value\"")
// 		}
// 		if !strings.Contains(updatedResult, "\"key2\":\"value2\"") {
// 			t.Errorf("expected updated result to contain %q", "\"key2\":\"value2\"")
// 		}
// 		if !strings.Contains(updatedResult, "\"key3\":\"value3\"") {
// 			t.Errorf("expected updated result to contain %q", "\"key3\":\"value3\"")
// 		}
//
// 		// Check comma count (should be 2 for 3 fields)
// 		commaCount := strings.Count(updatedResult, ",")
// 		if commaCount != 2 {
// 			t.Errorf("expected updated result to contain exactly 2 commas, got %d", commaCount)
// 		}
// 	})
// }
