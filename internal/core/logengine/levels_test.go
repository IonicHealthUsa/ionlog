package logengine

import (
	"testing"
)

func TestLevelConstants(t *testing.T) {
	t.Run("Trace constant has correct value", func(t *testing.T) {
		if Trace != -2 {
			t.Errorf("Trace constant has incorrect value: got %d, want %d", Trace, -2)
		}
	})

	t.Run("Debug constant has correct value", func(t *testing.T) {
		if Debug != -1 {
			t.Errorf("Debug constant has incorrect value: got %d, want %d", Debug, -1)
		}
	})

	t.Run("Info constant has correct value", func(t *testing.T) {
		if Info != 0 {
			t.Errorf("Info constant has incorrect value: got %d, want %d", Info, 0)
		}
	})

	t.Run("Warn constant has correct value", func(t *testing.T) {
		if Warn != 1 {
			t.Errorf("Warn constant has incorrect value: got %d, want %d", Warn, 1)
		}
	})

	t.Run("Error constant has correct value", func(t *testing.T) {
		if Error != 2 {
			t.Errorf("Error constant has incorrect value: got %d, want %d", Error, 2)
		}
	})

	t.Run("Panic constant has correct value", func(t *testing.T) {
		if Panic != 3 {
			t.Errorf("Panic constant has incorrect value: got %d, want %d", Panic, 3)
		}
	})

	t.Run("Fatal constant has correct value", func(t *testing.T) {
		if Fatal != 4 {
			t.Errorf("Fatal constant has incorrect value: got %d, want %d", Fatal, 4)
		}
	})
}

func TestLevelString(t *testing.T) {
	t.Run("Trace level returns correct string", func(t *testing.T) {
		if s := Trace.String(); s != "TRACE" {
			t.Errorf("Trace.String() returned incorrect value: got %s, want %s", s, "TRACE")
		}
	})

	t.Run("Debug level returns correct string", func(t *testing.T) {
		if s := Debug.String(); s != "DEBUG" {
			t.Errorf("Debug.String() returned incorrect value: got %s, want %s", s, "DEBUG")
		}
	})

	t.Run("Info level returns correct string", func(t *testing.T) {
		if s := Info.String(); s != "INFO" {
			t.Errorf("Info.String() returned incorrect value: got %s, want %s", s, "INFO")
		}
	})

	t.Run("Warn level returns correct string", func(t *testing.T) {
		if s := Warn.String(); s != "WARN" {
			t.Errorf("Warn.String() returned incorrect value: got %s, want %s", s, "WARN")
		}
	})

	t.Run("Error level returns correct string", func(t *testing.T) {
		if s := Error.String(); s != "ERROR" {
			t.Errorf("Error.String() returned incorrect value: got %s, want %s", s, "ERROR")
		}
	})

	t.Run("Panic level returns correct string", func(t *testing.T) {
		if s := Panic.String(); s != "PANIC" {
			t.Errorf("Panic.String() returned incorrect value: got %s, want %s", s, "PANIC")
		}
	})

	t.Run("Fatal level returns correct string", func(t *testing.T) {
		if s := Fatal.String(); s != "FATAL" {
			t.Errorf("Fatal.String() returned incorrect value: got %s, want %s", s, "FATAL")
		}
	})

	t.Run("Custom level returns numeric string", func(t *testing.T) {
		customLevel := Level(10)
		expected := "10"
		if s := customLevel.String(); s != expected {
			t.Errorf("Level(10).String() returned incorrect value: got %s, want %s", s, expected)
		}
	})

	t.Run("Negative custom level returns numeric string", func(t *testing.T) {
		customLevel := Level(-10)
		expected := "-10"
		if s := customLevel.String(); s != expected {
			t.Errorf("Level(-10).String() returned incorrect value: got %s, want %s", s, expected)
		}
	})
}
