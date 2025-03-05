package benchmark

import (
	"log/slog"
	"sync"
	"testing"

	"github.com/IonicHealthUsa/ionlog"
)

func BenchmarkBasicIonlog(b *testing.B) {
	ionlog.SetLogAttributes()

	// Start the logger service
	ionlog.Start()
	defer ionlog.Stop()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ionlog.Info("log test")
	}
}

func BenchmarkIonlogStress(b *testing.B) {
	ionlog.SetLogAttributes(
		ionlog.WithTargets(ionlog.DefaultOutput),
		ionlog.WithLogFileRotation(ionlog.DefaultLogFolder, 10*ionlog.Gibibyte, ionlog.Daily),
	)

	ionlog.Start()
	defer ionlog.Stop()

	b.ResetTimer()

	wg := sync.WaitGroup{}
	b.Logf("b.N: %d", b.N)
	for range b.N {
		wg.Add(100)
		for range 100 {
			go func() {
				for range 1000 {
					ionlog.Info("log test")
				}
				wg.Done()
			}()
		}
	}
	wg.Wait()
}

func BenchmarkDefaultSlog(b *testing.B) {
	for i := 0; i < b.N; i++ {
		slog.Info("log test")
	}
}
