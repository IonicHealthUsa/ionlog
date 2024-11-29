package ionlog

import (
	"fmt"
	"io"
	"log/slog"
)

type ionLogger struct {
	logHandler    *slog.Logger
	writerHandler ionWriter
}

type customAttrs func(i *ionLogger)

var logger *ionLogger

// init initializes the logger with default values
func init() {
	logger = &ionLogger{}
	logger.logHandler = slog.New(createDefaultLogHandler())
}

func createDefaultLogHandler() slog.Handler {
	return slog.NewJSONHandler(
		&logger.writerHandler,
		&slog.HandlerOptions{Level: slog.LevelDebug},
	)
}

// SetLogAttributes sets the log SetLogAttributes
// fns is a variadic parameter that accepts customAttrs
func SetLogAttributes(fns ...customAttrs) {
	for _, fn := range fns {
		fn(logger)
	}
}

func WithTargets(w ...io.Writer) customAttrs {
	return func(i *ionLogger) {
		i.writerHandler.writeTargets = w
	}
}

func WithAttrs(attrs map[string]string) customAttrs {
	return func(i *ionLogger) {
		index := 0
		slogAttrs := make([]slog.Attr, len(attrs))
		for k, v := range attrs {
			slogAttrs[index] = slog.String(k, v)
			index++
		}
		handler := createDefaultLogHandler().WithAttrs(slogAttrs)
		i.logHandler = slog.New(handler)
	}
}

func Info(msg string, args ...any) {
	logger.logHandler.Info(fmt.Sprintf(msg, args...), getRecordInformation()...)
}

func Error(msg string, args ...any) {
	logger.logHandler.Error(fmt.Sprintf(msg, args...), getRecordInformation()...)
}

func Warn(msg string, args ...any) {
	logger.logHandler.Warn(fmt.Sprintf(msg, args...), getRecordInformation()...)
}

func Debug(msg string, args ...any) {
	logger.logHandler.Debug(fmt.Sprintf(msg, args...), getRecordInformation()...)
}
