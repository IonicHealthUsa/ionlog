package logrotation

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/IonicHealthUsa/ionlog/internal/infrastructure/filesystem"
	"github.com/IonicHealthUsa/ionlog/internal/ionservice"
)

type service struct {
	ctx             context.Context
	cancel          context.CancelFunc
	serviceWg       sync.WaitGroup
	incomingReports bool
	serviceStatus   ionservice.ServiceStatus
}

type logRotation struct {
	filesystem.Filesystem
	service

	logFile       io.WriteCloser
	folder        string
	maxFolderSize uint
	rotation      PeriodicRotation
}

type ILogRotation interface {
	io.Writer
	ionservice.IService

	SetLogRotationSettings(folder string, maxFolderSize uint, rotation PeriodicRotation)
}

func NewLogFileRotation() ILogRotation {
	r := &logRotation{}

	r.Filesystem = filesystem.NewFileSystem(
		os.Stat,
		os.Mkdir,
		os.ReadDir,
		os.IsNotExist,
		os.OpenFile,
		os.Remove,
	)

	r.ctx, r.cancel = context.WithCancel(context.Background())

	r.folder = ""
	r.maxFolderSize = NoMaxFolderSize
	r.rotation = NoAutoRotate

	return r
}

// Write writes the log message to the log file.
func (l *logRotation) Write(p []byte) (n int, err error) {
	if l.logFile == nil {
		return 0, ErrLogFileNotSet
	}

	return l.logFile.Write(p)
}

func (l *logRotation) SetLogRotationSettings(folder string, maxFolderSize uint, rotation PeriodicRotation) {
	l.folder = folder
	l.maxFolderSize = maxFolderSize
	l.rotation = rotation
	l.autoRotate()
	l.autoCheckFolderSize()
}

// closeFile closes the log file.
func (l *logRotation) closeFile() {
	if l.logFile != nil {
		if err := l.logFile.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error to close current log file: %v\n", err)
		}
		l.logFile = nil
	}
}

func (l *logRotation) setLogFile(file io.WriteCloser) {
	if file == nil {
		fmt.Fprint(os.Stderr, "Cannot set the log file: file is not valid\n")
		return
	}

	l.closeFile()
	l.logFile = file
}

func (l *logRotation) autoRotate() {
	if err := l.assertFolder(); err != nil {
		fmt.Fprintf(os.Stderr, "Error in assert folder: %v", err)
		return
	}

	fileName, err := l.getMostRecentLogFile()

	if err == ErrNoLogFileFound {
		l.createNewFile()
		return
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	fileDate, err := l.getFileDate(fileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	if l.checkRotation(fileDate) {
		l.createNewFile()
		return
	}

	// no rotaion needed, check if file is open
	if l.logFile == nil {
		actualFile, err := l.OpenFile(filepath.Join(l.folder, fileName), os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		l.setLogFile(actualFile)
	}
}

func (l *logRotation) autoCheckFolderSize() {
	if l.maxFolderSize == NoMaxFolderSize {
		return
	}

	size, err := l.getFolderSize()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	if size <= l.maxFolderSize {
		return
	}

	oldestFile, err := l.getOldestLogFile()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	if err = l.RemoveFile(filepath.Join(l.folder, oldestFile)); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	// check if it need to create a new file
	files, err := l.getAllfiles()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	if len(files) == 0 {
		l.createNewFile()
		return
	}
}
