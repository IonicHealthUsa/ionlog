package benchmark

import (
	"log/slog"
	"testing"

	"github.com/IonicHealthUsa/ionlog"
)

var fakeMessage = "We shall not cease from exploration and the end of all our exploring will be to arrive where we started and know the place for the first time."

func BenchmarkInfo(b *testing.B) {
	slog.SetDefault(slog.New(slog.NewTextHandler(ionlog.DefaultOutput, &slog.HandlerOptions{Level: slog.LevelError})))

	ionlog.SetLogAttributes(
		ionlog.SetReportsBufferSizer(1000),
	)

	// Start the logger service
	ionlog.Start()
	defer ionlog.Stop()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ionlog.Info(fakeMessage)
	}
}

func BenchmarkInfoLogFile(b *testing.B) {
	slog.SetDefault(slog.New(slog.NewTextHandler(ionlog.DefaultOutput, &slog.HandlerOptions{Level: slog.LevelError})))

	ionlog.SetLogAttributes(
		ionlog.SetReportsBufferSizer(1000),
		ionlog.WithLogFileRotation(ionlog.DefaultLogFolder, 10*ionlog.Gibibyte, ionlog.Daily),
	)

	// Start the logger service
	ionlog.Start()
	defer ionlog.Stop()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ionlog.Info(fakeMessage)
	}
}

func BenchmarkInfoPararel(b *testing.B) {
	slog.SetDefault(slog.New(slog.NewTextHandler(ionlog.DefaultOutput, &slog.HandlerOptions{Level: slog.LevelError})))

	ionlog.SetLogAttributes(
		ionlog.SetReportsBufferSizer(1000),
		ionlog.WithLogFileRotation(ionlog.DefaultLogFolder, 10*ionlog.Gibibyte, ionlog.Daily),
	)

	// Start the logger service
	ionlog.Start()
	defer ionlog.Stop()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ionlog.Info(fakeMessage)
		}
	})
}
