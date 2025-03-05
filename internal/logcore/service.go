package logcore

import (
	"log/slog"
	"sync"

	"github.com/IonicHealthUsa/ionlog/internal/ionservice"
)

// Start starts the logger service, it blocks until the service is stopped
func (i *ionLogger) Start(startSync *sync.WaitGroup) {
	slog.Info("Logger service starting...")

	i.serviceWg.Add(1)
	defer i.serviceWg.Done()

	i.setServiceStatus(ionservice.Running)
	defer i.setServiceStatus(ionservice.Stopped)

	if i.logRotate != nil {
		rotateSync := sync.WaitGroup{}
		rotateSync.Add(1)
		go i.logRotate.Start(&rotateSync)
		rotateSync.Wait()
	}

	if startSync != nil {
		startSync.Done()
	}

	for {
		select {
		case <-i.ctx.Done():
			slog.Debug("Logger service stopped by context")

			i.syncReports()
			return

		case r := <-i.reports:
			i.log(r.level, r.msg, r.args...)
		}
	}
}

// Stop stops the logger by canceling the context and waiting for the worker to finish
func (i *ionLogger) Stop() {
	slog.Info("Logger service stopping...")

	i.cancel()
	i.serviceWg.Wait()
	slog.Info("All reports have been processed and the logger service has stopped.")

	if i.logRotate != nil {
		i.logRotate.Stop()
		slog.Info("Logger rotation service has stopped.")
	}
}

// Status returns the status of the logger service
func (i *ionLogger) Status() ionservice.ServiceStatus {
	return i.serviceStatus
}

func (i *ionLogger) setServiceStatus(status ionservice.ServiceStatus) {
	i.serviceStatus = status
}
