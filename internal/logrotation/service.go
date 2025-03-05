package logrotation

import (
	"log/slog"
	"sync"
	"time"

	"github.com/IonicHealthUsa/ionlog/internal/ionservice"
)

func (l *logRotation) Start(startSync *sync.WaitGroup) {
	slog.Info("Logger rotation service starting...")

	l.serviceWg.Add(1)
	defer l.serviceWg.Done()

	l.setServiceStatus(ionservice.Running)
	defer l.setServiceStatus(ionservice.Stopped)

	if startSync != nil {
		startSync.Done()
	}

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-l.ctx.Done():
			slog.Info("Logger rotation service stopped by context.")
			return

		case <-ticker.C:
			l.autoRotate()
			l.autoCheckFolderSize()
		}
	}
}

func (l *logRotation) Stop() {
	slog.Debug("Logger rotation service stopping...")

	l.cancel()
	l.serviceWg.Wait()
}

func (l *logRotation) Status() ionservice.ServiceStatus {
	return l.serviceStatus
}

func (l *logRotation) setServiceStatus(status ionservice.ServiceStatus) {
	l.serviceStatus = status
}
