// Package ioncore provides the core functionalities of the logger.
// It is responsible for handling the logger service, the log writer, and the log engine.
package logcore

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"time"

	"github.com/IonicHealthUsa/ionlog/internal/infrastructure/memory"
	"github.com/IonicHealthUsa/ionlog/internal/interfaces"
	"github.com/IonicHealthUsa/ionlog/internal/logrotation"
)

type controlFlow struct {
	ctx                   context.Context
	cancel                context.CancelFunc
	servicesWg            sync.WaitGroup
	reportsWg             sync.WaitGroup
	blockIncommingReports bool
}

type autoRotateInfo struct {
	logRotateService logrotation.ILogFileRotation
	rotationPeriod   logrotation.PeriodicRotation
	folder           string
	maxFolderSize    uint
}

type ionLogger struct {
	controlFlow
	autoRotateInfo

	logsMemory    memory.IRecordMemory
	logEngine     *slog.Logger
	writerHandler ionWriter
	reports       chan ionReport
	serviceStatus interfaces.ServiceStatus
}

type IIonLogger interface {
	interfaces.IService

	LogsMemory() memory.IRecordMemory

	SetRotationPeriod(period logrotation.PeriodicRotation)
	SetFolder(folder string)
	SetFolderMaxSize(folderMaxSize uint)

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
	l.reports = make(chan ionReport, 10000)
	l.logEngine = slog.New(l.CreateDefaultLogHandler())
	l.rotationPeriod = logrotation.NoAutoRotate
	l.logsMemory = memory.NewRecordMemory()

	return l
}

// Logger returns the logger instance
func Logger() IIonLogger {
	return logger
}

func (i *ionLogger) SetFolderMaxSize(folderMaxSize uint) {
	i.autoRotateInfo.maxFolderSize = folderMaxSize
}

func (i *ionLogger) SetRotationPeriod(period logrotation.PeriodicRotation) {
	i.rotationPeriod = period
}

func (i *ionLogger) LogsMemory() memory.IRecordMemory {
	return i.logsMemory
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
	if i.blockIncommingReports {
		return
	}
	i.reportsWg.Add(1)
	select {
	case <-time.After(timeout):
		slog.Warn(fmt.Sprintf("Failed to send the report (timeout=%v): %v", timeout, r))
		i.reportsWg.Done() // Will not be processed, so decrement the counter.
		return
	case i.reports <- r:
	}
}

// handleIonReports handles the reports sent to the logger
// When the context is canceled, it will log all the reports in the queue before returning
func (i *ionLogger) handleIonReports() {
	for {
		select {
		case <-i.ctx.Done():
			slog.Debug("Logger stopped by context")
			i.blockIncommingReports = true
			for len(i.reports) > 0 {
				r := <-i.reports
				i.reportsWg.Done()
				i.log(r.level, r.msg, r.args...)
			}
			return

		case r := <-i.reports:
			i.reportsWg.Done()
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
