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

// Info logs a message with level info. Arguments are handled in the manner of fmt.Printf.
func Info(msg string, args ...any) {
	ioncore.Logger().SendReport(ioncore.NewIonReport(slog.LevelInfo, fmt.Sprintf(msg, args...), ioncore.GetRecordInformation()))
}

// Error logs a message with level error. Arguments are handled in the manner of fmt.Printf.
func Error(msg string, args ...any) {
	ioncore.Logger().SendReport(ioncore.NewIonReport(slog.LevelError, fmt.Sprintf(msg, args...), ioncore.GetRecordInformation()))
}

// Warn logs a message with level warn. Arguments are handled in the manner of fmt.Printf.
func Warn(msg string, args ...any) {
	ioncore.Logger().SendReport(ioncore.NewIonReport(slog.LevelWarn, fmt.Sprintf(msg, args...), ioncore.GetRecordInformation()))
}

// Debug logs a message with level debug. Arguments are handled in the manner of fmt.Printf.
func Debug(msg string, args ...any) {
	ioncore.Logger().SendReport(ioncore.NewIonReport(slog.LevelDebug, fmt.Sprintf(msg, args...), ioncore.GetRecordInformation()))
}

// LogOnce logs a message with level info only once. Arguments are handled in the manner of fmt.Printf.
// Each function call will log the message only once.
// Avoid using it in a sintax like this: LogOnce("Logging..."); LogOnce("Logging...")
func LogOnce(msg string, args ...any) {
	callInfo := ioncore.GetRecordInformation()
	pkg := string(callInfo[0].(slog.Attr).Value.String())
	function := string(callInfo[1].(slog.Attr).Value.String())
	file := string(callInfo[2].(slog.Attr).Value.String())
	line := int(callInfo[3].(slog.Attr).Value.Int64())

	recordMsg := fmt.Sprintf(msg, args...)

	proceed := recordhistory.LogOnce(
		ioncore.Logger().History(),
		pkg,
		function,
		file,
		line,
		recordMsg,
	)

	if proceed {
		ioncore.Logger().SendReport(ioncore.NewIonReport(slog.LevelInfo, recordMsg, callInfo))
	}
}

// LogOnChange logs a message with level info only when the message changes. Arguments are handled in the manner of fmt.Printf.
// Each function call will log the message only when it changes.
// Avoid using it in a sintax like this: LogOnChange("Logging..."); LogOnChange("Logging...")
func LogOnChange(msg string, args ...any) {
	callInfo := ioncore.GetRecordInformation()
	pkg := string(callInfo[0].(slog.Attr).Value.String())
	function := string(callInfo[1].(slog.Attr).Value.String())
	file := string(callInfo[2].(slog.Attr).Value.String())
	line := int(callInfo[3].(slog.Attr).Value.Int64())

	recordMsg := fmt.Sprintf(msg, args...)

	proceed := recordhistory.LogOnChange(
		ioncore.Logger().History(),
		pkg,
		function,
		file,
		line,
		recordMsg,
	)

	if proceed {
		ioncore.Logger().SendReport(ioncore.NewIonReport(slog.LevelInfo, recordMsg, callInfo))
	}
}
