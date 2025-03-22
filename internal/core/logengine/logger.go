package logengine

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/IonicHealthUsa/ionlog/internal/core/logbuilder"
	"github.com/IonicHealthUsa/ionlog/internal/core/runtimeinfo"
	"github.com/IonicHealthUsa/ionlog/internal/infrastructure/memory"
)

type Report struct {
	Time       string
	Level      Level
	Msg        string
	CallerInfo runtimeinfo.CallerInfo
}

type logger struct {
	builder    logbuilder.ILogBuilder
	logsMemory memory.IRecordMemory
	closed     bool
	reports    chan Report
	writer     IWriter

	staticFields map[string]string
	traceMode    bool

	reportLock sync.Mutex
}

type ILogger interface {
	AsyncReport(r Report)
	Report(r Report)
	FlushReports()
	HandleReports(ctx context.Context)
	Writer() IWriter
	Memory() memory.IRecordMemory
	SetStaticFields(attrs map[string]string)
	SetReportQueueSize(size uint)
	SetTraceMode(mode bool)
	TraceMode() bool
}

func NewLogger() ILogger {
	logger := &logger{}

	logger.builder = logbuilder.NewLogBuilder()
	logger.logsMemory = memory.NewRecordMemory()
	logger.reports = make(chan Report, 100)
	logger.writer = NewWriter()

	return logger
}

func (l *logger) AsyncReport(r Report) {
	if l.closed {
		return
	}
	select {
	case l.reports <- r:
	case <-time.After(1 * time.Millisecond):
	}
}

func (l *logger) Report(r Report) {
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

	l.writer.Write(l.builder.Compile())
}

func (l *logger) FlushReports() {
	for len(l.reports) > 0 {
		l.Report(<-l.reports)
	}
}

func (l *logger) HandleReports(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			l.closed = true
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

func (l *logger) SetStaticFields(attrs map[string]string) {
	l.staticFields = attrs
}

func (l *logger) SetReportQueueSize(size uint) {
	l.reports = make(chan Report, size)
}

func (l *logger) SetTraceMode(mode bool) {
	l.traceMode = mode
}

func (l *logger) TraceMode() bool {
	return l.traceMode
}
