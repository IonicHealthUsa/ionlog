package ionlog

import (
	"fmt"
	"log/slog"

	ioncore "github.com/IonicHealthUsa/ionlog/internal/core"
	recordhistory "github.com/IonicHealthUsa/ionlog/internal/record_history"
)

func Start() {
	ioncore.Logger().Start()
}

func Stop() {
	ioncore.Logger().Stop()
}

// Info logs a message with level info.
func Info(msg string) {
	ioncore.Logger().SendReport(ioncore.NewIonReport(slog.LevelInfo, msg, ioncore.GetRecordInformation()))
}

// Error logs a message with level error.
func Error(msg string) {
	ioncore.Logger().SendReport(ioncore.NewIonReport(slog.LevelError, msg, ioncore.GetRecordInformation()))
}

// Warn logs a message with level warn.
func Warn(msg string) {
	ioncore.Logger().SendReport(ioncore.NewIonReport(slog.LevelWarn, msg, ioncore.GetRecordInformation()))
}

// Debug logs a message with level debug.
func Debug(msg string) {
	ioncore.Logger().SendReport(ioncore.NewIonReport(slog.LevelDebug, msg, ioncore.GetRecordInformation()))
}

// LogOnceInfo logs a message with level info only once.
func LogOnceInfo(msg string) {
	logOnce(slog.LevelInfo, msg)
}

// LogOnceError logs a message with level info only once.
func LogOnceError(msg string) {
	logOnce(slog.LevelError, msg)
}

// LogOnceWarn logs a message with level warn only once.
func LogOnceWarn(msg string) {
	logOnce(slog.LevelWarn, msg)
}

// LogOnceDebug logs a message with level debug only once.
func LogOnceDebug(msg string) {
	logOnce(slog.LevelDebug, msg)
}

// LogOnChangeInfo logs a message with level info only when the message changes.
func LogOnChangeInfo(msg string) {
	logOnChange(slog.LevelInfo, msg)
}

// LogOnChangeError logs a message with level error only when the message changes.
func LogOnChangeError(msg string) {
	logOnChange(slog.LevelError, msg)
}

// LogOnChangeWarn logs a message with level warn only when the message changes.
func LogOnChangeWarn(msg string) {
	logOnChange(slog.LevelWarn, msg)
}

// LogOnChangeDebug logs a message with level debug only when the message changes.
func LogOnChangeDebug(msg string) {
	logOnChange(slog.LevelDebug, msg)
}

// Infof logs a message with level info. Arguments are handled in the manner of fmt.Printf.
func Infof(msg string, args ...any) {
	ioncore.Logger().SendReport(ioncore.NewIonReport(slog.LevelInfo, fmt.Sprintf(msg, args...), ioncore.GetRecordInformation()))
}

// Errorf logs a message with level error. Arguments are handled in the manner of fmt.Printf.
func Errorf(msg string, args ...any) {
	ioncore.Logger().SendReport(ioncore.NewIonReport(slog.LevelError, fmt.Sprintf(msg, args...), ioncore.GetRecordInformation()))
}

// Warnf logs a message with level warn. Arguments are handled in the manner of fmt.Printf.
func Warnf(msg string, args ...any) {
	ioncore.Logger().SendReport(ioncore.NewIonReport(slog.LevelWarn, fmt.Sprintf(msg, args...), ioncore.GetRecordInformation()))
}

// Debugf logs a message with level debug. Arguments are handled in the manner of fmt.Printf.
func Debugf(msg string, args ...any) {
	ioncore.Logger().SendReport(ioncore.NewIonReport(slog.LevelDebug, fmt.Sprintf(msg, args...), ioncore.GetRecordInformation()))
}

// LogOnceInfof logs a message with level info only once.
// Arguments are handled in the manner of fmt.Printf.
func LogOnceInfof(msg string, args ...any) {
	logOnce(slog.LevelInfo, fmt.Sprintf(msg, args...))
}

// LogOnceErrorf logs a message with level error only once.
// Arguments are handled in the manner of fmt.Printf.
func LogOnceErrorf(msg string, args ...any) {
	logOnce(slog.LevelError, fmt.Sprintf(msg, args...))
}

// LogOnceWarnf logs a message with level warn only once.
// Arguments are handled in the manner of fmt.Printf.
func LogOnceWarnf(msg string, args ...any) {
	logOnce(slog.LevelWarn, fmt.Sprintf(msg, args...))
}

// LogOnceDebugf logs a message with level debug only once.
// Arguments are handled in the manner of fmt.Printf.
func LogOnceDebugf(msg string, args ...any) {
	logOnce(slog.LevelDebug, fmt.Sprintf(msg, args...))
}

// LogOnChangeInfof logs a message with level info only when the message changes.
// Arguments are handled in the manner of fmt.Printf.
func LogOnChangeInfof(msg string, args ...any) {
	logOnChange(slog.LevelInfo, fmt.Sprintf(msg, args...))
}

// LogOnChangeErrorf logs a message with level error only when the message changes.
// Arguments are handled in the manner of fmt.Printf.
func LogOnChangeErrorf(msg string, args ...any) {
	logOnChange(slog.LevelError, fmt.Sprintf(msg, args...))
}

// LogOnChangeWarnf logs a message with level warn only when the message changes.
// Arguments are handled in the manner of fmt.Printf.
func LogOnChangeWarnf(msg string, args ...any) {
	logOnChange(slog.LevelWarn, fmt.Sprintf(msg, args...))
}

// LogOnChangeDebugf logs a message with level debug only when the message changes.
// Arguments are handled in the manner of fmt.Printf.
func LogOnChangeDebugf(msg string, args ...any) {
	logOnChange(slog.LevelDebug, fmt.Sprintf(msg, args...))
}

// logOnce logs a message with level info only once. Arguments are handled in the manner of fmt.Printf.
// Each function call will log the message only once.
// Avoid using it in a sintax like this: LogOnce("Logging..."); LogOnce("Logging...")
func logOnce(level slog.Level, recordMsg string) {
	callInfo := ioncore.GetRecordInformation()
	pkg := string(callInfo[0].(slog.Attr).Value.String())
	function := string(callInfo[1].(slog.Attr).Value.String())
	file := string(callInfo[2].(slog.Attr).Value.String())
	line := int(callInfo[3].(slog.Attr).Value.Int64())

	proceed := recordhistory.LogOnce(
		ioncore.Logger().History(),
		pkg,
		function,
		file,
		line,
		recordMsg,
	)

	if proceed {
		ioncore.Logger().SendReport(ioncore.NewIonReport(level, recordMsg, callInfo))
	}
}

// logOnChange logs a message with level info only when the message changes. Arguments are handled in the manner of fmt.Printf.
// Each function call will log the message only when it changes.
// Avoid using it in a sintax like this: LogOnChange("Logging..."); LogOnChange("Logging...")
func logOnChange(level slog.Level, recordMsg string) {
	callInfo := ioncore.GetRecordInformation()
	pkg := string(callInfo[0].(slog.Attr).Value.String())
	function := string(callInfo[1].(slog.Attr).Value.String())
	file := string(callInfo[2].(slog.Attr).Value.String())
	line := int(callInfo[3].(slog.Attr).Value.Int64())

	proceed := recordhistory.LogOnChange(
		ioncore.Logger().History(),
		pkg,
		function,
		file,
		line,
		recordMsg,
	)

	if proceed {
		ioncore.Logger().SendReport(ioncore.NewIonReport(level, recordMsg, callInfo))
	}
}
