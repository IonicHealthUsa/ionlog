package logcore

import (
	"path/filepath"
	"strings"
	"testing"
)

func BenchmarkDyFields(b *testing.B) {
	b.Run("GetCallerInfo", func(b *testing.B) {
		for range b.N {
			GetCallerInfo(1)
		}
	})
}

func TestGetCallerInfo(t *testing.T) {
	t.Run("should return current function information with skip=1", func(t *testing.T) {
		info := GetCallerInfo(1)

		// Check if file path ends with the correct test file name
		if info.File != "dynamic_fields_test.go" {
			t.Errorf("expected file to end with 'dynamic_fields_test.go', got %q", info.File)
		}

		// Check package name
		if info.PackageName != "logcore" {
			t.Errorf("expected package name 'logcore', got %q", info.PackageName)
		}

		// Check function name
		expectedFuncSuffix := "TestGetCallerInfo.func1"
		if info.Function != expectedFuncSuffix {
			t.Errorf("expected function name to end with %q, got %q", expectedFuncSuffix, info.Function)
		}

		// Line number is variable, just check if it's positive
		if info.Line <= 0 {
			t.Errorf("expected positive line number, got %d", info.Line)
		}
	})

	t.Run("should return caller's caller with skip=2", func(t *testing.T) {
		// Helper function to add a level to the call stack
		var helperCall = func() callerInfo {
			return GetCallerInfo(2) // Skip 2 levels to get the test function
		}

		info := helperCall()

		// Check if file path ends with the correct test file name
		if info.File != "dynamic_fields_test.go" {
			t.Errorf("expected file to end with 'dynamic_fields_test.go', got %q", info.File)
		}

		// Check package name
		if info.PackageName != "logcore" {
			t.Errorf("expected package name 'logcore', got %q", info.PackageName)
		}

		// Check function name - should be the test function name
		expectedFuncSuffix := "TestGetCallerInfo.func2"
		if info.Function != expectedFuncSuffix {
			t.Errorf("expected function name to end with %q, got %q", expectedFuncSuffix, info.Function)
		}
	})

	t.Run("should properly parse package and function names", func(t *testing.T) {
		info := GetCallerInfo(1)

		// Package name should not contain dots
		if strings.Contains(info.PackageName, ".") {
			t.Errorf("package name should not contain dots, got %q", info.PackageName)
		}

		// Function name might contain dots for method calls or nested functions
		// but should at least be non-empty
		if info.Function == "" {
			t.Errorf("function name should not be empty")
		}
	})
}

func exampleFunction() callerInfo {
	return GetCallerInfo(1)
}

func TestGetCallerInfoInDifferentFile(t *testing.T) {
	t.Run("should work when called from different functions", func(t *testing.T) {
		info := exampleFunction()

		// File path should point to the current file
		expectedFilename := filepath.Base(info.File)
		if expectedFilename != "dynamic_fields_test.go" {
			t.Errorf("expected file name 'dynamic_fields_test.go', got %q", expectedFilename)
		}

		if info.PackageName != "logcore" {
			t.Errorf("expected package name 'logcore', got %q", info.PackageName)
		}

		if info.Function != "exampleFunction" {
			t.Errorf("expected function name 'exampleFunction', got %q", info.Function)
		}
	})
}
