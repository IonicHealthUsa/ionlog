// Package ioncore provides the core functionalities of the logger.
// It is responsible for handling the logger service, the log writer, and the log engine.
package logcore

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/IonicHealthUsa/ionlog/internal/infrastructure/memory"
	"github.com/IonicHealthUsa/ionlog/internal/ionservice"
	"github.com/IonicHealthUsa/ionlog/internal/logrotation"
)

type service struct {
	ctx           context.Context
	cancel        context.CancelFunc
	serviceWg     sync.WaitGroup
	serviceStatus ionservice.ServiceStatus
}

type ionLogger struct {
	service

	logMux sync.Mutex

	logsMemory memory.IRecordMemory
	logRotate  logrotation.ILogRotation

	writerHandler ionWriter
	reports       chan *IonReport

	staticFields map[string]string
}

type IIonLogger interface {
	ionservice.IService

	LogsMemory() memory.IRecordMemory

	SetLogRotationSettings(folder string, maxFolderSize uint, rotation logrotation.PeriodicRotation)

	SetReportsBufferSize(size uint)

	SetTargets(targets ...io.Writer)

	SetStaticFields(attrs map[string]string)

	SendReport(r *IonReport)
	LogReport(r *IonReport)
}

const maxReports = 100

var logger *ionLogger

func init() {
	logger = newLogger()
}

func newLogger() *ionLogger {
	l := &ionLogger{}
	l.ctx, l.cancel = context.WithCancel(context.Background())
	l.reports = make(chan *IonReport, maxReports)

	l.logsMemory = memory.NewRecordMemory()

	return l
}

// Logger returns the logger instance
func Logger() IIonLogger {
	return logger
}

func ResetLogger() {
	logger = newLogger()
}

func (i *ionLogger) LogsMemory() memory.IRecordMemory {
	return i.logsMemory
}

// SetLogRotationSettings is a proxy to the log rotation service
func (i *ionLogger) SetLogRotationSettings(folder string, maxFolderSize uint, rotation logrotation.PeriodicRotation) {
	i.logRotate = logrotation.NewLogFileRotation()
	i.logRotate.SetLogRotationSettings(folder, maxFolderSize, rotation)
	i.writerHandler.AddTarget(i.logRotate)
}

func (i *ionLogger) SetReportsBufferSize(size uint) {
	i.reports = make(chan *IonReport, size)
}

func (i *ionLogger) SetTargets(targets ...io.Writer) {
	i.writerHandler.SetTargets(targets...)
}

func (i *ionLogger) SetStaticFields(attrs map[string]string) {
	i.staticFields = attrs
}

func (i *ionLogger) SendReport(r *IonReport) {
	select {
	case i.reports <- r:
		return

	case <-time.After(1000 * time.Millisecond):
		fmt.Fprintf(os.Stderr, "Report timed out queue length: %d\n", len(i.reports))
		return
	}
}

// LogReport a synchronous version of SendReport
func (i *ionLogger) LogReport(r *IonReport) {
	i.log(r)
}

func (i *ionLogger) syncReports() {
	for len(i.reports) > 0 {
		r := <-i.reports
		i.log(r)
	}
}

func (i *ionLogger) log(r *IonReport) {
	i.logMux.Lock()
	defer i.logMux.Unlock()

	msg := i.createLog(r)
	i.writerHandler.Write(msg)
}

func (i *ionLogger) createLog(r *IonReport) []byte {
	payload := map[string]string{
		"time":     r.Datetime.Format(time.RFC3339),
		"level":    r.Level.String(),
		"msg":      r.Msg,
		"file":     r.File,
		"package":  r.PackageName,
		"function": r.Function,
		"line":     strconv.Itoa(r.Line),
	}

	maps.Copy(payload, i.staticFields)

	msg, err := json.Marshal(payload)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal payload: %v\n", err)
	}
	return append(msg, '\n')
}
