package benchmark

import (
	"testing"

	"github.com/IonicHealthUsa/ionlog"
)

var fakeMessage = "We shall not cease from exploration and the end of all our exploring will be to arrive where we started and know the place for the first time."

func BenchmarkBasicLogs(b *testing.B) {
	// Start the logger service
	ionlog.Start()
	defer ionlog.Stop()

	b.ResetTimer()

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
}

func BenchmarkBasicLogsFormated(b *testing.B) {
	// Start the logger service
	ionlog.Start()
	defer ionlog.Stop()

	b.ResetTimer()

	b.Run("Trace", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.Tracef("my msg is: %v", fakeMessage)
		}
	})

	b.Run("Debug", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.Debugf("my msg is: %v", fakeMessage)
		}
	})

	b.Run("Info", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.Infof("my msg is: %v", fakeMessage)
		}
	})

	b.Run("Error", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.Errorf("my msg is: %v", fakeMessage)
		}
	})

	b.Run("Warn", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.Warnf("my msg is: %v", fakeMessage)
		}
	})
}

func BenchmarkBasicLogsParallel(b *testing.B) {
	// Start the logger service
	ionlog.Start()
	defer ionlog.Stop()

	b.ResetTimer()

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
}

func BenchmarkBasicLogsFormatedParallel(b *testing.B) {
	// Start the logger service
	ionlog.Start()
	defer ionlog.Stop()

	b.ResetTimer()

	b.Run("Trace", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ionlog.Tracef("my msg is: %v", fakeMessage)
			}
		})
	})

	b.Run("Debug", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ionlog.Debugf("my msg is: %v", fakeMessage)
			}
		})
	})

	b.Run("Info", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ionlog.Infof("my msg is: %v", fakeMessage)
			}
		})
	})

	b.Run("Error", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ionlog.Errorf("my msg is: %v", fakeMessage)
			}
		})
	})

	b.Run("Warn", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ionlog.Warnf("my msg is: %v", fakeMessage)
			}
		})
	})
}
