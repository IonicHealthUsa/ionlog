package logcore

import (
	"log/slog"

	ionservice "github.com/IonicHealthUsa/ionlog/internal/interfaces"
	"github.com/IonicHealthUsa/ionlog/internal/logrotation"
)

// Start starts the logger service, it blocks until the service is stopped
func (i *ionLogger) Start() error {
	i.serviceStatus = ionservice.Running
	defer func() { i.serviceStatus = ionservice.Stopped }()

	// user has chosen to auto rotate the log file
	if i.rotationPeriod != logrotation.NoAutoRotate {
		i.logRotateService = logrotation.NewLogFileRotation(i.folder, i.maxFolderSize, i.rotationPeriod)

		// logRotateService is manages a file, so it is a target...
		i.SetTargets(append(i.Targets(), i.logRotateService)...)

		// block until the log rotate service sets up the file to write to.
		i.logRotateService.BlockWrite()
		i.servicesWg.Add(1)
		go func() {
			defer i.servicesWg.Done()

			if err := i.logRotateService.Start(); err != nil {
				slog.Error(err.Error())
				return
			}
		}()
	}

	i.servicesWg.Add(1)
	go func() {
		defer i.servicesWg.Done()

		i.handleIonReports()
	}()

	return nil
}

// Status returns the status of the logger service
func (i *ionLogger) Status() ionservice.ServiceStatus {
	return i.serviceStatus
}

// Stop stops the logger by canceling the context and waiting for the worker to finish
func (i *ionLogger) Stop() {
	slog.Debug("Logger service stopping...")

	i.cancel()

	i.reportsWg.Wait()
	slog.Debug("All reports have been processed")

	if i.logRotateService != nil {
		i.logRotateService.Stop()
	}

	i.servicesWg.Wait()
}
