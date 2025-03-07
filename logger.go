package ionlog

import (
	"fmt"
	"sync"
	"time"

	"github.com/IonicHealthUsa/ionlog/internal/logcore"
	"github.com/IonicHealthUsa/ionlog/internal/usecases"
)

func Start() {
	startSync := sync.WaitGroup{}
	startSync.Add(1)
	go logcore.Logger().Start(&startSync)
	startSync.Wait()
}

func Stop() {
	logcore.Logger().Stop()
	logcore.ResetLogger()
}

// Info logs a message with level info.
func Info(msg string) {
	callerInfo := logcore.GetCallerInfo(2)

	report := logcore.IonReport{}
	report.Datetime = time.Now()
	report.Level = logcore.Info
	report.Msg = msg
	report.File = callerInfo.File
	report.PackageName = callerInfo.PackageName
	report.Function = callerInfo.Function
	report.Line = callerInfo.Line

	logcore.Logger().SendReport(&report)
}

// Infof logs a message with level info. Arguments are handled in the manner of fmt.Printf.
func Infof(msg string, args ...any) {
	callerInfo := logcore.GetCallerInfo(2)

	report := logcore.IonReport{}
	report.Datetime = time.Now()
	report.Level = logcore.Info
	report.Msg = fmt.Sprintf(msg, args...)
	report.File = callerInfo.File
	report.PackageName = callerInfo.PackageName
	report.Function = callerInfo.Function
	report.Line = callerInfo.Line

	logcore.Logger().SendReport(&report)
}

// Error logs a message with level error.
func Error(msg string) {
	callerInfo := logcore.GetCallerInfo(2)

	report := logcore.IonReport{}
	report.Datetime = time.Now()
	report.Level = logcore.Error
	report.Msg = msg
	report.File = callerInfo.File
	report.PackageName = callerInfo.PackageName
	report.Function = callerInfo.Function
	report.Line = callerInfo.Line

	logcore.Logger().SendReport(&report)
}

// Errorf logs a message with level error. Arguments are handled in the manner of fmt.Printf.
func Errorf(msg string, args ...any) {
	callerInfo := logcore.GetCallerInfo(2)

	report := logcore.IonReport{}
	report.Datetime = time.Now()
	report.Level = logcore.Error
	report.Msg = fmt.Sprintf(msg, args...)
	report.File = callerInfo.File
	report.PackageName = callerInfo.PackageName
	report.Function = callerInfo.Function
	report.Line = callerInfo.Line

	logcore.Logger().SendReport(&report)
}

// Warn logs a message with level warn.
func Warn(msg string) {
	callerInfo := logcore.GetCallerInfo(2)

	report := logcore.IonReport{}
	report.Datetime = time.Now()
	report.Level = logcore.Warn
	report.Msg = msg
	report.File = callerInfo.File
	report.PackageName = callerInfo.PackageName
	report.Function = callerInfo.Function
	report.Line = callerInfo.Line

	logcore.Logger().SendReport(&report)
}

// Warnf logs a message with level warn. Arguments are handled in the manner of fmt.Printf.
func Warnf(msg string, args ...any) {
	callerInfo := logcore.GetCallerInfo(2)

	report := logcore.IonReport{}
	report.Datetime = time.Now()
	report.Level = logcore.Warn
	report.Msg = fmt.Sprintf(msg, args...)
	report.File = callerInfo.File
	report.PackageName = callerInfo.PackageName
	report.Function = callerInfo.Function
	report.Line = callerInfo.Line

	logcore.Logger().SendReport(&report)
}

// Debug logs a message with level debug.
func Debug(msg string) {
	callerInfo := logcore.GetCallerInfo(2)

	report := logcore.IonReport{}
	report.Datetime = time.Now()
	report.Level = logcore.Debug
	report.Msg = msg
	report.File = callerInfo.File
	report.PackageName = callerInfo.PackageName
	report.Function = callerInfo.Function
	report.Line = callerInfo.Line

	logcore.Logger().SendReport(&report)
}

// Debugf logs a message with level debug. Arguments are handled in the manner of fmt.Printf.
func Debugf(msg string, args ...any) {
	callerInfo := logcore.GetCallerInfo(2)

	report := logcore.IonReport{}
	report.Datetime = time.Now()
	report.Level = logcore.Debug
	report.Msg = fmt.Sprintf(msg, args...)
	report.File = callerInfo.File
	report.PackageName = callerInfo.PackageName
	report.Function = callerInfo.Function
	report.Line = callerInfo.Line

	logcore.Logger().SendReport(&report)
}

func LogOnceInfo(msg string) {
	logOnce(logcore.Info, msg)
}

func LogOnceInfof(msg string, args ...any) {
	logOnce(logcore.Info, fmt.Sprintf(msg, args...))
}

func LogOnceError(msg string) {
	logOnce(logcore.Error, msg)
}

func LogOnceErrorf(msg string, args ...any) {
	logOnce(logcore.Error, fmt.Sprintf(msg, args...))
}

func LogOnceWarn(msg string) {
	logOnce(logcore.Warn, msg)
}

func LogOnceWarnf(msg string, args ...any) {
	logOnce(logcore.Warn, fmt.Sprintf(msg, args...))
}

func LogOnceDebug(msg string) {
	logOnce(logcore.Debug, msg)
}

func LogOnceDebugf(msg string, args ...any) {
	logOnce(logcore.Debug, fmt.Sprintf(msg, args...))
}

func logOnce(level logcore.Level, recordMsg string) {
	callerInfo := logcore.GetCallerInfo(3)

	proceed := usecases.LogOnce(
		logcore.Logger().LogsMemory(),
		recordMsg,
		callerInfo.File,
		callerInfo.PackageName,
		callerInfo.Function,
	)

	if !proceed {
		return
	}

	report := logcore.IonReport{}
	report.Datetime = time.Now()
	report.Level = level
	report.Msg = recordMsg
	report.File = callerInfo.File
	report.PackageName = callerInfo.PackageName
	report.Function = callerInfo.Function
	report.Line = callerInfo.Line

	logcore.Logger().SendReport(&report)
}
