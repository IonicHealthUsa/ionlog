package ionlog

import (
	"fmt"
	"io"
	"os"

	"github.com/IonicHealthUsa/ionlog/internal/ionservice"
	"github.com/IonicHealthUsa/ionlog/internal/logcore"
	"github.com/IonicHealthUsa/ionlog/internal/logrotation"
)

type customAttrs func(i logcore.IIonLogger)

// SetLogAttributes sets the log SetLogAttributes
// fns is a variadic parameter that accepts customAttrs
func SetLogAttributes(fns ...customAttrs) {
	if logcore.Logger().Status() == ionservice.Running {
		fmt.Fprint(os.Stderr, "Logger is already running, cannot set attributes\n")
		return
	}

	for _, fn := range fns {
		fn(logcore.Logger())
	}
}

// WithTargets sets the write targets for the logger, every log will be written to these targets.
func WithTargets(w ...io.Writer) customAttrs {
	return func(i logcore.IIonLogger) {
		i.SetTargets(w...)
	}
}

// WithStaticFields sets the static fields for the logger, every log will have these fields.
// usage: WithStaicFields(map[string]string{"key": "value", "key2": "value2", ...})
func WithStaticFields(attrs map[string]string) customAttrs {
	return func(i logcore.IIonLogger) {
		i.SetStaticFields(attrs)
	}
}

// WithLogFileRotation enables log file rotation, specifying the directory where log files will be stored, the maximum size of the log folder in bytes, and the rotation frequency.
func WithLogFileRotation(folder string, folderMaxSize uint, period logrotation.PeriodicRotation) customAttrs {
	return func(i logcore.IIonLogger) {
		i.SetLogRotationSettings(folder, folderMaxSize, period)
	}
}

func SetReportsBufferSize(size uint) customAttrs {
	return func(i logcore.IIonLogger) {
		i.SetReportsBufferSize(size)
	}
}
