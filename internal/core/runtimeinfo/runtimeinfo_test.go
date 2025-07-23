package runtimeinfo

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
		if info.File != "runtimeinfo_test.go" {
			t.Errorf("expected file to end with 'runtimeinfo_test.go', got %q", info.File)
		}

		// Check package name
		if info.Package != "runtimeinfo" {
			t.Errorf("expected package name 'runtimeinfo', got %q", info.Package)
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
		var helperCall = func() CallerInfo {
			return GetCallerInfo(2) // Skip 2 levels to get the test function
		}

		info := helperCall()

		// Check if file path ends with the correct test file name
		if info.File != "runtimeinfo_test.go" {
			t.Errorf("expected file to end with 'runtimeinfo_test.go', got %q", info.File)
		}

		// Check package name
		if info.Package != "runtimeinfo" {
			t.Errorf("expected package name 'runtimeinfo', got %q", info.Package)
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
		if strings.Contains(info.Package, ".") {
			t.Errorf("package name should not contain dots, got %q", info.Package)
		}

		// Function name might contain dots for method calls or nested functions
		// but should at least be non-empty
		if info.Function == "" {
			t.Errorf("function name should not be empty")
		}
	})

	t.Run("should return empty caller info when skip is more than function on stack", func(t *testing.T) {
		info := GetCallerInfo(5)

		if info.File != "" {
			t.Errorf("expected file to be empty, but got %q", info.File)
		}

		if info.Package != "" {
			t.Errorf("expected package to be empty, but got %q", info.Package)
		}

		if info.Function != "" {
			t.Errorf("expected function to be empty, but got %q", info.Function)
		}

		if info.Line != 0 {
			t.Errorf("expected line to be 0, but got %q", info.Line)
		}
	})
}

func exampleFunction() CallerInfo {
	return GetCallerInfo(1)
}

func TestGetCallerInfoInDifferentFile(t *testing.T) {
	t.Run("should work when called from different functions", func(t *testing.T) {
		info := exampleFunction()

		// File path should point to the current file
		expectedFilename := filepath.Base(info.File)
		if expectedFilename != "runtimeinfo_test.go" {
			t.Errorf("expected file name 'runtimeinfo_test.go', got %q", expectedFilename)
		}

		if info.Package != "runtimeinfo" {
			t.Errorf("expected package name 'runtimeinfo', got %q", info.Package)
		}

		if info.Function != "exampleFunction" {
			t.Errorf("expected function name 'exampleFunction', got %q", info.Function)
		}
	})
}
