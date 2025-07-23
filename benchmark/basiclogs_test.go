package benchmark

import (
	"testing"

	"github.com/IonicHealthUsa/ionlog"
)

var fakeMessage = "We shall not cease from exploration and the end of all our exploring will be to arrive where we started and know the place for the first time."

func BenchmarkBasicLogs(b *testing.B) {
	ionlog.SetAttributes(
		ionlog.WithQueueSize(1000),
		ionlog.WithTraceMode(true),
	)

	ionlog.Start()
	defer ionlog.Stop()

	b.Run("Trace", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.Trace(fakeMessage)
		}
	})

	b.Run("Debug", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.Debug(fakeMessage)
		}
	})

	b.Run("Info", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.Info(fakeMessage)
		}
	})

	b.Run("Error", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.Error(fakeMessage)
		}
	})

	b.Run("Warn", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.Warn(fakeMessage)
		}
	})

	b.Run("Tracef", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.Tracef("log: %v", fakeMessage)
		}
	})

	b.Run("Debugf", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.Debugf("log: %v", fakeMessage)
		}
	})

	b.Run("Infof", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.Infof("log: %v", fakeMessage)
		}
	})

	b.Run("Errorf", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.Errorf("log: %v", fakeMessage)
		}
	})

	b.Run("Warnf", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.Warnf("log: %v", fakeMessage)
		}
	})
}

func BenchmarkBasicLogsParallel(b *testing.B) {
	ionlog.SetAttributes(
		ionlog.WithQueueSize(1000),
		ionlog.WithTraceMode(true),
	)

	ionlog.Start()
	defer ionlog.Stop()

	b.Run("Trace", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ionlog.Trace(fakeMessage)
			}
		})
	})

	b.Run("Debug", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ionlog.Debug(fakeMessage)
			}
		})
	})

	b.Run("Info", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ionlog.Info(fakeMessage)
			}
		})
	})

	b.Run("Error", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ionlog.Error(fakeMessage)
			}
		})
	})

	b.Run("Warn", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ionlog.Warn(fakeMessage)
			}
		})
	})

	b.Run("Tracef", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ionlog.Tracef("log: %v", fakeMessage)
			}
		})
	})

	b.Run("Debugf", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ionlog.Debugf("log: %v", fakeMessage)
			}
		})
	})

	b.Run("Infof", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ionlog.Infof("log: %v", fakeMessage)
			}
		})
	})

	b.Run("Errorf", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ionlog.Errorf("log: %v", fakeMessage)
			}
		})
	})

	b.Run("Warnf", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ionlog.Warnf("log: %v", fakeMessage)
			}
		})
	})
}
