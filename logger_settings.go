package ionlog

import (
	"io"
	"time"

	"github.com/IonicHealthUsa/ionlog/internal/core/rotationengine"
	"github.com/IonicHealthUsa/ionlog/internal/observability/loki"
	"github.com/IonicHealthUsa/ionlog/internal/service"
)

type customAttrs func(i service.ICoreService)

// SetAttributes sets the log SetAttributes
// fns is a variadic parameter that accepts customAttrs
func SetAttributes(fns ...customAttrs) {
	Flush()

	for _, fn := range fns {
		fn(logger)
	}
}

// WithWriters sets the write targets for the logger,
// every log will be written to these targets.
func WithWriters(w ...io.Writer) customAttrs {
	return func(i service.ICoreService) {
		i.LogEngine().Writer().AddWriter(w...)
	}
}

// WithoutWriters deletes the write targets for the logger.
func WithoutWriters(w ...io.Writer) customAttrs {
	return func(i service.ICoreService) {
		i.LogEngine().Writer().DeleteWriter(w...)
	}
}

// WithStaticFields sets the static fields for the logger, every log will have these fields.
// usage: WithStaicFields(map[string]string{"key": "value", "key2": "value2", ...})
func WithStaticFields(attrs map[string]string) customAttrs {
	return func(i service.ICoreService) {
		i.LogEngine().AddStaticFields(attrs)
	}
}

// WithoutStaticFields remove the static fields for the logger.
// Use the key of the static field to remove.
func WithoutStaticFields(fields ...string) customAttrs {
	return func(i service.ICoreService) {
		i.LogEngine().DeleteStaticField(fields...)
	}
}

// WithLogFileRotation enables log file rotation,
// specifying the directory where log files will be stored,
// the maximum size of the log folder in bytes, and the rotation frequency.
func WithLogFileRotation(
	folder string,
	folderMaxSize uint,
	period rotationengine.PeriodicRotation,
) customAttrs {
	return func(i service.ICoreService) {
		i.CreateRotationService(folder, folderMaxSize, period)
	}
}

// WithQueueSize sets the size of the reports queue,
// which stores logs before sending them to a file descriptor.
func WithQueueSize(size uint) customAttrs {
	return func(i service.ICoreService) {
		i.LogEngine().SetReportQueueSize(size)
	}
}

// WithTraceMode enables trace log mode.
// For default, the trace mode is disable,
// to enable is need pass a true boolean.
func WithTraceMode(mode bool) customAttrs {
	return func(i service.ICoreService) {
		i.LogEngine().SetTraceMode(mode)
	}
}

// WithLokiIntegration sets up Loki integration for the logger
func WithLokiIntegration(config LokiConfig) customAttrs {
	return func(i service.ICoreService) {
		// Convert public config to internal config
		internalConfig := config.toInternalConfig()

		integration, err := loki.NewLokiIntegration(internalConfig)
		if err != nil {
			// Log error but don't fail the setup
			// The logger will continue to work without Loki
			return
		}

		// Store the integration globally for graceful shutdown
		lokiMutex.Lock()
		lokiIntegration = integration
		lokiMutex.Unlock()

		// Add Loki writer to the logger
		i.LogEngine().Writer().AddWriter(integration.Writer())
	}
}

// GetLokiIntegration returns the current Loki integration (if any)
func GetLokiIntegration() *loki.LokiIntegration {
	return getLokiIntegration()
}

// WithLokiShutdownTimeout sets the timeout for graceful Loki shutdown
// This overrides both the default and environment variable settings
func WithLokiShutdownTimeout(timeout time.Duration) customAttrs {
	return func(i service.ICoreService) {
		setLokiShutdownTimeout(timeout)
	}
}
