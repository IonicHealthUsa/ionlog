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

func Start() {
	startSync := sync.WaitGroup{}
	startSync.Add(1)
	go logger.Start(&startSync)
	startSync.Wait()
}

func Stop() {
	logger.Stop()
	logger = service.NewCoreService() // Reset the logger
}

// Info logs a message with level info.
func Info(msg string) {
	logger.LogEngine().AsyncReport(
		logengine.Report{
			Time:       time.Now().Format(time.RFC3339),
			Level:      logengine.Info,
			Msg:        msg,
			CallerInfo: runtimeinfo.GetCallerInfo(2),
		},
	)
}

// Infof logs a message with level info. Arguments are handled in the manner of fmt.Printf.
func Infof(msg string, args ...any) {
	logger.LogEngine().AsyncReport(
		logengine.Report{
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
		logengine.Report{
			Time:       time.Now().Format(time.RFC3339),
			Level:      logengine.Error,
			Msg:        msg,
			CallerInfo: runtimeinfo.GetCallerInfo(2),
		},
	)
}

// Errorf logs a message with level error. Arguments are handled in the manner of fmt.Printf.
func Errorf(msg string, args ...any) {
	logger.LogEngine().AsyncReport(
		logengine.Report{
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
		logengine.Report{
			Time:       time.Now().Format(time.RFC3339),
			Level:      logengine.Warn,
			Msg:        msg,
			CallerInfo: runtimeinfo.GetCallerInfo(2),
		},
	)
}

// Warnf logs a message with level warn. Arguments are handled in the manner of fmt.Printf.
func Warnf(msg string, args ...any) {
	logger.LogEngine().AsyncReport(
		logengine.Report{
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
		logengine.Report{
			Time:       time.Now().Format(time.RFC3339),
			Level:      logengine.Debug,
			Msg:        msg,
			CallerInfo: runtimeinfo.GetCallerInfo(2),
		},
	)
}

// Debugf logs a message with level debug. Arguments are handled in the manner of fmt.Printf.
func Debugf(msg string, args ...any) {
	logger.LogEngine().AsyncReport(
		logengine.Report{
			Time:       time.Now().Format(time.RFC3339),
			Level:      logengine.Debug,
			Msg:        fmt.Sprintf(msg, args...),
			CallerInfo: runtimeinfo.GetCallerInfo(2),
		},
	)
}

func Trace(msg string) {
	logger.LogEngine().Report(
		logengine.Report{
			Time:       time.Now().Format(time.RFC3339),
			Level:      logengine.Trace,
			Msg:        msg,
			CallerInfo: runtimeinfo.GetCallerInfo(2),
		},
	)
}

func Tracef(msg string, args ...any) {
	logger.LogEngine().Report(
		logengine.Report{
			Time:       time.Now().Format(time.RFC3339),
			Level:      logengine.Trace,
			Msg:        fmt.Sprintf(msg, args...),
			CallerInfo: runtimeinfo.GetCallerInfo(2),
		},
	)
}

func LogOnceInfo(msg string) {
	logOnce(logengine.Info, msg)
}

func LogOnceInfof(msg string, args ...any) {
	logOnce(logengine.Info, fmt.Sprintf(msg, args...))
}

func LogOnceError(msg string) {
	logOnce(logengine.Error, msg)
}

func LogOnceErrorf(msg string, args ...any) {
	logOnce(logengine.Error, fmt.Sprintf(msg, args...))
}

func LogOnceWarn(msg string) {
	logOnce(logengine.Warn, msg)
}

func LogOnceWarnf(msg string, args ...any) {
	logOnce(logengine.Warn, fmt.Sprintf(msg, args...))
}

func LogOnceDebug(msg string) {
	logOnce(logengine.Debug, msg)
}

func LogOnceDebugf(msg string, args ...any) {
	logOnce(logengine.Debug, fmt.Sprintf(msg, args...))
}

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
		logengine.Report{
			Time:       time.Now().Format(time.RFC3339),
			Level:      level,
			Msg:        recordMsg,
			CallerInfo: callerInfo,
		},
	)
}
