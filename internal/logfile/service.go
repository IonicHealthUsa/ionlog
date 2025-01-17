package ionlogfile

import (
	"fmt"
	"log/slog"
	"time"

	ionservice "github.com/IonicHealthUsa/ionlog/internal/service"
)

// Start starts the log file rotation service. It blocks until the service is stopped.
func (l *logFileRotation) Start() error {
	l.serviceStatus = ionservice.Running
	defer func() { l.serviceStatus = ionservice.Stopped }()

	if err := validateRotation(l.rotation); err != nil {
		return err
	}

	f, err := l.getActualFile()
	l.logFile = f
	defer l.closeFile()

	l.UnblockWrite()

	if err != nil {
		return err
	}

	var ticker *time.Ticker

	// every ticker check if the log file needs to be rotated
	switch l.rotation {
	case Daily:
		ticker = time.NewTicker(8 * time.Hour) // every 8 hours
	case Weekly:
		ticker = time.NewTicker(3 * 24 * time.Hour) // every 3 days
	case Monthly:
		ticker = time.NewTicker(7 * 24 * time.Hour) // every 7 days
	default:
		slog.Error(fmt.Sprintf("rotation was validated but it's invalid: %v", l.rotation))
		return ErrInvalidRotation
	}

	defer ticker.Stop()

	for {
		select {
		case <-l.ctx.Done():
			slog.Debug("logfile system stopped by context")
			return nil

		case <-ticker.C:
			err := func() error {
				l.writeMutex.Lock()
				defer l.writeMutex.Unlock()

				if err := l.logFile.Close(); err != nil {
					slog.Warn(err.Error())
				}

				f, err := l.getActualFile()
				l.logFile = f

				if err != nil {
					return err
				}

				return nil
			}()

			if err != nil {
				return err
			}
		}
	}
}

// Stop stops the log file rotation service.
func (l *logFileRotation) Stop() {
	l.writeMutex.Lock()
	defer l.writeMutex.Unlock()

	l.cancel()
}

// Status returns the status of the log file rotation service.
func (l *logFileRotation) Status() ionservice.ServiceStatus {
	return l.serviceStatus
}
