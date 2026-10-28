package rotationengine

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/IonicHealthUsa/ionlog/internal/infrastructure/filesystem"
)

type rotationEngine struct {
	filesystem.Filesystem

	logFileMu     sync.RWMutex
	logFile       *os.File
	folder        string
	maxFolderSize uint
	rotation      PeriodicRotation
	maxFileSize   uint
	maxFiles      uint
	now           func() time.Time
	ready         bool
	files         []logFile // oldest first; the open file is last
	total         uint
	checkedAt     time.Time
}

type IRotationEngine interface {
	io.Writer
	CloseLogFile()
}

// Option sets the size at which a log file continues in a new part.
type Option func(*rotationEngine)

// WithMaxFileSize breaks a log file at size bytes, capped at the max folder size.
func WithMaxFileSize(size uint) Option {
	return func(r *rotationEngine) {
		if size > 0 {
			r.maxFileSize, r.maxFiles = size, 0
		}
	}
}

// WithMaxFiles breaks a log file at max folder size / count bytes.
func WithMaxFiles(count uint) Option {
	return func(r *rotationEngine) {
		if count > 0 {
			r.maxFiles, r.maxFileSize = count, 0
		}
	}
}

func NewRotationEngine(folder string, maxFolderSize uint, rotation PeriodicRotation, opts ...Option) IRotationEngine {
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
	r.maxFiles = DefaultMaxFiles
	r.now = time.Now
	for _, opt := range opts {
		opt(r)
	}
	if maxFolderSize != NoMaxFolderSize {
		if r.maxFiles > 0 {
			r.maxFileSize = max(maxFolderSize/r.maxFiles, 1)
		}
		r.maxFileSize = min(r.maxFileSize, maxFolderSize)
	}
	if rotation < NoAutoRotate || rotation > Monthly {
		fmt.Fprintf(os.Stderr, "[ionlog internal log] rotation value is invalid\n")
	}
	if err := r.load(); err != nil {
		fmt.Fprintf(os.Stderr, "[ionlog internal log] %s\n", err.Error())
	}

	return r
}

// Write writes the log message to the log file.
func (r *rotationEngine) Write(p []byte) (n int, err error) {
	r.logFileMu.Lock()
	defer r.logFileMu.Unlock()

	if !r.ready {
		if err := r.load(); err != nil {
			return 0, err
		}
	}
	if err := r.autoRotate(uint(len(p))); err != nil {
		r.ready = false
		return 0, err
	}
	r.autoCheckFolderSize(uint(len(p)))

	n, err = r.logFile.Write(p)
	r.files[len(r.files)-1].size += uint(n)
	r.total += uint(n)
	return n, err
}

// load asserts the folder and reads its log files.
func (r *rotationEngine) load() error {
	if err := r.assertFolder(); err != nil {
		return err
	}

	files, err := r.getAllfiles()
	if err != nil {
		return err
	}

	r.closeFileLocked()
	r.files, r.total, r.ready = files, 0, true
	for _, f := range files {
		r.total += f.size
	}
	return nil
}

func (r *rotationEngine) CloseLogFile() {
	r.closeFile()
}

// closeFile closes the log file.
func (r *rotationEngine) closeFile() {
	r.logFileMu.Lock()
	defer r.logFileMu.Unlock()
	r.closeFileLocked()
}

// closeFileLocked closes the log file. Callers must hold logFileMu.
func (r *rotationEngine) closeFileLocked() {
	if r.logFile != nil {
		if err := r.logFile.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "[ionlog internal log] Error to close current log file: %v\n", err)
		}
		r.logFile = nil
	}
}

func (r *rotationEngine) setLogFile(file *os.File) {
	if file == nil {
		fmt.Fprintf(os.Stderr, "[ionlog internal log] Cannot set the log file: file is not valid\n")
		return
	}

	r.closeFileLocked()
	r.logFile = file
}

func (r *rotationEngine) autoRotate(n uint) error {
	if now := r.now(); r.logFile != nil && now.Sub(r.checkedAt).Abs() >= time.Second {
		r.checkedAt = now
		// reload if the open file is no longer in the folder
		info, err := r.Stat(filepath.Join(r.folder, r.files[len(r.files)-1].name))
		openInfo, openErr := r.logFile.Stat()
		if err != nil || openErr != nil || !os.SameFile(info, openInfo) {
			if err := r.load(); err != nil {
				return err
			}
		}
	}

	if len(r.files) == 0 {
		return r.createNewFile(r.now(), 0)
	}

	last := r.files[len(r.files)-1]
	if r.now().After(last.date) && r.checkRotation(last.date) {
		return r.createNewFile(r.now(), 0)
	}

	if r.maxFileSize > 0 && last.size > 0 && last.size+n > r.maxFileSize {
		return r.createNewFile(last.date, last.part+1)
	}

	// no rotaion needed, check if file is open
	if r.logFile == nil {
		actualFile, err := r.OpenFile(filepath.Join(r.folder, last.name), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		r.setLogFile(actualFile)
	}
	return nil
}

func (r *rotationEngine) autoCheckFolderSize(n uint) {
	if r.maxFolderSize == NoMaxFolderSize {
		return
	}

	for r.total+n > r.maxFolderSize && len(r.files) > 1 {
		oldestFile := r.files[0]
		if err := r.RemoveFile(filepath.Join(r.folder, oldestFile.name)); err != nil && !r.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "[ionlog internal log] %s\n", err.Error())
		}
		r.files = r.files[1:]
		r.total -= oldestFile.size
	}
}
