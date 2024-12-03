package ionlog

import (
	"context"
	"fmt"
	"log/slog"
)

type ionLogger struct {
	logHandler    *slog.Logger
	writerHandler ionWriter
	ctx           context.Context
	reports       chan string
}

var logger *ionLogger = newLogger()

func newLogger() *ionLogger {
	l := &ionLogger{}

  l.reports = make(chan string)
	l.logHandler = slog.New(l.createDefaultLogHandler())

	return l
}

func (l *ionLogger) createDefaultLogHandler() slog.Handler {
	return slog.NewJSONHandler(
		&l.writerHandler,
		&slog.HandlerOptions{Level: slog.LevelDebug},
	)
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
