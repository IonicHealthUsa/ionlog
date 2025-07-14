package benchmark

import (
	"fmt"
	"testing"

	"github.com/IonicHealthUsa/ionlog"
)

func BenchmarkLogOnceNoChange(b *testing.B) {
	ionlog.SetAttributes(ionlog.WithQueueSize(1000))
	ionlog.Start()
	defer ionlog.Stop()

	b.ResetTimer()

	b.Run("LogOnceDebug", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.LogOnceDebug(fakeMessage)
		}
	})

	b.Run("LogOnceInfo", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.LogOnceInfo(fakeMessage)
		}
	})

	b.Run("LogOnceWarn", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.LogOnceWarn(fakeMessage)
		}
	})

	b.Run("LogOnceError", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.LogOnceError(fakeMessage)
		}
	})

	b.Run("LogOnceDebugf", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.LogOnceDebugf("my msg is: %v", fakeMessage)
		}
	})

	b.Run("LogOnceInfof", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.LogOnceInfof("my msg is: %v", fakeMessage)
		}
	})

	b.Run("LogOnceWarnf", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.LogOnceWarnf("my msg is: %v", fakeMessage)
		}
	})

	b.Run("LogOnceErrorf", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.LogOnceErrorf("my msg is: %v", fakeMessage)
		}
	})
}

func BenchmarkLogOnceWithChange(b *testing.B) {
	ionlog.SetAttributes(ionlog.WithQueueSize(1000))
	ionlog.Start()
	defer ionlog.Stop()

	b.Run("LogOnceDebug", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.LogOnceDebug(fmt.Sprintf("%d: %v", i, fakeMessage))
		}
	})

	b.Run("LogOnceInfo", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.LogOnceInfo(fmt.Sprintf("%d: %v", i, fakeMessage))
		}
	})

	b.Run("LogOnceWarn", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.LogOnceWarn(fmt.Sprintf("%d: %v", i, fakeMessage))
		}
	})

	b.Run("LogOnceError", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.LogOnceError(fmt.Sprintf("%d: %v", i, fakeMessage))
		}
	})

	b.Run("LogOnceDebugf", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.LogOnceDebugf("%d: %v", i, fakeMessage)
		}
	})

	b.Run("LogOnceInfof", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.LogOnceInfof("%d: %v", i, fakeMessage)
		}
	})

	b.Run("LogOnceWarnf", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.LogOnceWarnf("%d: %v", i, fakeMessage)
		}
	})

	b.Run("LogOnceErrorf", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.LogOnceErrorf("%d: %v", i, fakeMessage)
		}
	})
}

func BenchmarkLogOnceNoChangeParallel(b *testing.B) {
	ionlog.SetAttributes(ionlog.WithQueueSize(1000))
	ionlog.Start()
	defer ionlog.Stop()

	b.Run("LogOnceDebug", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ionlog.LogOnceDebug(fakeMessage)
			}
		})
	})

	b.Run("LogOnceInfo", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ionlog.LogOnceInfo(fakeMessage)
			}
		})
	})

	b.Run("LogOnceWarn", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ionlog.LogOnceWarn(fakeMessage)
			}
		})
	})

	b.Run("LogOnceError", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ionlog.LogOnceError(fakeMessage)
			}
		})
	})

	b.Run("LogOnceDebugf", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ionlog.LogOnceDebugf("my msg is: %v", fakeMessage)
			}
		})
	})

	b.Run("LogOnceInfof", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ionlog.LogOnceInfof("my msg is: %v", fakeMessage)
			}
		})
	})

	b.Run("LogOnceWarnf", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ionlog.LogOnceWarnf("my msg is: %v", fakeMessage)
			}
		})
	})

	b.Run("LogOnceErrorf", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				ionlog.LogOnceErrorf("my msg is: %v", fakeMessage)
			}
		})
	})
}

func BenchmarkLogOnceWithChangeParallel(b *testing.B) {
	ionlog.SetAttributes(ionlog.WithQueueSize(1000))
	ionlog.Start()
	defer ionlog.Stop()

	b.Run("LogOnceDebug", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				ionlog.LogOnceDebug(fmt.Sprintf("%d: %v", i, fakeMessage))
				i++
			}
		})
	})

	b.Run("LogOnceInfo", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				ionlog.LogOnceInfo(fmt.Sprintf("%d: %v", i, fakeMessage))
				i++
			}
		})
	})

	b.Run("LogOnceWarn", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				ionlog.LogOnceWarn(fmt.Sprintf("%d: %v", i, fakeMessage))
				i++
			}
		})
	})

	b.Run("LogOnceError", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				ionlog.LogOnceError(fmt.Sprintf("%d: %v", i, fakeMessage))
				i++
			}
		})
	})

	b.Run("LogOnceDebugf", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				ionlog.LogOnceDebugf("%d: %v", i, fakeMessage)
				i++
			}
		})
	})

	b.Run("LogOnceInfof", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				ionlog.LogOnceInfof("%d: %v", i, fakeMessage)
				i++
			}
		})
	})

	b.Run("LogOnceWarnf", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				ionlog.LogOnceWarnf("%d: %v", i, fakeMessage)
				i++
			}
		})
	})

	b.Run("LogOnceErrorf", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				ionlog.LogOnceErrorf("%d: %v", i, fakeMessage)
				i++
			}
		})
	})

	b.Run("LogOnceDebugf", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				ionlog.LogOnceDebugf("%d: %v", i, fakeMessage)
				i++
			}
		})
	})
}
