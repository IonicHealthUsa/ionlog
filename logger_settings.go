package ionlog

import (
	"fmt"
	"io"
	"os"

	"github.com/IonicHealthUsa/ionlog/internal/core/rotationengine"
	"github.com/IonicHealthUsa/ionlog/internal/service"
)

type customAttrs func(i service.ICoreService)

func SetTraceMode(mode bool) {
	logger.LogEngine().SetTraceMode(mode)
}

// SetAttributes sets the log SetAttributes
// fns is a variadic parameter that accepts customAttrs
func SetAttributes(fns ...customAttrs) {
	if logger.Status() == service.Running {
		fmt.Fprint(os.Stderr, "Logger is already running, cannot set attributes\n")
		return
	}

	for _, fn := range fns {
		fn(logger)
	}
}

// WithWriters sets the write targets for the logger, every log will be written to these targets.
func WithWriters(w ...io.Writer) customAttrs {
	return func(i service.ICoreService) {
		i.LogEngine().Writer().SetWriters(w...)
	}
}

// WithStaticFields sets the static fields for the logger, every log will have these fields.
// usage: WithStaicFields(map[string]string{"key": "value", "key2": "value2", ...})
func WithStaticFields(attrs map[string]string) customAttrs {
	return func(i service.ICoreService) {
		i.LogEngine().SetStaticFields(attrs)
	}
}

// WithLogFileRotation enables log file rotation, specifying the directory where log files will be stored, the maximum size of the log folder in bytes, and the rotation frequency.
func WithLogFileRotation(folder string, folderMaxSize uint, period rotationengine.PeriodicRotation) customAttrs {
	return func(i service.ICoreService) {
		i.CreateRotationService(folder, folderMaxSize, period)
	}
}

func SetQueueSize(size uint) customAttrs {
	return func(i service.ICoreService) {
		i.LogEngine().SetReportQueueSize(size)
	}
}
