package benchmark

import (
	"fmt"
	"testing"

	"github.com/IonicHealthUsa/ionlog"
)

func BenchmarkLogOnceNoChange(b *testing.B) {
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
}

func BenchmarkLogOnceNoChangeFormated(b *testing.B) {
	ionlog.Start()
	defer ionlog.Stop()

	b.ResetTimer()

	b.Run("LogOnceDebug", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.LogOnceDebugf("my msg is: %v", fakeMessage)
		}
	})

	b.Run("LogOnceInfo", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.LogOnceInfof("my msg is: %v", fakeMessage)
		}
	})

	b.Run("LogOnceWarn", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.LogOnceWarnf("my msg is: %v", fakeMessage)
		}
	})

	b.Run("LogOnceError", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.LogOnceErrorf("my msg is: %v", fakeMessage)
		}
	})
}

func BenchmarkLogOnceWithChange(b *testing.B) {
	ionlog.Start()
	defer ionlog.Stop()

	b.ResetTimer()

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
}

func BenchmarkLogOnceWithChangeFormated(b *testing.B) {
	ionlog.Start()
	defer ionlog.Stop()

	b.ResetTimer()

	b.Run("LogOnceDebug", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.LogOnceDebugf("%d: %v", i, fakeMessage)
		}
	})

	b.Run("LogOnceInfo", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.LogOnceInfof("%d: %v", i, fakeMessage)
		}
	})

	b.Run("LogOnceWarn", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.LogOnceWarnf("%d: %v", i, fakeMessage)
		}
	})

	b.Run("LogOnceError", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ionlog.LogOnceErrorf("%d: %v", i, fakeMessage)
		}
	})
}
