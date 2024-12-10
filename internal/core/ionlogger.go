package ioncore

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"time"

	ionlogfile "github.com/IonicHealthUsa/ionlog/internal/logfile"
	recordhistory "github.com/IonicHealthUsa/ionlog/internal/record_history"
	ionservice "github.com/IonicHealthUsa/ionlog/internal/service"
)

type controlFlow struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

type autoRotateInfo struct {
	logRotateService ionlogfile.ILogFileRotation
	rotationPeriod   ionlogfile.PeriodicRotation
	folder           string
}

type ionLogger struct {
	controlFlow
	autoRotateInfo

	history       *recordhistory.RecordHistory
	logEngine     *slog.Logger
	writerHandler ionWriter
	reports       chan ionReport
	serviceStatus ionservice.ServiceStatus
}

type IIonLogger interface {
	ionservice.IService

  History() *recordhistory.RecordHistory

	SetRotationPeriod(period ionlogfile.PeriodicRotation)
	SetFolder(folder string)

	LogEngine() *slog.Logger
	SetLogEngine(handler *slog.Logger)

	Targets() []io.Writer
	SetTargets(targets ...io.Writer)

	CreateDefaultLogHandler() slog.Handler
	SendReport(r ionReport)
}

const timeout = 10 * time.Millisecond

var logger *ionLogger

func init() {
	logger = newLogger()

	// using internaly
	slog.SetDefault(slog.New(slog.NewJSONHandler(DefaultOutput, &slog.HandlerOptions{Level: slog.LevelDebug})))
}

func newLogger() *ionLogger {
	l := &ionLogger{}
	l.ctx, l.cancel = context.WithCancel(context.Background())
	l.reports = make(chan ionReport)
	l.logEngine = slog.New(l.CreateDefaultLogHandler())
	l.rotationPeriod = ionlogfile.NoAutoRotate
	l.history = recordhistory.NewRecordHistory()

	return l
}

// Logger returns the logger instance
func Logger() IIonLogger {
	return logger
}

func (i *ionLogger) SetRotationPeriod(period ionlogfile.PeriodicRotation) {
	i.rotationPeriod = period
}

func (i *ionLogger) History() *recordhistory.RecordHistory {
	return i.history
}

func (i *ionLogger) SetFolder(folder string) {
	i.folder = folder
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
	select {
	case <-time.After(timeout):
		slog.Warn(fmt.Sprintf("Failed to send the report (timeout=%v): %v", timeout, r))
		return
	case i.reports <- r:
	}
}

func (i *ionLogger) handleIonReports() {
	for {
		select {
		case <-i.ctx.Done():
			slog.Debug("Logger stopped by context")
			return

		case r := <-i.reports:
			i.log(r.level, r.msg, r.args...)
		}
	}
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
