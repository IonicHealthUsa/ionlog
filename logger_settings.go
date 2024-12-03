package ionlog

import (
	"context"
	"io"
	"log/slog"
)

type customAttrs func(i *ionLogger)

// SetLogAttributes sets the log SetLogAttributes
// fns is a variadic parameter that accepts customAttrs
func SetLogAttributes(fns ...customAttrs) {
	for _, fn := range fns {
		fn(logger)
	}
}

// WithTargets sets the write targets for the logger, every log will be written to these targets.
func WithTargets(ctx context.Context, w ...io.Writer) customAttrs {
	return func(i *ionLogger) {
		i.ctx = ctx
		i.writerHandler.writeTargets = w
	}
}

// WithStaicFields sets the static fields for the logger, every log will have these fields.
// usage: WithStaicFields(map[string]string{"key": "value", "key2": "value2", ...})
func WithStaicFields(attrs map[string]string) customAttrs {
	return func(i *ionLogger) {
		index := 0
		slogAttrs := make([]slog.Attr, len(attrs))
		for k, v := range attrs {
			slogAttrs[index] = slog.String(k, v)
			index++
		}
		handler := logger.createDefaultLogHandler().WithAttrs(slogAttrs)
		i.logHandler = slog.New(handler)
	}
}
