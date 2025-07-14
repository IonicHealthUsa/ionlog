package styles

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)

// customWriter type of customs writers
type customWriter struct{}

// Logentry logs in JSON format
type logEntry map[string]string

// ANSI color for terminal
const (
	reset    = "\033[0m"
	bold     = "\033[1m"
	red      = "\033[31m"
	green    = "\033[32m"
	yellow   = "\033[33m"
	blue     = "\033[34m"
	magenta  = "\033[35m"
	cyan     = "\033[36m"
	white    = "\033[37m"
	bgRed    = "\033[41m"
	bgGreen  = "\033[42m"
	bgYellow = "\033[43m"
	bgBlue   = "\033[44m"
)

func (c *customWriter) Write(p []byte) (int, error) {
	log, err := processLogLine(p)
	if err != nil {
		return 0, fmt.Errorf("failed to process log line: %w", err)
	}

	return os.Stdout.Write(log)
}

var (
	CustomOutput = &customWriter{}
)

var logEntryKeyDefault = []string{"time", "level", "msg", "file", "package", "function", "line"}

func processLogLine(line []byte) ([]byte, error) {
	if line == nil {
		return nil, ErrNilLine
	}

	var entry logEntry
	err := json.Unmarshal(line, &entry)
	if err != nil {
		return nil, err
	}

	timestamp := formatTimestamp(entry["time"])
	functionName := formatFunctionName(entry["function"])
	levelColor := getLevelColor(entry["level"])
	staticField := formatStaticField(entry)

	formatLine := fmt.Sprintf("%s%s%s%s %s%s%s [%s%s%s %s] %s%s%s (%s%s:%s%s) %s\n",
		bold, white, timestamp, reset,
		levelColor, entry["level"], reset,

		cyan, entry["package"], reset,
		functionName,

		levelColor, entry["msg"], reset,

		magenta, entry["file"],
		entry["line"], reset,

		staticField,
	)

	return []byte(formatLine), nil
}

func formatTimestamp(timeStr string) string {
	t, err := time.Parse(time.RFC3339Nano, timeStr)
	if err != nil {
		return timeStr
	}
	return t.Format(time.RFC3339)
}

func formatFunctionName(function string) string {
	parts := strings.Split(function, ".")
	if len(parts) > 1 {
		return blue + parts[len(parts)-1] + reset
	}
	return blue + function + reset
}

func getLevelColor(level string) string {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return white
	case "INFO":
		return green
	case "WARN":
		return yellow
	case "ERROR":
		return red
	case "FATAL", "PANIC":
		return bgRed + bold + white
	case "TRACE":
		return cyan
	default:
		return reset
	}
}
func formatStaticField(entry map[string]string) string {
	numStaticFields := len(entry) - len(logEntryKeyDefault)
	if numStaticFields == 0 {
		return ""
	}

	var staticField strings.Builder
	staticField.Grow(numStaticFields * 40) // expected 40 bytes for each static field
	for k, v := range entry {
		if !slices.Contains(logEntryKeyDefault, k) {
			staticField.WriteString(k)
			staticField.WriteString(":")
			staticField.WriteString(v)
			staticField.WriteString(" ")
		}
	}

	return staticField.String()
}
