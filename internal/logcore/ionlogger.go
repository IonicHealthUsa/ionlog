// Package ioncore provides the core functionalities of the logger.
// It is responsible for handling the logger service, the log writer, and the log engine.
package logcore

import (
	"context"
	"io"
	"log/slog"
	"sync"
	"time"

	"github.com/IonicHealthUsa/ionlog/internal/infrastructure/memory"
	"github.com/IonicHealthUsa/ionlog/internal/ionservice"
	"github.com/IonicHealthUsa/ionlog/internal/logrotation"
)

type service struct {
	ctx             context.Context
	cancel          context.CancelFunc
	serviceWg       sync.WaitGroup
	incomingReports bool
	serviceStatus   ionservice.ServiceStatus
}

type ionLogger struct {
	service

	logsMemory memory.IRecordMemory
	logRotate  logrotation.ILogRotation

	logEngine     *slog.Logger
	writerHandler ionWriter
	reports       chan ionReport
}

type IIonLogger interface {
	ionservice.IService

	LogsMemory() memory.IRecordMemory

	SetLogRotationSettings(folder string, maxFolderSize uint, rotation logrotation.PeriodicRotation)

	SetReportsBufferSizer(size uint)

	LogEngine() *slog.Logger
	SetLogEngine(handler *slog.Logger)

	Targets() []io.Writer
	SetTargets(targets ...io.Writer)

	CreateDefaultLogHandler() slog.Handler
	SendReport(r ionReport)
}

const maxReports = 100

var logger *ionLogger

func init() {
	logger = newLogger()

	// using internaly
	slog.SetDefault(slog.New(slog.NewTextHandler(DefaultOutput, &slog.HandlerOptions{Level: slog.LevelDebug})))
}

func newLogger() *ionLogger {
	l := &ionLogger{}
	l.ctx, l.cancel = context.WithCancel(context.Background())
	l.reports = make(chan ionReport, maxReports)
	l.logEngine = slog.New(l.CreateDefaultLogHandler())

	l.logsMemory = memory.NewRecordMemory()
	l.logRotate = nil

	return l
}

// Logger returns the logger instance
func Logger() IIonLogger {
	return logger
}

func (i *ionLogger) LogsMemory() memory.IRecordMemory {
	return i.logsMemory
}

// SetLogRotationSettings is a proxy to the log rotation service
func (i *ionLogger) SetLogRotationSettings(folder string, maxFolderSize uint, rotation logrotation.PeriodicRotation) {
	i.logRotate = logrotation.NewLogFileRotation()
	i.logRotate.SetLogRotationSettings(folder, maxFolderSize, rotation)
	i.SetTargets(append(i.Targets(), i.logRotate)...)
}

func (i *ionLogger) SetReportsBufferSizer(size uint) {
	i.reports = make(chan ionReport, size)
}

func (i *ionLogger) LogEngine() *slog.Logger {
	return i.logEngine
}

func (i *ionLogger) SetLogEngine(handler *slog.Logger) {
	i.logEngine = handler
}

func (i *ionLogger) Targets() []io.Writer {
	return i.writerHandler.writeTargets
}

func (i *ionLogger) SetTargets(targets ...io.Writer) {
	i.writerHandler.SetTargets(targets...)
}

// CreateDefaultLogHandler creates a default log handler for the logger
func (i *ionLogger) CreateDefaultLogHandler() slog.Handler {
	return slog.NewJSONHandler(
		&i.writerHandler,
		&slog.HandlerOptions{Level: slog.LevelDebug},
	)
}

// SendReport sends the report to the Logger handler, it will wait for 10ms before returning.
func (i *ionLogger) SendReport(r ionReport) {
	if i.incomingReports {
		return
	}

	if len(i.reports) == maxReports {
		return
	}

	select {
	case i.reports <- r:
		return

	case <-time.After(1 * time.Millisecond):
		return
	}

}

func (i *ionLogger) syncReports() {
	i.incomingReports = true
	for len(i.reports) > 0 {
		r := <-i.reports
		i.log(r.level, r.msg, r.args...)
	}
	i.incomingReports = false
}

func (i *ionLogger) log(level slog.Level, msg string, args ...any) {
	switch level {
	case slog.LevelDebug:
		i.logEngine.Debug(msg, args...)
	case slog.LevelInfo:
		i.logEngine.Info(msg, args...)
	case slog.LevelWarn:
		i.logEngine.Warn(msg, args...)
	case slog.LevelError:
		i.logEngine.Error(msg, args...)

	default:
		slog.Warn("Unknown log level")
	}
}
