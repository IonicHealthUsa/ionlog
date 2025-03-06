package tests

import (
	"fmt"
	"log/slog"
	"sync"
	"testing"

	"github.com/IonicHealthUsa/ionlog"
)

var fakeMessage = "We shall not cease from exploration and the end of all our exploring will be to arrive where we started and know the place for the first time."

func TestStressWithLogFile(t *testing.T) {
	// supress internal logging in ionlog
	slog.SetDefault(slog.New(slog.NewTextHandler(ionlog.DefaultOutput, &slog.HandlerOptions{Level: slog.LevelError})))

	ionlog.SetLogAttributes(
		ionlog.WithLogFileRotation(ionlog.DefaultLogFolder, 10*ionlog.Gibibyte, ionlog.Daily),
	)

	ionlog.Start()
	defer ionlog.Stop()

	tests := 1000
	threads := 100
	logsPerThread := 1000

	wg := sync.WaitGroup{}
	for range tests {
		wg.Add(threads)
		for range threads {
			go func() {
				for range logsPerThread {
					ionlog.Info(fakeMessage)
				}
				wg.Done()
			}()
		}
	}
	wg.Wait()
}

func TestStressInfo(t *testing.T) {
	// supress internal logging in ionlog
	slog.SetDefault(slog.New(slog.NewTextHandler(ionlog.DefaultOutput, &slog.HandlerOptions{Level: slog.LevelError})))

	ionlog.Start()
	defer ionlog.Stop()

	tests := 1000
	threads := 100
	logsPerThread := 1000

	wg := sync.WaitGroup{}
	for range tests {
		wg.Add(threads)
		for range threads {
			go func() {
				for range logsPerThread {
					ionlog.Info(fakeMessage)
				}
				wg.Done()
			}()
		}
	}
	wg.Wait()
}

func TestStressLogOnceInfo(t *testing.T) {
	// supress internal logging in ionlog
	slog.SetDefault(slog.New(slog.NewTextHandler(ionlog.DefaultOutput, &slog.HandlerOptions{Level: slog.LevelError})))

	ionlog.Start()
	defer ionlog.Stop()

	tests := 1000
	threads := 100
	logsPerThread := 1000

	wg := sync.WaitGroup{}
	for range tests {
		wg.Add(threads)
		for range threads {
			go func() {
				for i := range logsPerThread {
					ionlog.LogOnceInfo(fmt.Sprintf("%s-%d", fakeMessage, i))
				}
				wg.Done()
			}()
		}
	}
	wg.Wait()
}
