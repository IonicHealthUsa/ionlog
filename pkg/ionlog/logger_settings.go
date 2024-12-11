package ionlog

import (
	"io"
	"log/slog"

	ioncore "github.com/IonicHealthUsa/ionlog/internal/core"
	ionlogfile "github.com/IonicHealthUsa/ionlog/internal/logfile"
	ionservice "github.com/IonicHealthUsa/ionlog/internal/service"
)

type customAttrs func(i ioncore.IIonLogger)

const (
	Daily   = ionlogfile.Daily
	Weekly  = ionlogfile.Weekly
	Monthly = ionlogfile.Monthly
)

// SetLogAttributes sets the log SetLogAttributes
// fns is a variadic parameter that accepts customAttrs
func SetLogAttributes(fns ...customAttrs) {
	if ioncore.Logger().Status() == ionservice.Running {
		slog.Warn("Logger is already running, cannot set attributes")
		return
	}

	for _, fn := range fns {
		fn(ioncore.Logger())
	}
}

// WithTargets sets the write targets for the logger, every log will be written to these targets.
func WithTargets(w ...io.Writer) customAttrs {
	return func(i ioncore.IIonLogger) {
		i.SetTargets(w...)
	}
}

// WithStaticFields sets the static fields for the logger, every log will have these fields.
// usage: WithStaicFields(map[string]string{"key": "value", "key2": "value2", ...})
func WithStaticFields(attrs map[string]string) customAttrs {
	return func(i ioncore.IIonLogger) {
		index := 0
		slogAttrs := make([]slog.Attr, len(attrs))
		for k, v := range attrs {
			slogAttrs[index] = slog.String(k, v)
			index++
		}
		handler := i.CreateDefaultLogHandler().WithAttrs(slogAttrs)
		i.SetLogEngine(slog.New(handler))
	}
}

// WithLogFileRotation sets the log file rotation period and the folder where the log files will be stored.
func WithLogFileRotation(folder string, period ionlogfile.PeriodicRotation) customAttrs {
	return func(i ioncore.IIonLogger) {
		i.SetRotationPeriod(period)
		i.SetFolder(folder)
	}
}
