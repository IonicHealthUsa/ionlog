package logengine

import (
	"context"
	"fmt"
	"maps"
	"os"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/IonicHealthUsa/ionlog/internal/core/logbuilder"
	"github.com/IonicHealthUsa/ionlog/internal/core/runtimeinfo"
	"github.com/IonicHealthUsa/ionlog/internal/infrastructure/memory"
)

type ReportType struct {
	Time       string
	Level      Level
	Msg        string
	CallerInfo runtimeinfo.CallerInfo
}

type logger struct {
	builder    logbuilder.ILogBuilder
	logsMemory memory.IRecordMemory
	closed     bool
	reports    chan ReportType
	writer     IWriter

	staticFields map[string]string
	traceMode    bool

	reportLock sync.Mutex
	closeLock  sync.Mutex
}

type ILogger interface {
	AsyncReport(r ReportType)
	Report(r ReportType)
	FlushReports()
	HandleReports(ctx context.Context)
	Writer() IWriter
	Memory() memory.IRecordMemory
	AddStaticFields(attrs map[string]string)
	DeleteStaticField(fields ...string)
	SetReportQueueSize(size uint)
	SetTraceMode(mode bool)
	TraceMode() bool
}

func NewLogger() ILogger {
	logger := &logger{}

	logger.builder = logbuilder.NewLogBuilder()
	logger.logsMemory = memory.NewRecordMemory()
	logger.reports = make(chan ReportType, 100)
	logger.writer = NewWriter()

	return logger
}

func (l *logger) closeReport() {
	l.closeLock.Lock()
	defer l.closeLock.Unlock()
	l.closed = true
}

func (l *logger) getStatusCloseReport() bool {
	l.closeLock.Lock()
	defer l.closeLock.Unlock()
	return l.closed
}

func (l *logger) AsyncReport(r ReportType) {
	if l.getStatusCloseReport() {
		return
	}
	select {
	case l.reports <- r:
	case <-time.After(1 * time.Second):
		fmt.Fprintf(os.Stderr, "logger reports channel is full\n")
	}
}

func (l *logger) Report(r ReportType) {
	l.reportLock.Lock()
	defer l.reportLock.Unlock()

	if l.staticFields != nil {
		for key, value := range l.staticFields {
			l.builder.AddFields(key, value)
		}
	}

	l.builder.AddFields(
		"time", r.Time,
		"level", r.Level.String(),
		"msg", r.Msg,
		"file", r.CallerInfo.File,
		"package", r.CallerInfo.Package,
		"function", r.CallerInfo.Function,
		"line", strconv.Itoa(r.CallerInfo.Line),
	)

	_, _ = l.writer.Write(l.builder.Compile())
}

func (l *logger) FlushReports() {
	for {
		select {
		case r := <-l.reports:
			l.Report(r)

		case <-time.After(1 * time.Millisecond):
			return
		}
	}
}

func (l *logger) HandleReports(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			l.closeReport()
			return

		case r := <-l.reports:
			l.Report(r)
		}
	}
}

func (l *logger) Writer() IWriter {
	return l.writer
}

func (l *logger) Memory() memory.IRecordMemory {
	return l.logsMemory
}

func (l *logger) AddStaticFields(attrs map[string]string) {
	l.reportLock.Lock()
	defer l.reportLock.Unlock()

	if l.staticFields == nil {
		l.staticFields = attrs
		return
	}

	maps.Copy(l.staticFields, attrs)
}

func (l *logger) DeleteStaticField(fields ...string) {
	l.reportLock.Lock()
	defer l.reportLock.Unlock()

	maps.DeleteFunc(l.staticFields, func(k string, v string) bool {
		return slices.Contains(fields, k)
	})
}

func (l *logger) SetReportQueueSize(size uint) {
	l.reportLock.Lock()
	defer l.reportLock.Unlock()
	l.reports = make(chan ReportType, size)
}

func (l *logger) SetTraceMode(mode bool) {
	l.reportLock.Lock()
	defer l.reportLock.Unlock()
	l.traceMode = mode
}

func (l *logger) TraceMode() bool {
	l.reportLock.Lock()
	defer l.reportLock.Unlock()
	return l.traceMode
}
