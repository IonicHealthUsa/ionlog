package styles

import (
	"encoding/json"
	"fmt"
	"maps"
	"reflect"
	"testing"
	"time"

	"github.com/IonicHealthUsa/ionlog/internal/core/logengine"
	"github.com/IonicHealthUsa/ionlog/internal/core/runtimeinfo"
)

func TestWrite(t *testing.T) {
	r := logengine.ReportType{
		Time:       time.Now().Format(time.RFC3339),
		Level:      logengine.Info,
		Msg:        "Hello World",
		CallerInfo: runtimeinfo.GetCallerInfo(1),
	}

	reportLog := fmt.Sprintf(`{"time":"%s","level":"%s","msg":"%s","file":"%s","package":"%s","function":"%s","line":"%d"}
`, r.Time, r.Level, r.Msg, r.CallerInfo.File, r.CallerInfo.Package, r.CallerInfo.Function, r.CallerInfo.Line)

	t.Run("should write slice of byte on stdout", func(t *testing.T) {
		processedLog, err := processLogLine([]byte(reportLog))
		if err != nil {
			t.Errorf("expected no error, but got %q", err)
		}

		l, err := CustomOutput.Write([]byte(reportLog))
		if err != nil {
			t.Errorf("expected no error, but got %q", err)
		}
		if l != len(processedLog) {
			t.Errorf("expected report log to be %v, but got %v", len(processedLog), l)
		}
	})
}

func TestProcessLogline(t *testing.T) {
	t.Run("should return nil when line is nil", func(t *testing.T) {
		format, err := processLogLine(nil)
		if err == nil {
			t.Errorf("expected an error when line is nil, but got nil")
		}

		if format != nil {
			t.Errorf("expected nil slice of byte, but got %q", format)
		}
	})

	t.Run("should return nil when could not decode the json", func(t *testing.T) {
		line := []byte(`"key":"value"`)

		log, err := processLogLine(line)
		if err == nil {
			t.Errorf("expected an error when decoding json, but got nil")
		}

		if log != nil {
			t.Errorf("expected nil slice of byte, but got %q", log)
		}
	})

	t.Run("should return the correct format for each level type", func(t *testing.T) {
		testCase := [...]struct {
			report          logengine.ReportType
			reportLog       string
			expectFormatLog string
		}{
			{
				report: logengine.ReportType{
					Time:       time.Now().Format(time.RFC3339),
					Level:      logengine.Debug,
					Msg:        "Hello World",
					CallerInfo: runtimeinfo.GetCallerInfo(1),
				},
			},
			{
				report: logengine.ReportType{
					Time:       time.Now().Format(time.RFC3339),
					Level:      logengine.Info,
					Msg:        "Hello World",
					CallerInfo: runtimeinfo.GetCallerInfo(1),
				},
			},
			{
				report: logengine.ReportType{
					Time:       time.Now().Format(time.RFC3339),
					Level:      logengine.Warn,
					Msg:        "Hello World",
					CallerInfo: runtimeinfo.GetCallerInfo(1),
				},
			},
			{
				report: logengine.ReportType{
					Time:       time.Now().Format(time.RFC3339),
					Level:      logengine.Error,
					Msg:        "Hello World",
					CallerInfo: runtimeinfo.GetCallerInfo(1),
				},
			},
			{
				report: logengine.ReportType{
					Time:       time.Now().Format(time.RFC3339),
					Level:      logengine.Fatal,
					Msg:        "Hello World",
					CallerInfo: runtimeinfo.GetCallerInfo(1),
				},
			},
			{
				report: logengine.ReportType{
					Time:       time.Now().Format(time.RFC3339),
					Level:      logengine.Panic,
					Msg:        "Hello World",
					CallerInfo: runtimeinfo.GetCallerInfo(1),
				},
			},
			{
				report: logengine.ReportType{
					Time:       time.Now().Format(time.RFC3339),
					Level:      logengine.Trace,
					Msg:        "Hello World",
					CallerInfo: runtimeinfo.GetCallerInfo(1),
				},
			},
		}

		for _, tt := range testCase {
			t.Run(tt.report.Level.String(), func(t *testing.T) {
				timestamp := formatTimestamp(tt.report.Time)
				levelColor := getLevelColor(tt.report.Level.String())
				functionName := formatFunctionName(tt.report.CallerInfo.Function)

				tt.expectFormatLog = fmt.Sprintf("%s %s [%s %s] %s (%s:%d%s) \n",
					bold+white+timestamp+reset,
					levelColor+tt.report.Level.String()+reset,

					cyan+tt.report.CallerInfo.Package+reset,
					functionName,

					levelColor+tt.report.Msg+reset,

					magenta+tt.report.CallerInfo.File,
					tt.report.CallerInfo.Line, reset,
				)

				tt.reportLog = fmt.Sprintf(`{"time":"%s","level":"%s","msg":"%s","file":"%s","package":"%s","function":"%s","line":"%d"}
`, tt.report.Time, tt.report.Level, tt.report.Msg, tt.report.CallerInfo.File, tt.report.CallerInfo.Package, tt.report.CallerInfo.Function, tt.report.CallerInfo.Line)

				gotLog, err := processLogLine([]byte(tt.reportLog))
				if err != nil {
					t.Errorf("expected no error, but got %q", err)
				}

				if !reflect.DeepEqual([]byte(tt.expectFormatLog), gotLog) {
					t.Errorf("expected log to be %q, but got %q", tt.expectFormatLog, gotLog)
				}
			})
		}
	})

	t.Run("should return the correct format with static fields", func(t *testing.T) {
		report := logengine.ReportType{
			Time:       time.Now().Format(time.RFC3339),
			Level:      logengine.Info,
			Msg:        "Hello World",
			CallerInfo: runtimeinfo.GetCallerInfo(1),
		}

		var entry logEntry
		reportLog := fmt.Sprintf(`{"test":"123","time":"%s","level":"%s","msg":"%s","file":"%s","package":"%s","function":"%s","line":"%d"}
`, report.Time, report.Level, report.Msg, report.CallerInfo.File, report.CallerInfo.Package, report.CallerInfo.Function, report.CallerInfo.Line)
		if err := json.Unmarshal([]byte(reportLog), &entry); err != nil {
			t.Errorf("expected no error, but got %q", err)
		}

		staticFieldMap := map[string]string{"test": "123"}
		maps.Copy(entry, staticFieldMap)

		timestamp := formatTimestamp(report.Time)
		levelColor := getLevelColor(report.Level.String())
		functionName := formatFunctionName(report.CallerInfo.Function)
		staticField := formatStaticField(entry)

		expectFormatLog := fmt.Sprintf("%s %s [%s %s] %s (%s:%d%s) %s\n",
			bold+white+timestamp+reset,
			levelColor+report.Level.String()+reset,

			cyan+report.CallerInfo.Package+reset,
			functionName,

			levelColor+report.Msg+reset,

			magenta+report.CallerInfo.File,
			report.CallerInfo.Line, reset,
			staticField,
		)

		gotLog, err := processLogLine([]byte(reportLog))
		if err != nil {
			t.Errorf("expected no error, but got %q", err)
		}

		if !reflect.DeepEqual([]byte(expectFormatLog), gotLog) {
			t.Errorf("expected log to be %q, but got %q", expectFormatLog, gotLog)
		}
	})
}

func BenchmarkProcessLogLine(b *testing.B) {
	report := logengine.ReportType{
		Time:       time.Now().Format(time.RFC3339),
		Level:      logengine.Info,
		Msg:        "Hello World",
		CallerInfo: runtimeinfo.GetCallerInfo(1),
	}

	var entry logEntry
	reportLog := fmt.Sprintf(`{"test":"123","ionic":"health","time":"%s","level":"%s","msg":"%s","file":"%s","package":"%s","function":"%s","line":"%d"}
`, report.Time, report.Level, report.Msg, report.CallerInfo.File, report.CallerInfo.Package, report.CallerInfo.Function, report.CallerInfo.Line)
	if err := json.Unmarshal([]byte(reportLog), &entry); err != nil {
		b.Errorf("expected no error, but got %q", err)
	}

	b.ResetTimer()

	for range b.N {
		_, _ = processLogLine([]byte(reportLog))
	}
}

func TestFormatStaticField(t *testing.T) {
	t.Run("should return empty when does not exist static fields", func(t *testing.T) {
		expectedFormatStaticFields := ""
		entry := map[string]string{
			"time":     "17-06-2025",
			"level":    "DEBUG",
			"msg":      "test 123",
			"file":     "format.go",
			"package":  "styles",
			"function": "formatStaticField",
			"line":     "123",
		}

		gotFormatStaticFields := formatStaticField(entry)

		if gotFormatStaticFields != expectedFormatStaticFields {
			t.Errorf("expcted the static field to be %q, but got %q", expectedFormatStaticFields, gotFormatStaticFields)
		}
	})

	t.Run("should return the static field as string", func(t *testing.T) {
		expectedFormatStaticFields := "computer-id:q "
		entry := map[string]string{
			"time":        "17-06-2025",
			"level":       "DEBUG",
			"msg":         "test 123",
			"file":        "format.go",
			"package":     "styles",
			"function":    "formatStaticField",
			"line":        "123",
			"computer-id": "q",
		}

		gotFormatStaticFields := formatStaticField(entry)

		if gotFormatStaticFields != expectedFormatStaticFields {
			t.Errorf("expcted the static field to be %q, but got %q", expectedFormatStaticFields, gotFormatStaticFields)
		}
	})
}

func BenchmarkFormatStaticField(b *testing.B) {
	entry := map[string]string{
		"time":     "17-06-2025",
		"level":    "DEBUG",
		"msg":      "test 123",
		"file":     "format.go",
		"packet":   "styles",
		"function": "formatStaticField",
		"line":     "123",
		"static1":  "field1",
		"static2":  "field2",
		"static3":  "field3",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = formatStaticField(entry)
	}
}

func TestFormatTimestamp(t *testing.T) {
	t.Run("should return the correct timestamp", func(t *testing.T) {
		timeNow := time.Now()
		timeStr := timeNow.Format(time.RFC3339)

		expectedTimeStr, err := time.Parse(time.RFC3339Nano, timeStr)
		if err != nil {
			t.Errorf("expected no error to parse the time, but got %q", err)
		}

		if format := formatTimestamp(timeStr); format != expectedTimeStr.Format(time.RFC3339) {
			t.Errorf("expected time format to be %q, but got %q", expectedTimeStr.Format(time.RFC3339), format)
		}
	})

	t.Run("should return the arg timeStr", func(t *testing.T) {
		timeStr := "123456789.ABC"

		if format := formatTimestamp(timeStr); format != timeStr {
			t.Errorf("expected time format to be %q, but got %q", timeStr, format)
		}
	})
}

func BenchmarkFormatTimestamp(b *testing.B) {
	timeNow := time.Now()
	timeStr := timeNow.Format(time.RFC3339)

	b.ResetTimer()

	for range b.N {
		_ = formatTimestamp(timeStr)
	}
}

func TestFormatFunctionName(t *testing.T) {
	t.Run("should return the last function", func(t *testing.T) {
		function := "func1.func2.func3"
		expectedFormat := blue + "func3" + reset

		if format := formatFunctionName(function); format != expectedFormat {
			t.Errorf("expected format of function to be %q, but got %q", expectedFormat, format)
		}
	})

	t.Run("should return the correct function name format", func(t *testing.T) {
		function := "func1"
		expectedFormat := blue + "func1" + reset

		if format := formatFunctionName(function); format != expectedFormat {
			t.Errorf("expected format of function to be %q, but got %q", expectedFormat, format)
		}
	})
}

func BenchmarkFormatFunctionName(b *testing.B) {
	function := "func1.func2.func3"

	b.ResetTimer()

	for range b.N {
		_ = formatFunctionName(function)
	}
}

func TestGetLevelColor(t *testing.T) {
	testCase := [...]struct {
		level         string
		expectedColor string
	}{
		{
			level:         "DEBUG",
			expectedColor: white,
		},
		{
			level:         "INFO",
			expectedColor: green,
		},
		{
			level:         "WARN",
			expectedColor: yellow,
		},
		{
			level:         "ERROR",
			expectedColor: red,
		},
		{
			level:         "FATAL",
			expectedColor: bgRed + bold + white,
		},
		{
			level:         "PANIC",
			expectedColor: bgRed + bold + white,
		},
		{
			level:         "TRACE",
			expectedColor: cyan,
		},
		{
			level:         "others",
			expectedColor: reset,
		},
	}

	t.Run("should return the correct color for each level type", func(t *testing.T) {
		for _, tt := range testCase {
			if color := getLevelColor(tt.level); color != tt.expectedColor {
				t.Errorf("expected the color of %q to be %q, but got %q", tt.level, tt.expectedColor, color)
			}
		}

	})
}

func BenchmarkGetLevelColor(b *testing.B) {
	b.ResetTimer()

	for range b.N {
		_ = getLevelColor("INFO")
	}
}
