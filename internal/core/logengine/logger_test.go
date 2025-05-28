package logengine

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/IonicHealthUsa/ionlog/internal/core/runtimeinfo"
)

func TestNewLogger(t *testing.T) {
	t.Run("should return logger instance", func(t *testing.T) {
		l := NewLogger()
		if l == nil {
			t.Error("NewLogger did not returned a interface of logger")
		}
		if reflect.ValueOf(l).IsNil() {
			t.Error("expected a value to logger")
		}

		_l, ok := l.(*logger)
		if !ok {
			t.Fatal("NewLogger did not returned a instance of logger")
		}

		if _l.builder == nil {
			t.Error("expected the momory was instance")
		}
		if reflect.ValueOf(_l.builder).IsNil() {
			t.Error("expected the builder was not nil")
		}

		if _l.logsMemory == nil {
			t.Error("expected the momory was instance")
		}
		if reflect.ValueOf(_l.logsMemory).IsNil() {
			t.Error("expected the momory was not nil")
		}

		if _l.reports == nil {
			t.Error("expected a chan to reports")
		}

		if _l.writer == nil {
			t.Error("expected the writes was instace")
		}
		if reflect.ValueOf(_l.writer).IsNil() {
			t.Error("expected the write was not nil")
		}
	})
}

func TestCloseReport(t *testing.T) {
	t.Run("should close the asynchronous report", func(t *testing.T) {
		l := NewLogger()
		_l, ok := l.(*logger)
		if !ok {
			t.Fatalf("NewLogger did not returned a instance of logger")
		}

		if _l.closed {
			t.Errorf("expected the closed report to be %v, but got %v", false, _l.closed)
		}

		_l.closeReport()

		if !_l.closed {
			t.Errorf("expected the closed report to be %v, but got %v", true, _l.closed)
		}
	})
}

func TestGetStatusCloseReport(t *testing.T) {
	t.Run("should get the status of reports closed", func(t *testing.T) {
		l := NewLogger()
		_l, ok := l.(*logger)
		if !ok {
			t.Fatalf("NewLogger did not returned a instance of logger")
		}

		if _l.getStatusCloseReport() {
			t.Errorf("expected the closed report to be %v, but got %v", false, _l.getStatusCloseReport())
		}

		_l.closed = true

		if !_l.getStatusCloseReport() {
			t.Errorf("expected the closed report to be %v, but got %v", true, _l.getStatusCloseReport())
		}
	})
}

func TestAsyncReport(t *testing.T) {
	r := Report{
		Time:       time.Now().Format(time.RFC3339),
		Level:      Info,
		Msg:        "Hello World",
		CallerInfo: runtimeinfo.GetCallerInfo(1),
	}

	t.Run("should receive the report", func(t *testing.T) {
		l := NewLogger()
		_l, ok := l.(*logger)
		if !ok {
			t.Fatalf("NewLogger did not returned a instance of logger")
		}

		l.AsyncReport(r)

		select {
		case report := <-_l.reports:
			if report.Time != r.Time {
				t.Errorf("expected time to be %q, but got %q", r.Time, report.Time)
			}
			if report.Level != r.Level {
				t.Errorf("expected level to be %q, but got %q", r.Level, report.Level)
			}
			if report.Msg != r.Msg {
				t.Errorf("expected message to be %q, but got %q", r.Msg, report.Msg)
			}
			if report.CallerInfo.File != r.CallerInfo.File {
				t.Errorf("expected file info to be %q, but got %q", r.CallerInfo.File, report.CallerInfo.File)
			}
			if report.CallerInfo.Line != r.CallerInfo.Line {
				t.Errorf("expected line info to be %q, but got %q", r.CallerInfo.Line, report.CallerInfo.Line)
			}
			if report.CallerInfo.Package != r.CallerInfo.Package {
				t.Errorf("expected package info to be %q, but got %q", r.CallerInfo.Package, report.CallerInfo.Package)
			}
			if report.CallerInfo.Function != r.CallerInfo.Function {
				t.Errorf("expected function info to be %q, but got %q", r.CallerInfo.Function, report.CallerInfo.Function)
			}
		case <-time.After(time.Second):
			t.Error("expected a report, but timeout")
		}
	})

	t.Run("should timeout when logger is closed", func(t *testing.T) {
		l := NewLogger()
		_l, ok := l.(*logger)
		if !ok {
			t.Fatalf("NewLogger did not returned a instance of logger")
		}

		_l.closed = true
		l.AsyncReport(r)

		select {
		case <-_l.reports:
			t.Error("expected no report, but got")
		case <-time.After(time.Second):
		}
	})

	t.Run("should timeout when report channel is full", func(t *testing.T) {
		l := NewLogger()
		_l, ok := l.(*logger)
		if !ok {
			t.Fatalf("NewLogger did not returned a instance of logger")
		}
		_l.reports = make(chan Report)

		l.AsyncReport(r)

		select {
		case <-_l.reports:
			t.Error("expected no report, but got")
		case <-time.After(time.Second):
		}
	})
}

type mockBufferWriter struct {
	lock sync.Mutex
	buf  bytes.Buffer
}

func (m *mockBufferWriter) Write(p []byte) (n int, err error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.buf.Write(p)
}

func (m *mockBufferWriter) String() string {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.buf.String()
}

func TestReport(t *testing.T) {
	r := Report{
		Time:       time.Now().Format(time.RFC3339),
		Level:      Info,
		Msg:        "Hello World",
		CallerInfo: runtimeinfo.GetCallerInfo(1),
	}

	reportLog := fmt.Sprintf(`"time":"%s","level":"%s","msg":"%s","file":"%s","package":"%s","function":"%s","line":"%d"}
`, r.Time, r.Level, r.Msg, r.CallerInfo.File, r.CallerInfo.Package, r.CallerInfo.Function, r.CallerInfo.Line)

	t.Run("should timout when mutex is lock", func(t *testing.T) {
		l := NewLogger()
		_l, ok := l.(*logger)
		if !ok {
			t.Fatalf("NewLogger did not returned a instance of logger")
		}

		buf := &mockBufferWriter{}
		_l.writer.SetWriters(buf)

		expectedReport := "{" + reportLog

		_l.reportLock.Lock()

		go l.Report(r)
		time.Sleep(10 * time.Millisecond)

		if buf.String() != "" {
			t.Errorf("expected nothing on buffer, but got %q", buf.String())
		}

		_l.reportLock.Unlock()
		time.Sleep(10 * time.Millisecond)

		if buf.String() != expectedReport {
			t.Errorf("expected read on buffer %q, but got %q", expectedReport, buf.String())
		}
	})

	t.Run("should write the information of report", func(t *testing.T) {
		l := NewLogger()
		_l, ok := l.(*logger)
		if !ok {
			t.Fatalf("NewLogger did not returned a instance of logger")
		}

		buf := &mockBufferWriter{}

		_l.writer.SetWriters(buf)

		l.Report(r)

		expectedReport := "{" + reportLog

		if buf.String() != expectedReport {
			t.Errorf("expected read on buffer %q, but got %q", expectedReport, buf.String())
		}
	})

	t.Run("should write the key and value when staticFields is not empty", func(t *testing.T) {
		l := NewLogger()
		_l, ok := l.(*logger)
		if !ok {
			t.Fatalf("newlogger did not returned a instance of logger")
		}

		attrs := make(map[string]string, 1)
		attrs["hello"] = "world"
		helloReport := `{"hello":"world",`
		expectedReport := helloReport + reportLog

		_l.staticFields = attrs

		buf := &bytes.Buffer{}
		_l.writer.SetWriters(buf)

		l.Report(r)

		if buf.String() != expectedReport {
			t.Errorf("expected read on buffer %q, but got %q", expectedReport, buf.String())
		}
	})
}

func TestFlushReports(t *testing.T) {
	r := Report{
		Time:       time.Now().Format(time.RFC3339),
		Level:      Info,
		Msg:        "Hello World",
		CallerInfo: runtimeinfo.GetCallerInfo(1),
	}

	reportLog := fmt.Sprintf(`{"time":"%s","level":"%s","msg":"%s","file":"%s","package":"%s","function":"%s","line":"%d"}
`, r.Time, r.Level, r.Msg, r.CallerInfo.File, r.CallerInfo.Package, r.CallerInfo.Function, r.CallerInfo.Line)

	t.Run("should not flush any report when buffer reports is empty", func(t *testing.T) {
		l := NewLogger()
		_l, ok := l.(*logger)
		if !ok {
			t.Fatalf("NewLogger did not returned a instance of logger")
		}

		buf := &bytes.Buffer{}
		_l.writer.SetWriters(buf)

		l.FlushReports()

		if buf.String() != "" {
			t.Errorf("expected nothing on buffer, but got %q", buf.String())
		}
	})

	t.Run("should flush the report", func(t *testing.T) {
		l := NewLogger()
		_l, ok := l.(*logger)
		if !ok {
			t.Fatalf("NewLogger did not returned a instance of logger")
		}

		buf := &bytes.Buffer{}
		_l.writer.SetWriters(buf)

		_l.reports <- r

		l.FlushReports()

		if buf.String() != reportLog {
			t.Errorf("expected read on buffer %q, but got %q", reportLog, buf.String())
		}
	})
}

func TestHandleReports(t *testing.T) {
	r := Report{
		Time:       time.Now().Format(time.RFC3339),
		Level:      Info,
		Msg:        "Hello World",
		CallerInfo: runtimeinfo.GetCallerInfo(1),
	}

	reportLog := fmt.Sprintf(`{"time":"%s","level":"%s","msg":"%s","file":"%s","package":"%s","function":"%s","line":"%d"}
`, r.Time, r.Level, r.Msg, r.CallerInfo.File, r.CallerInfo.Package, r.CallerInfo.Function, r.CallerInfo.Line)

	t.Run("should handle the report and close the logger", func(t *testing.T) {
		l := NewLogger()
		_l, ok := l.(*logger)
		if !ok {
			t.Fatalf("NewLogger did not returned a instance of logger")
		}

		buf := &mockBufferWriter{}
		_l.writer.SetWriters(buf)

		ctx, cancel := context.WithCancel(context.Background())
		go l.HandleReports(ctx)

		_l.reports <- r
		time.Sleep(time.Millisecond)

		if buf.String() != reportLog {
			t.Errorf("expected read on buffer %q, but got %q", reportLog, buf.String())
		}

		cancel()
		time.Sleep(time.Millisecond)
		if !_l.getStatusCloseReport() {
			t.Error("expected report handle reports to be closed, but remain open")
		}

		buf.buf.Reset()
		_l.reports <- r
		time.Sleep(time.Millisecond)

		if buf.String() != "" {
			t.Errorf("expected nothing on buffer, but got %q", buf.String())
		}
	})
}

func TestWriter(t *testing.T) {
	t.Run("should return the writers set on logger", func(t *testing.T) {
		l := NewLogger()

		w := l.Writer()
		if w == nil {
			t.Error("expected the writes was instace")
		}
		if reflect.ValueOf(w).IsNil() {
			t.Error("expected the write was empty")
		}
	})
}

func TestMemory(t *testing.T) {
	t.Run("should return the record memory", func(t *testing.T) {
		l := NewLogger()

		m := l.Memory()
		if m == nil {
			t.Error("expected the memory was instace")
		}
		if reflect.ValueOf(m).IsNil() {
			t.Error("expected the memory was empty")
		}
	})
}

func TestSetStaticFields(t *testing.T) {
	t.Run("should set the static fields", func(t *testing.T) {
		l := NewLogger()
		_l, ok := l.(*logger)
		if !ok {
			t.Fatalf("newlogger did not returned a instance of logger")
		}

		attrs := make(map[string]string, 2)
		attrs["hello"] = "world"
		attrs["ionic"] = "health"

		l.SetStaticFields(attrs)

		if len(_l.staticFields) != len(attrs) {
			t.Errorf("expected the size of static field to be %q, but got %q", len(attrs), len(_l.staticFields))
		}

		for key, value := range attrs {
			if _l.staticFields[key] != value {
				t.Errorf("expected the value of static fields with key=%v to be %q, but got %q", key, value, _l.staticFields[key])
			}
		}

		for key, value := range _l.staticFields {
			if value != attrs[key] {
				t.Errorf("expected the l.staticFields[%q]=%q was set on attrs, but got %q", key, value, attrs[key])
			}
		}
	})
}

func TestSetReportQueueSize(t *testing.T) {
	t.Run("should set the size of record reports", func(t *testing.T) {
		l := NewLogger()
		_l, ok := l.(*logger)
		if !ok {
			t.Fatalf("newlogger did not returned a instance of logger")
		}

		size := 100
		l.SetReportQueueSize(uint(size))

		if cap(_l.reports) != size {
			t.Errorf("expected the size of report to be %v, but got %q", size, cap(_l.reports))
		}
	})
}

func TestSetTraceMode(t *testing.T) {
	t.Run("should set trace mode to true", func(t *testing.T) {
		l := NewLogger()
		_l, ok := l.(*logger)
		if !ok {
			t.Fatalf("newlogger did not returned a instance of logger")
		}

		mode := true
		l.SetTraceMode(mode)

		if _l.traceMode != mode {
			t.Errorf("expected the trace mode to be %v, but got %v", mode, _l.traceMode)
		}
	})

	t.Run("should set the trace mode to false", func(t *testing.T) {
		l := NewLogger()
		_l, ok := l.(*logger)
		if !ok {
			t.Fatalf("newlogger did not returned a instance of logger")
		}

		mode := false
		l.SetTraceMode(mode)

		if _l.traceMode != mode {
			t.Errorf("expected the trace mode to be %v, but got %v", mode, _l.traceMode)
		}
	})
}

func TestTraceMode(t *testing.T) {
	t.Run("should return the trace mode (true)", func(t *testing.T) {
		l := NewLogger()
		_l, ok := l.(*logger)
		if !ok {
			t.Fatalf("newlogger did not returned a instance of logger")
		}

		mode := true
		_l.traceMode = mode

		if l.TraceMode() != mode {
			t.Errorf("expected the trace mode to be %v, but got %v", mode, l.TraceMode())
		}
	})

	t.Run("should return the trace mode (false)", func(t *testing.T) {
		l := NewLogger()
		_l, ok := l.(*logger)
		if !ok {
			t.Fatalf("newlogger did not returned a instance of logger")
		}

		mode := false
		_l.traceMode = mode

		if l.TraceMode() != mode {
			t.Errorf("expected the trace mode to be %v, but got %v", mode, l.TraceMode())
		}
	})
}
