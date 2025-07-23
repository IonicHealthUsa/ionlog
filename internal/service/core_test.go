package service

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/IonicHealthUsa/ionlog/internal/core/logengine"
	"github.com/IonicHealthUsa/ionlog/internal/core/rotationengine"
	"github.com/IonicHealthUsa/ionlog/internal/core/runtimeinfo"
)

func TestNewCoreService(t *testing.T) {
	t.Run("should return the core service", func(t *testing.T) {
		cs := NewCoreService()
		if cs == nil {
			t.Error("expected core service no to be nil")
		}
		if reflect.ValueOf(cs).IsNil() {
			t.Error("expected core service no to be nil")
		}

		_cs, ok := cs.(*coreService)
		if !ok {
			t.Error("expected a instance of core service to implement ICoreService")
		}

		if _cs.serviceStatus != 0 {
			t.Errorf("expected service status to be '0', but got %q", _cs.serviceStatus)
		}

		if reflect.ValueOf(_cs.logEngine).IsNil() {
			t.Error("expected log engine no to be nil")
		}

		if _cs.rotationService != nil {
			t.Error("expected rotation service to be nil, but got a instance")
		}

		if reflect.ValueOf(cs.LogEngine()).IsNil() {
			t.Error("expected the log engine interface no to be nil")
		}
	})
}

func TestLogEngine(t *testing.T) {
	t.Run("should return a interface of logger", func(t *testing.T) {
		cs := NewCoreService()

		logger := cs.LogEngine()

		if logger == nil {
			t.Errorf("expected a interface of logger")
		}

		if reflect.ValueOf(logger).IsNil() {
			t.Errorf("expected a interface of logger")
		}
	})

	t.Run("should return a nil interface of logger when logger is not instanced", func(t *testing.T) {
		cs := &coreService{}

		logger := cs.LogEngine()

		if logger != nil {
			t.Errorf("expected a nil interface of logger")
		}
	})
}

func TestCreateRotationService(t *testing.T) {
	folderName := "hello_world"
	t.Run("should add a writer", func(t *testing.T) {
		cs := NewCoreService()
		_cs, ok := cs.(*coreService)
		if !ok {
			t.Fatal("expected a instance of core service to implement ICoreService")
		}

		cs.CreateRotationService(folderName, 10, rotationengine.Daily)

		if _cs.rotationService == nil {
			t.Error("expected a interface of rotation service")
		}

		if reflect.ValueOf(_cs.rotationService).IsNil() {
			t.Error("expected a interface of rotation service")
		}

		_r, ok := _cs.rotationService.(*rotationService)
		if !ok {
			t.Fatal("expected a instance of rotation service")
		}

		if _r.rotationEngine == nil {
			t.Error("expected a interface of ratation engine")
		}

		if reflect.ValueOf(_r.rotationEngine).IsNil() {
			t.Error("expected a interface of rotation engine")
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a directory, but did not exist")
		}

		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})
}

type mockBufferWriter struct {
	lock sync.Mutex
	cond *sync.Cond
	buf  bytes.Buffer
}

func newMockBufferWriter() *mockBufferWriter {
	m := &mockBufferWriter{}
	m.cond = sync.NewCond(&m.lock)

	return m
}

func (m *mockBufferWriter) Write(p []byte) (n int, err error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.cond.Signal()
	return m.buf.Write(p)
}

func (m *mockBufferWriter) String() string {
	m.lock.Lock()
	defer m.lock.Unlock()

	for m.buf.Len() == 0 {
		m.cond.Wait()
	}

	return m.buf.String()
}

func (m *mockBufferWriter) Len() int {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.buf.Len()
}

func TestStart_Core(t *testing.T) {
	r := logengine.ReportType{
		Time:       time.Now().Format(time.RFC3339),
		Level:      logengine.Info,
		Msg:        "Hello World",
		CallerInfo: runtimeinfo.GetCallerInfo(1),
	}

	reportLog := fmt.Sprintf(`{"time":"%s","level":"%s","msg":"%s","file":"%s","package":"%s","function":"%s","line":"%d"}
`, r.Time, r.Level, r.Msg, r.CallerInfo.File, r.CallerInfo.Package, r.CallerInfo.Function, r.CallerInfo.Line)

	t.Run("should receive the message on buffer", func(t *testing.T) {
		cs := NewCoreService()
		_cs, ok := cs.(*coreService)
		if !ok {
			t.Fatal("NewCoreService() did not return *coreService")
		}

		buf := newMockBufferWriter()
		cs.LogEngine().Writer().AddWriter(buf)

		startSync := sync.WaitGroup{}
		startSync.Add(1)
		go cs.Start(&startSync)

		cs.LogEngine().AsyncReport(r)
		startSync.Wait()

		if buf.String() != reportLog {
			t.Errorf("expected the report log to be %q, but got %q", reportLog, buf.String())
		}

		_cs.cancel()
	})

	t.Run("should ionlog started", func(t *testing.T) {
		cs := NewCoreService()
		_cs, ok := cs.(*coreService)
		if !ok {
			t.Fatal("NewCoreService() did not return *coreService")
		}

		startSync := sync.WaitGroup{}
		startSync.Add(1)
		go cs.Start(&startSync)
		startSync.Wait()

		_cs.serviceStatusLock.Lock()
		if _cs.serviceStatus != Running {
			t.Errorf("expected the status of logger service to be %q, but got %q", Running, _cs.serviceStatus)
		}
		_cs.serviceStatusLock.Unlock()

		_cs.cancel()
	})

	t.Run("should ionlog started with rotation service", func(t *testing.T) {
		folderName := "rotaion_start"
		maxFolderSize := rotationengine.GB
		rotation := rotationengine.Daily

		cs := NewCoreService()
		_cs, ok := cs.(*coreService)
		if !ok {
			t.Fatal("NewCoreService() did not return *coreService")
		}

		cs.CreateRotationService(folderName, maxFolderSize, rotation)
		_rs, ok := _cs.rotationService.(*rotationService)
		if !ok {
			t.Fatal("NewRotationService() did not return *rotationService")
		}

		startSync := sync.WaitGroup{}
		startSync.Add(1)
		go cs.Start(&startSync)
		startSync.Wait()

		_cs.serviceStatusLock.Lock()
		if _cs.serviceStatus != Running {
			t.Errorf("expected the status of logger service to be %q, but got %q", Running, _cs.serviceStatus)
		}
		_cs.serviceStatusLock.Unlock()

		_rs.serviceStatusLock.Lock()
		if _rs.serviceStatus != Running {
			t.Errorf("expected the status of rotation service to be %q, but got %q", Running, _rs.serviceStatus)
		}
		_rs.serviceStatusLock.Unlock()

		_cs.cancel()

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a direcory %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})

	t.Run("should send the report after ionlog start", func(t *testing.T) {
		cs := NewCoreService()
		_cs, ok := cs.(*coreService)
		if !ok {
			t.Fatal("NewCoreService() did not return *coreService")
		}

		buf := newMockBufferWriter()
		cs.LogEngine().Writer().AddWriter(buf)

		cs.LogEngine().AsyncReport(r)

		if buf.Len() != 0 {
			t.Errorf("expected the buffer length to be 0, but got %d", buf.Len())
		}

		startSync := sync.WaitGroup{}
		startSync.Add(1)
		go cs.Start(&startSync)
		startSync.Wait()

		if buf.String() != reportLog {
			t.Errorf("expected the report log to be %q, but got %q", reportLog, buf.String())
		}

		_cs.cancel()
	})
}

func TestStop_Core(t *testing.T) {
	t.Run("should core service stop", func(t *testing.T) {
		cs := NewCoreService()
		_cs, ok := cs.(*coreService)
		if !ok {
			t.Fatal("NreCoreService() did not return *coreService")
		}

		startSync := sync.WaitGroup{}
		startSync.Add(1)
		go cs.Start(&startSync)
		startSync.Wait()

		_cs.serviceStatusLock.Lock()
		if _cs.serviceStatus != Running {
			t.Errorf("expected the status of logger service to be %q, but got %q", Running, _cs.serviceStatus)
		}
		_cs.serviceStatusLock.Unlock()

		cs.Stop()
		time.Sleep(time.Millisecond)

		_cs.serviceStatusLock.Lock()
		if _cs.serviceStatus != Stopped {
			t.Errorf("expected the status of logger service to be %q, but got %q", Stopped, _cs.serviceStatus)
		}
		_cs.serviceStatusLock.Unlock()
	})

	t.Run("should core service stop with rotaion service", func(t *testing.T) {
		folderName := "rotaion_stop"
		maxFolderSize := rotationengine.GB
		rotation := rotationengine.Daily

		cs := NewCoreService()
		_cs, ok := cs.(*coreService)
		if !ok {
			t.Fatal("NreCoreService() did not return *coreService")
		}

		cs.CreateRotationService(folderName, maxFolderSize, rotation)
		_rs, ok := _cs.rotationService.(*rotationService)
		if !ok {
			t.Fatal("NewRotationService() did not return *rotationService")
		}

		startSync := sync.WaitGroup{}
		startSync.Add(1)
		go cs.Start(&startSync)
		startSync.Wait()

		_cs.serviceStatusLock.Lock()
		if _cs.serviceStatus != Running {
			t.Errorf("expected the status of logger service to be %q, but got %q", Running, _cs.serviceStatus)
		}
		_cs.serviceStatusLock.Unlock()

		_rs.serviceStatusLock.Lock()
		if _rs.serviceStatus != Running {
			t.Errorf("expected the status of rotation service to be %q, but got %q", Running, _rs.serviceStatus)
		}
		_rs.serviceStatusLock.Unlock()

		cs.Stop()
		time.Sleep(time.Millisecond)

		if _cs.serviceStatus != Stopped {
			t.Errorf("expected the status of logger service to be %q, but got %q", Stopped, _cs.serviceStatus)
		}
		if _rs.serviceStatus != Stopped {
			t.Errorf("expected the status of rotation service to be %q, but got %q", Stopped, _rs.serviceStatus)
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a direcory %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})
}

func TestStatusCore(t *testing.T) {
	t.Run("should return running status", func(t *testing.T) {
		cs := NewCoreService()
		_cs, ok := cs.(*coreService)
		if !ok {
			t.Error("expected a instance of core service to implement ICoreService")
		}

		_cs.serviceStatus = Running

		if cs.Status() != Running {
			t.Errorf("expected status to be %q, but got %q", Running, cs.Status())
		}
	})

	t.Run("should return stopped status", func(t *testing.T) {
		cs := NewCoreService()
		_cs, ok := cs.(*coreService)
		if !ok {
			t.Error("expected a instance of core service to implement ICoreService")
		}

		_cs.serviceStatus = Stopped

		if cs.Status() != Stopped {
			t.Errorf("expected status to be %q, but got %q", Stopped, cs.Status())
		}
	})
}

func TestSetServiceStatus(t *testing.T) {
	t.Run("should set runnig status on service status", func(t *testing.T) {
		cs := NewCoreService()
		_cs, ok := cs.(*coreService)
		if !ok {
			t.Error("expected a instance of core service to implement ICoreService")
		}

		_cs.setServiceStatus(Running)

		if _cs.serviceStatus != Running {
			t.Errorf("expected service status to be %q, but got %q", Running, _cs.serviceStatus)
		}
	})

	t.Run("should set stopped status on service status", func(t *testing.T) {
		cs := NewCoreService()
		_cs, ok := cs.(*coreService)
		if !ok {
			t.Error("expected a instance of core service to implement ICoreService")
		}

		_cs.setServiceStatus(Stopped)

		if _cs.serviceStatus != Stopped {
			t.Errorf("expected service status to be %q, but got %q", Stopped, _cs.serviceStatus)
		}
	})
}
