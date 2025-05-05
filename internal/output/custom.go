package output

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type customWriter struct{}

// Logentry logs in JSON format
type logEntry struct {
	ComputerID string `json:"computer-id"`
	Time       string `json:"time"`
	Level      string `json:"level"`
	Message    string `json:"msg"`
	File       string `json:"file"`
	Package    string `json:"package"`
	Function   string `json:"function"`
	Line       string `json:"line"`
}

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

var (
	CustomOutput = &customWriter{}
)

func (c *customWriter) Write(p []byte) (int, error) {
	return DefaultOutput.Write(processLogLine(p))
}

func processLogLine(line []byte) []byte {
	var entry logEntry

	if line == nil {
		return nil
	}

	err := json.Unmarshal(line, &entry)
	if err != nil {
		fmt.Printf("%sErro '%s' ao processar linha: %s%s\n", yellow, err.Error(), line, reset)
		return nil
	}

	timestamp := formatTimestamp(entry.Time)

	levelColor := getLevelColor(entry.Level)

	formatLine := fmt.Sprintf("%s %s [%s %s] %s (%s:%s)\n",
		bold+white+timestamp+reset,
		levelColor+entry.Level+reset,

		cyan+entry.Package+reset,
		formatFunctionName(entry.Function),

		levelColor+entry.Message+reset,

		magenta+entry.File,
		entry.Line+reset,
	)

	return []byte(formatLine)
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
	if len(parts) > 0 {
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
