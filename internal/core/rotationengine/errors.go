package rotationengine

import "errors"

var (
	ErrLogFileNotSet             = errors.New("log file not set")
	ErrCouldNotCheckFolderStatus = errors.New("could not check folder status")
	ErrNoLogFileFound            = errors.New("no log file found")
)
