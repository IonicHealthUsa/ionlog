package ionlog

import (
	"context"
	"io"
	"log/slog"
	"os"
	"testing"
)

func TestSetLogAttributes(t *testing.T) {
	t.Run("WithTargets", func(t *testing.T) {
		logger = newLogger()
		logger.logHandler = slog.New(logger.createDefaultLogHandler())

		w := []io.Writer{Stdout(), os.Stderr}
		ctx := context.Background()

		SetLogAttributes(WithTargets(ctx, w...))

		for _, writer := range logger.writerHandler.writeTargets {
			if writer == nil {
				t.Errorf("Expected writer to be not nil")
			}
		}
	})
}
