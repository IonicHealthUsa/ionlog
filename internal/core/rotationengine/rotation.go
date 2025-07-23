package rotationengine

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/IonicHealthUsa/ionlog/internal/infrastructure/filesystem"
)

type rotationEngine struct {
	filesystem.Filesystem

	logFile       io.WriteCloser
	folder        string
	maxFolderSize uint
	rotation      PeriodicRotation
}

type IRotationEngine interface {
	io.Writer
	AutoChecks()
	CloseLogFile()
}

func NewRotationEngine(folder string, maxFolderSize uint, rotation PeriodicRotation) IRotationEngine {
	r := &rotationEngine{}

	r.Filesystem = filesystem.NewFileSystem(
		os.Stat,
		os.Mkdir,
		os.ReadDir,
		os.IsNotExist,
		os.OpenFile,
		os.Remove,
	)

	r.folder = folder
	r.maxFolderSize = maxFolderSize
	r.rotation = rotation
	r.AutoChecks()

	return r
}

// Write writes the log message to the log file.
func (r *rotationEngine) Write(p []byte) (n int, err error) {
	if r.logFile == nil {
		return 0, ErrLogFileNotSet
	}

	return r.logFile.Write(p)
}

func (r *rotationEngine) AutoChecks() {
	r.autoRotate()
	r.autoCheckFolderSize()
}

func (r *rotationEngine) CloseLogFile() {
	r.closeFile()
}

// closeFile closes the log file.
func (r *rotationEngine) closeFile() {
	if r.logFile != nil {
		if err := r.logFile.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Error to close current log file: %v\n", err)
		}
		r.logFile = nil
	}
}

func (r *rotationEngine) setLogFile(file io.WriteCloser) {
	if file == nil {
		fmt.Fprint(os.Stderr, "Cannot set the log file: file is not valid\n")
		return
	}

	r.closeFile()
	r.logFile = file
}

func (r *rotationEngine) autoRotate() {
	if err := r.assertFolder(); err != nil {
		fmt.Fprintf(os.Stderr, "Error in assert folder: %v", err)
		return
	}

	fileName, err := r.getMostRecentLogFile()

	if err == ErrNoLogFileFound {
		r.createNewFile()
		return
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	fileDate, err := r.getFileDate(fileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	if r.checkRotation(fileDate) {
		r.createNewFile()
		return
	}

	// no rotaion needed, check if file is open
	if r.logFile == nil {
		actualFile, err := r.OpenFile(filepath.Join(r.folder, fileName), os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		r.setLogFile(actualFile)
	}
}

func (r *rotationEngine) autoCheckFolderSize() {
	if r.maxFolderSize == NoMaxFolderSize {
		return
	}

	size, err := r.getFolderSize()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	if size <= r.maxFolderSize {
		return
	}

	oldestFile, err := r.getOldestLogFile()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	if err = r.RemoveFile(filepath.Join(r.folder, oldestFile)); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	// check if it need to create a new file
	files, err := r.getAllfiles()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	if len(files) == 0 {
		r.createNewFile()
		return
	}
}
