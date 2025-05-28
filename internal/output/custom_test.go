package output

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/IonicHealthUsa/ionlog/internal/core/logengine"
	"github.com/IonicHealthUsa/ionlog/internal/core/runtimeinfo"
)

func TestWrite(t *testing.T) {
	r := logengine.Report{
		Time:       time.Now().Format(time.RFC3339),
		Level:      logengine.Info,
		Msg:        "Hello World",
		CallerInfo: runtimeinfo.GetCallerInfo(1),
	}

	reportLog := fmt.Sprintf(`"time":"%s","level":"%s","msg":"%s","file":"%s","package":"%s","function":"%s","line":"%d"}
`, r.Time, r.Level, r.Msg, r.CallerInfo.File, r.CallerInfo.Package, r.CallerInfo.Function, r.CallerInfo.Line)

	t.Run("should write slice of byte on stdout", func(t *testing.T) {
		l, err := CustomOutput.Write([]byte(reportLog))
		if err != nil {
			t.Errorf("expected no error, but got %q", err)
		}
		if l != 0 {
			t.Errorf("expected report log to be %q, but got %q", 0, l)
		}
	})
}

func TestProcessLogline(t *testing.T) {
	t.Run("should return nil when line is nil", func(t *testing.T) {
		if format := processLogLine(nil); format != nil {
			t.Errorf("expected nil slice of byte, but got %q", format)
		}
	})

	t.Run("should return nil when could not decode the json", func(t *testing.T) {
		line := []byte(`"key":"value"`)

		if format := processLogLine(line); format != nil {
			t.Errorf("expected nil slice of byte, but got %q", format)
		}
	})

	t.Run("should return the correct format for each level type", func(t *testing.T) {
		testCase := [...]struct {
			report          logengine.Report
			reportLog       string
			expectFormatLog string
		}{
			{
				report: logengine.Report{
					Time:       time.Now().Format(time.RFC3339),
					Level:      logengine.Debug,
					Msg:        "Hello World",
					CallerInfo: runtimeinfo.GetCallerInfo(1),
				},
			},
			{
				report: logengine.Report{
					Time:       time.Now().Format(time.RFC3339),
					Level:      logengine.Info,
					Msg:        "Hello World",
					CallerInfo: runtimeinfo.GetCallerInfo(1),
				},
			},
			{
				report: logengine.Report{
					Time:       time.Now().Format(time.RFC3339),
					Level:      logengine.Warn,
					Msg:        "Hello World",
					CallerInfo: runtimeinfo.GetCallerInfo(1),
				},
			},
			{
				report: logengine.Report{
					Time:       time.Now().Format(time.RFC3339),
					Level:      logengine.Error,
					Msg:        "Hello World",
					CallerInfo: runtimeinfo.GetCallerInfo(1),
				},
			},
			{
				report: logengine.Report{
					Time:       time.Now().Format(time.RFC3339),
					Level:      logengine.Fatal,
					Msg:        "Hello World",
					CallerInfo: runtimeinfo.GetCallerInfo(1),
				},
			},
			{
				report: logengine.Report{
					Time:       time.Now().Format(time.RFC3339),
					Level:      logengine.Panic,
					Msg:        "Hello World",
					CallerInfo: runtimeinfo.GetCallerInfo(1),
				},
			},
			{
				report: logengine.Report{
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

				tt.expectFormatLog = fmt.Sprintf("%s %s [%s %s] %s (%s:%d%s)\n",
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

				gotLog := processLogLine([]byte(tt.reportLog))

				if !reflect.DeepEqual([]byte(tt.expectFormatLog), gotLog) {
					t.Errorf("expected log to be %q, but got %q", tt.expectFormatLog, gotLog)
				}
			})
		}
	})
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
