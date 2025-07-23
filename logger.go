package ionlog

import (
	"fmt"
	"sync"
	"time"

	"github.com/IonicHealthUsa/ionlog/internal/core/logengine"
	"github.com/IonicHealthUsa/ionlog/internal/core/runtimeinfo"
	"github.com/IonicHealthUsa/ionlog/internal/service"
	"github.com/IonicHealthUsa/ionlog/internal/usecases"
)

// Start begin the ionlog reports when it does not running
func Start() {
	startSync := sync.WaitGroup{}
	startSync.Add(1)
	go logger.Start(&startSync)
	startSync.Wait()
}

// Stop stop the ionlog reports and reset the logger
func Stop() {
	logger.Stop()
	logger = service.NewCoreService() // Reset the logger
}

// Flush flushes the reports to the output writers.
func Flush() {
	logger.LogEngine().FlushReports()
}

// Info logs a message with level info.
func Info(msg string) {
	logger.LogEngine().AsyncReport(
		logengine.ReportType{
			Time:       time.Now().Format(time.RFC3339),
			Level:      logengine.Info,
			Msg:        msg,
			CallerInfo: runtimeinfo.GetCallerInfo(2),
		},
	)
}

// Infof logs a message with level info.
// Arguments are handled in the manner of fmt.Printf.
func Infof(msg string, args ...any) {
	logger.LogEngine().AsyncReport(
		logengine.ReportType{
			Time:       time.Now().Format(time.RFC3339),
			Level:      logengine.Info,
			Msg:        fmt.Sprintf(msg, args...),
			CallerInfo: runtimeinfo.GetCallerInfo(2),
		},
	)
}

// Error logs a message with level error.
func Error(msg string) {
	logger.LogEngine().AsyncReport(
		logengine.ReportType{
			Time:       time.Now().Format(time.RFC3339),
			Level:      logengine.Error,
			Msg:        msg,
			CallerInfo: runtimeinfo.GetCallerInfo(2),
		},
	)
}

// Errorf logs a message with level error.
// Arguments are handled in the manner of fmt.Printf.
func Errorf(msg string, args ...any) {
	logger.LogEngine().AsyncReport(
		logengine.ReportType{
			Time:       time.Now().Format(time.RFC3339),
			Level:      logengine.Error,
			Msg:        fmt.Sprintf(msg, args...),
			CallerInfo: runtimeinfo.GetCallerInfo(2),
		},
	)
}

// Warn logs a message with level warn.
func Warn(msg string) {
	logger.LogEngine().AsyncReport(
		logengine.ReportType{
			Time:       time.Now().Format(time.RFC3339),
			Level:      logengine.Warn,
			Msg:        msg,
			CallerInfo: runtimeinfo.GetCallerInfo(2),
		},
	)
}

// Warnf logs a message with level warn.
// Arguments are handled in the manner of fmt.Printf.
func Warnf(msg string, args ...any) {
	logger.LogEngine().AsyncReport(
		logengine.ReportType{
			Time:       time.Now().Format(time.RFC3339),
			Level:      logengine.Warn,
			Msg:        fmt.Sprintf(msg, args...),
			CallerInfo: runtimeinfo.GetCallerInfo(2),
		},
	)
}

// Debug logs a message with level debug.
func Debug(msg string) {
	logger.LogEngine().AsyncReport(
		logengine.ReportType{
			Time:       time.Now().Format(time.RFC3339),
			Level:      logengine.Debug,
			Msg:        msg,
			CallerInfo: runtimeinfo.GetCallerInfo(2),
		},
	)
}

// Debugf logs a message with level debug.
// Arguments are handled in the manner of fmt.Printf.
func Debugf(msg string, args ...any) {
	logger.LogEngine().AsyncReport(
		logengine.ReportType{
			Time:       time.Now().Format(time.RFC3339),
			Level:      logengine.Debug,
			Msg:        fmt.Sprintf(msg, args...),
			CallerInfo: runtimeinfo.GetCallerInfo(2),
		},
	)
}

// Trace logs a message with level trace only when trace mode is enable.
func Trace(msg string) {
	if !logger.LogEngine().TraceMode() {
		return
	}
	logger.LogEngine().Report(
		logengine.ReportType{
			Time:       time.Now().Format(time.RFC3339),
			Level:      logengine.Trace,
			Msg:        msg,
			CallerInfo: runtimeinfo.GetCallerInfo(2),
		},
	)
}

// Tracef logs a message with level trace only when trace mode is enable.
// Arguments are handled in the manner of fmt.Printf.
func Tracef(msg string, args ...any) {
	if !logger.LogEngine().TraceMode() {
		return
	}
	logger.LogEngine().Report(
		logengine.ReportType{
			Time:       time.Now().Format(time.RFC3339),
			Level:      logengine.Trace,
			Msg:        fmt.Sprintf(msg, args...),
			CallerInfo: runtimeinfo.GetCallerInfo(2),
		},
	)
}

// LogOnceInfo logs a message with level info only once time.
func LogOnceInfo(msg string) {
	logOnce(logengine.Info, msg)
}

// LogOnceInfof logs a message with level info only once time.
// Arguments are handled in the manner of fmt.Printf.
func LogOnceInfof(msg string, args ...any) {
	logOnce(logengine.Info, fmt.Sprintf(msg, args...))
}

// LogOnceError logs a message with level error only once time.
// Arguments are handled in the manner of fmt.Printf.
func LogOnceError(msg string) {
	logOnce(logengine.Error, msg)
}

// LogOnceErrorf logs a message with level error only once time.
// Arguments are handled in the manner of fmt.Printf.
func LogOnceErrorf(msg string, args ...any) {
	logOnce(logengine.Error, fmt.Sprintf(msg, args...))
}

// LogOnceWarn logs a message with level warn only once time.
// Arguments are handled in the manner of fmt.Printf.
func LogOnceWarn(msg string) {
	logOnce(logengine.Warn, msg)
}

// LogOnceWarnf logs a message with level warn only once time.
// Arguments are handled in the manner of fmt.Printf.
func LogOnceWarnf(msg string, args ...any) {
	logOnce(logengine.Warn, fmt.Sprintf(msg, args...))
}

// LogOnceDebug logs a message with level debug only once time.
// Arguments are handled in the manner of fmt.Printf.
func LogOnceDebug(msg string) {
	logOnce(logengine.Debug, msg)
}

// LogOnceDebugf logs a message with level debug only once time.
// Arguments are handled in the manner of fmt.Printf.
func LogOnceDebugf(msg string, args ...any) {
	logOnce(logengine.Debug, fmt.Sprintf(msg, args...))
}

// logOnce send the information about the function
// which called the log level to report queue asynchronously.
func logOnce(level logengine.Level, recordMsg string) {
	callerInfo := runtimeinfo.GetCallerInfo(3)

	proceed := usecases.LogOnce(
		logger.LogEngine().Memory(),
		recordMsg,
		callerInfo.File,
		callerInfo.Package,
		callerInfo.Function,
	)

	if !proceed {
		return
	}

	logger.LogEngine().AsyncReport(
		logengine.ReportType{
			Time:       time.Now().Format(time.RFC3339),
			Level:      level,
			Msg:        recordMsg,
			CallerInfo: callerInfo,
		},
	)
}
