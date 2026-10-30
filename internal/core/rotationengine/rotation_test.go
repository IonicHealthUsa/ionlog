package rotationengine

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"
)

func TestNewRotationEngine(t *testing.T) {
	folderName := "rotation_engine"
	maxFolderSize := GB
	rotation := Daily

	t.Run("should return rotation engine interface", func(t *testing.T) {
		r := NewRotationEngine(folderName, maxFolderSize, rotation)
		if r == nil {
			t.Error("expected a interface of rotation engine")
		}
		if reflect.ValueOf(r).IsNil() {
			t.Error("expected a value of rotation engine")
		}

		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instace of rotation engine")
		}

		if _r.folder != folderName {
			t.Errorf("expected folder name to be %q, but got %q", folderName, _r.folder)
		}
		if _r.maxFolderSize != maxFolderSize {
			t.Errorf("expected max folder size to be %q, but got %q", maxFolderSize, _r.maxFolderSize)
		}
		if _r.rotation != rotation {
			t.Errorf("expected rotation to be %q, but got %q", rotation, _r.rotation)
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a direcory %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})
}

func TestWrite(t *testing.T) {
	folderName := "rotation_write"
	maxFolderSize := GB
	rotation := Daily

	msg := []byte("Hello World")

	t.Run("should failure when write on nil logfile", func(t *testing.T) {
		r := NewRotationEngine(folderName, maxFolderSize, rotation)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instace of rotation engine")
		}

		_r.logFile = nil
		_r.folder = ""

		n, err := r.Write(msg)
		if n != 0 {
			t.Errorf("expected length of message to be %q, but got %q", 0, n)
		}
		if err == nil {
			t.Error("expected a error, but got nil")
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a diretory %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})

	t.Run("should write on logfile", func(t *testing.T) {
		r := NewRotationEngine(folderName, maxFolderSize, rotation)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instace of rotation engine")
		}

		_r.createNewFile(time.Now(), 0)
		if _r.logFile == nil {
			t.Errorf("new logfile not created")
		}

		n, err := r.Write(msg)
		if n != len(msg) {
			t.Errorf("expected the length of message to be %q, but got %d", len(msg), n)
		}
		if err != nil {
			t.Errorf("expected no error, but got %q", err)
		}

		fileName, err := _r.getMostRecentLogFile()
		if err != nil {
			t.Errorf("expected no error, but got %q", err)
		}

		content, err := os.ReadFile(filepath.Join(_r.folder, fileName))
		if err != nil {
			t.Errorf("expected no error, but got %q", err)
		}
		if !reflect.DeepEqual(content, msg) {
			t.Errorf("expected the message writen on file was %q, but got %q", string(msg), string(content))
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a diretory %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})
}

func TestAutoChecks(t *testing.T) {

}

func TestCloseLogFile(t *testing.T) {
	folderName := "rotaion_close"
	maxFolderSize := GB
	rotation := Daily

	t.Run("should close the logfile", func(t *testing.T) {
		r := NewRotationEngine(folderName, maxFolderSize, rotation)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instace of rotation engine")
		}

		_r.createNewFile(time.Now(), 0)
		if _r.logFile == nil {
			t.Errorf("new logfile not created")
		}

		r.CloseLogFile()
		if _r.logFile != nil {
			t.Errorf("expected the logfile to be closed")
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a diretory %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})
}

func TestCloseFile(t *testing.T) {
	folderName := "rotaion_close"
	maxFolderSize := GB
	rotation := Daily

	t.Run("should close the logfile", func(t *testing.T) {
		r := NewRotationEngine(folderName, maxFolderSize, rotation)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instace of rotation engine")
		}

		_r.createNewFile(time.Now(), 0)
		if _r.logFile == nil {
			t.Errorf("new logfile not created")
		}

		_r.closeFile()
		if _r.logFile != nil {
			t.Errorf("expected the logfile to be closed")
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a diretory %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})
}

func TestSetLogFile(t *testing.T) {
	folderName := "rotaion_setlogfile"
	maxFolderSize := GB
	rotation := Daily

	t.Run("should not set the logfile when file is nil", func(t *testing.T) {
		r := NewRotationEngine(folderName, maxFolderSize, rotation)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instace of rotation engine")
		}

		_r.logFile = nil

		_r.setLogFile(nil)
		if _r.logFile != nil {
			t.Error("expected the logfile to be not set")
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a diretory %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})

	t.Run("should set the logfile", func(t *testing.T) {
		r := NewRotationEngine(folderName, maxFolderSize, rotation)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instace of rotation engine")
		}

		_r.logFile = nil

		file := &os.File{}

		_r.setLogFile(file)
		if _r.logFile != file {
			t.Error("expected the logfile to be set")
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a diretory %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})
}

func TestAutoRotate(t *testing.T) {
	folderName := "rotaion_autorotate"
	maxFolderSize := GB
	rotation := Daily
	rotations := []PeriodicRotation{Daily, Weekly, Monthly}

	t.Run("should create a new logfile when none exists", func(t *testing.T) {
		r := NewRotationEngine(folderName, maxFolderSize, rotation)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instace of rotation engine")
		}

		_r.createNewFile(time.Now(), 0)
		fileName, err := _r.getMostRecentLogFile()
		if err != nil {
			t.Errorf("expected no error, but got %q", err)
		}

		err = os.RemoveAll(filepath.Join(_r.folder, fileName))
		if err != nil {
			t.Errorf("expected no error, but got %q", err)
		}

		_r.load()

		_r.autoRotate(0)

		if _r.logFile == nil {
			t.Error("expected the logfile set")
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a diretory %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})

	t.Run("should create a new logfile when the file need rotate", func(t *testing.T) {
		for _, rot := range rotations {
			r := NewRotationEngine(folderName, maxFolderSize, rot)
			_r, ok := r.(*rotationEngine)
			if !ok {
				t.Fatal("NewRotationEngine() did not return a instace of rotation engine")
			}

			_r.createNewFile(time.Now(), 0)
			fileName, err := _r.getMostRecentLogFile()
			if err != nil {
				t.Errorf("expected no error, but got %q", err)
			}

			err = os.RemoveAll(filepath.Join(_r.folder, fileName))
			if err != nil {
				t.Errorf("expected no error, but got %q", err)
			}

			filename := fmt.Sprintf(logFilePattern, "2000-01-01")
			filePath := filepath.Join(_r.folder, filename)

			_, err = _r.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				return
			}

			_r.load()

			_r.autoRotate(0)

			if _r.logFile == nil {
				t.Error("expected the logfile set")
			}

			fileName, err = _r.getMostRecentLogFile()
			if err != nil {
				t.Errorf("expected no error, but got %q", err)
			}

			fileDate, err := _r.getFileDate(fileName)
			if err != nil {
				t.Errorf("expected no error, but got %q", err)
			}

			if fileDate.Format(time.DateOnly) != time.Now().Format(time.DateOnly) {
				t.Errorf("expected the most recently file to be %q, but got %q", time.Now().Format(time.DateOnly), fileDate.Format(time.DateOnly))
			}

			if _, err := os.Stat(folderName); err != nil {
				t.Errorf("expected a diretory %q, but did not exist", folderName)
			}
			if err := os.RemoveAll(folderName); err != nil {
				t.Error("expected remove all file and the directory")
			}
		}
	})

	t.Run("should not rotate the file when set NoAutoRotate", func(t *testing.T) {
		r := NewRotationEngine(folderName, maxFolderSize, NoAutoRotate)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instace of rotation engine")
		}

		_r.createNewFile(time.Now(), 0)
		fileName, err := _r.getMostRecentLogFile()
		if err != nil {
			t.Errorf("expected no error, but got %q", err)
		}

		err = os.RemoveAll(filepath.Join(_r.folder, fileName))
		if err != nil {
			t.Errorf("expected no error, but got %q", err)
		}

		filename := fmt.Sprintf(logFilePattern, "2000-01-01")
		filePath := filepath.Join(_r.folder, filename)

		_, err = _r.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}

		_r.load()

		_r.autoRotate(0)

		if _r.logFile == nil {
			t.Error("expected the logfile set")
		}

		fileName, err = _r.getMostRecentLogFile()
		if err != nil {
			t.Errorf("expected no error, but got %q", err)
		}

		fileDate, err := _r.getFileDate(fileName)
		if err != nil {
			t.Errorf("expected no error, but got %q", err)
		}

		if fileDate.Format(time.DateOnly) != "2000-01-01" {
			t.Errorf("expected the most recently file to be %q, but got %q", "2000-01-01", fileDate.Format(time.DateOnly))
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a diretory %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})

	t.Run("should set the logfile with the most recent log file when logfile is nil", func(t *testing.T) {
		r := NewRotationEngine(folderName, maxFolderSize, rotation)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instace of rotation engine")
		}

		_r.logFile = nil

		_r.autoRotate(0)

		if _r.logFile == nil {
			t.Error("expected the logfile set")
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a diretory %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})

	t.Run("should not set the logfile when the file is not a valid log file", func(t *testing.T) {
		r := NewRotationEngine(folderName, maxFolderSize, NoAutoRotate)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instace of rotation engine")
		}

		_r.createNewFile(time.Now(), 0)
		fileName, err := _r.getMostRecentLogFile()
		if err != nil {
			t.Errorf("expected no error, but got %q", err)
		}

		err = os.RemoveAll(filepath.Join(_r.folder, fileName))
		if err != nil {
			t.Errorf("expected no error, but got %q", err)
		}

		filename := fmt.Sprintf(logFilePattern, "new-file")
		filePath := filepath.Join(_r.folder, filename)

		_, err = _r.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}

		_r.load()

		_r.autoRotate(0)

		if _r.logFile == nil {
			t.Error("expected the logfile set")
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a diretory %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})
}

func TestAutoCheckFolderSize(t *testing.T) {
	folderName := "rotaion_autorotate"
	maxFolderSize := GB
	rotation := Daily

	t.Run("should not create a new log file when do not have max folder size", func(t *testing.T) {
		r := NewRotationEngine(folderName, NoMaxFolderSize, rotation)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instace of rotation engine")
		}

		files, err := os.ReadDir(_r.folder)
		if err != nil {
			t.Errorf("expected no error, but got %q", err)
		}

		for _, file := range files {
			err = os.Remove(filepath.Join(_r.folder, file.Name()))
			if err != nil {
				t.Errorf("expected no error, but got %q", err)
			}
		}

		_r.logFile = nil

		_r.autoCheckFolderSize(0)

		if _r.logFile != nil {
			t.Error("expected the logfile not set")
		}

		_, err = _r.getMostRecentLogFile()
		if err != ErrNoLogFileFound {
			t.Errorf("expected error to be %q, but got %q", ErrNoLogFileFound, err)
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a diretory %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})

	t.Run("should not create a new log file when the folder not set", func(t *testing.T) {
		r := NewRotationEngine(folderName, maxFolderSize, rotation)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instace of rotation engine")
		}

		files, err := os.ReadDir(_r.folder)
		if err != nil {
			t.Errorf("expected no error, but got %q", err)
		}

		for _, file := range files {
			err = os.Remove(filepath.Join(_r.folder, file.Name()))
			if err != nil {
				t.Errorf("expected no error, but got %q", err)
			}
		}

		_r.logFile = nil
		_r.folder = ""

		_r.autoCheckFolderSize(0)

		_r.folder = folderName

		if _r.logFile != nil {
			t.Error("expected the logfile not set")
		}

		_, err = _r.getMostRecentLogFile()
		if err != ErrNoLogFileFound {
			t.Errorf("expected error to be %q, but got %q", ErrNoLogFileFound, err)
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a diretory %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})

	t.Run("should not create a new log file when the max folder size do not reached", func(t *testing.T) {
		r := NewRotationEngine(folderName, maxFolderSize, rotation)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instace of rotation engine")
		}

		files, err := os.ReadDir(_r.folder)
		if err != nil {
			t.Errorf("expected no error, but got %q", err)
		}

		for _, file := range files {
			err = os.Remove(filepath.Join(_r.folder, file.Name()))
			if err != nil {
				t.Errorf("expected no error, but got %q", err)
			}
		}

		_r.logFile = nil

		_r.autoCheckFolderSize(0)

		if _r.logFile != nil {
			t.Error("expected the logfile not set")
		}

		_, err = _r.getMostRecentLogFile()
		if err != ErrNoLogFileFound {
			t.Errorf("expected error to be %q, but got %q", ErrNoLogFileFound, err)
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a diretory %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})

	t.Run("should not create a new log file when the folder do not have any log files", func(t *testing.T) {
		r := NewRotationEngine(folderName, uint(1), rotation)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instace of rotation engine")
		}

		files, err := os.ReadDir(_r.folder)
		if err != nil {
			t.Errorf("expected no error, but got %q", err)
		}

		for _, file := range files {
			err = os.Remove(filepath.Join(_r.folder, file.Name()))
			if err != nil {
				t.Errorf("expected no error, but got %q", err)
			}
		}

		// Create a generic file to popularize the folder
		file, err := os.Create(filepath.Join(_r.folder, "HelloWorld.txt"))
		if err != nil {
			t.Errorf("expected no error, but got %q", err)
		}
		if _, err := file.Write([]byte("Hello World!!!")); err != nil {
			t.Errorf("expected no error, but got %q", err)
		}

		_r.logFile = nil

		_r.autoCheckFolderSize(0)

		if _r.logFile != nil {
			t.Error("expected the logfile not set")
		}

		_, err = _r.getMostRecentLogFile()
		if err != ErrNoLogFileFound {
			t.Errorf("expected error to be %q, but got %q", ErrNoLogFileFound, err)
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a diretory %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})

	t.Run("should remove the old file and create a new file", func(t *testing.T) {
		r := NewRotationEngine(folderName, uint(1), rotation)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instace of rotation engine")
		}

		files, err := os.ReadDir(_r.folder)
		if err != nil {
			t.Errorf("expected no error, but got %q", err)
		}

		for _, file := range files {
			err = os.Remove(filepath.Join(_r.folder, file.Name()))
			if err != nil {
				t.Errorf("expected no error, but got %q", err)
			}
		}

		fileName := fmt.Sprintf(logFilePattern, "2000-01-01")
		filePath := filepath.Join(_r.folder, fileName)

		oldFileLog, err := _r.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			t.Errorf("expected no error, but got %q", err)
		}
		if _, err := oldFileLog.Write([]byte("Hello World!!!")); err != nil {
			t.Errorf("expected no error, but got %q", err)
		}

		_r.load()
		_r.autoRotate(0)

		_r.autoCheckFolderSize(0)

		if _r.logFile == nil {
			t.Error("expected the logfile set")
		}
		if _, err := os.Stat(filePath); !os.IsNotExist(err) {
			t.Errorf("expected the old file %q to be removed", filePath)
		}

		fileName, err = _r.getMostRecentLogFile()
		if err != nil {
			t.Errorf("expected no error, but got %q", err)
		}

		fileDate, err := _r.getFileDate(fileName)
		if err != nil {
			t.Errorf("expected no error, but got %q", err)
		}

		if fileDate.Format(time.DateOnly) != time.Now().Format(time.DateOnly) {
			t.Errorf("expected the most recently file to be %q, but got %q", time.Now().Format(time.DateOnly), fileDate.Format(time.DateOnly))
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a diretory %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})
}

func TestRotationOptions(t *testing.T) {
	tests := []struct {
		name          string
		maxFolderSize uint
		opts          []Option
		want          uint
	}{
		{"default", 1000, nil, 100},
		{"max files", 1000, []Option{WithMaxFiles(4)}, 250},
		{"max file size", 1000, []Option{WithMaxFileSize(300)}, 300},
		{"max file size capped", 1000, []Option{WithMaxFileSize(5000)}, 1000},
		{"zero max file size ignored", 1000, []Option{WithMaxFileSize(0)}, 100},
		{"zero max files ignored", 1000, []Option{WithMaxFiles(0)}, 100},
		{"last max files wins", 1000, []Option{WithMaxFileSize(300), WithMaxFiles(2)}, 500},
		{"last max file size wins", 1000, []Option{WithMaxFiles(2), WithMaxFileSize(300)}, 300},
		{"at least one byte", 5, []Option{WithMaxFiles(10)}, 1},
		{"no folder limit never breaks", NoMaxFolderSize, []Option{WithMaxFiles(4)}, 0},
		{"no folder limit with max file size", NoMaxFolderSize, []Option{WithMaxFileSize(300)}, 300},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRotationEngine(t.TempDir(), tt.maxFolderSize, Daily, tt.opts...).(*rotationEngine)
			if r.maxFileSize != tt.want {
				t.Errorf("expected max file size to be %d, but got %d", tt.want, r.maxFileSize)
			}
		})
	}
}

func newTestEngine(t *testing.T, folder string, maxFolderSize uint, now *time.Time) *rotationEngine {
	t.Helper()
	r := NewRotationEngine(folder, maxFolderSize, Daily).(*rotationEngine)
	r.now = func() time.Time { return *now }
	t.Cleanup(r.CloseLogFile)
	return r
}

func seed(t *testing.T, folder, name string, size int) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(folder, name), make([]byte, size), 0644); err != nil {
		t.Fatal(err)
	}
}

func write(t *testing.T, r *rotationEngine, size int) {
	t.Helper()
	if n, err := r.Write(make([]byte, size)); err != nil || n != size {
		t.Fatalf("Write() = %d, %v", n, err)
	}
}

func logFiles(t *testing.T, r *rotationEngine) (names []string, total uint) {
	t.Helper()
	files, err := r.getAllfiles()
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range files {
		names = append(names, f.name)
		total += f.size
	}
	return names, total
}

func TestWriteRotation(t *testing.T) {
	now := time.Date(2026, 10, 26, 12, 0, 0, 0, time.Local)

	t.Run("should keep the folder within its limit", func(t *testing.T) {
		folder := t.TempDir()
		seed(t, folder, "autogenerated-2026-10-24.log", 300)
		seed(t, folder, "autogenerated-2026-10-25.log", 300)
		seed(t, folder, "autogenerated-2026-10-25.log.bak", 5000)
		r := newTestEngine(t, folder, 1000, &now)

		for i := 0; i < 100; i++ {
			write(t, r, 30)
			if _, total := logFiles(t, r); total > 1000 {
				t.Fatalf("expected at most 1000 bytes, but got %d after write %d", total, i)
			}
		}
		if names, _ := logFiles(t, r); !strings.HasPrefix(names[0], "autogenerated-2026-10-26") {
			t.Errorf("expected the older days removed, but got %v", names)
		}
		if _, err := os.Stat(filepath.Join(folder, "autogenerated-2026-10-25.log.bak")); err != nil {
			t.Errorf("expected other files untouched, but got %v", err)
		}
	})

	t.Run("should continue a full file in a new part", func(t *testing.T) {
		folder := filepath.Join(t.TempDir(), "logs")
		r := newTestEngine(t, folder, 1000, &now)

		for i := 0; i < 5; i++ {
			write(t, r, 50)
		}
		names, total := logFiles(t, r)
		want := []string{"autogenerated-2026-10-26.log", "autogenerated-2026-10-26.1.log", "autogenerated-2026-10-26.2.log"}
		if !slices.Equal(names, want) || total != 250 {
			t.Errorf("expected %v with 250 bytes, but got %v with %d", want, names, total)
		}
		for name, perm := range map[string]os.FileMode{folder: 0755, filepath.Join(folder, want[0]): 0644} {
			if info, err := os.Stat(name); err != nil || info.Mode().Perm() != perm {
				t.Errorf("expected %s with mode %v, but got %v", name, perm, info)
			}
		}
	})

	t.Run("should write an oversized message to the open file", func(t *testing.T) {
		folder := t.TempDir()
		seed(t, folder, "autogenerated-2026-10-26.log", 0)
		r := newTestEngine(t, folder, 1000, &now)

		for _, size := range []int{200, 2000} {
			write(t, r, size)
			if names, total := logFiles(t, r); len(names) != 1 || total != uint(size) {
				t.Errorf("expected one file with %d bytes, but got %v with %d", size, names, total)
			}
		}
	})

	t.Run("should not limit the folder without a max folder size", func(t *testing.T) {
		folder := t.TempDir()
		seed(t, folder, "autogenerated-2026-10-25.log", 5000)
		r := newTestEngine(t, folder, NoMaxFolderSize, &now)

		write(t, r, 1000)
		if names, total := logFiles(t, r); len(names) != 2 || total != 6000 {
			t.Errorf("expected two files with 6000 bytes, but got %v with %d", names, total)
		}
	})

	t.Run("should stop tracking a file that cannot be removed", func(t *testing.T) {
		folder := t.TempDir()
		seed(t, folder, "autogenerated-2026-10-24.log", 600)
		r := newTestEngine(t, folder, 1000, &now)
		var removed []string
		r.RemoveFile = func(name string) error {
			removed = append(removed, filepath.Base(name))
			if filepath.Base(name) == "autogenerated-2026-10-24.log" {
				return os.ErrPermission
			}
			return os.Remove(name)
		}

		write(t, r, 500)
		write(t, r, 600)
		if want := []string{"autogenerated-2026-10-24.log", "autogenerated-2026-10-26.log"}; !slices.Equal(removed, want) {
			t.Errorf("expected removals %v, but got %v", want, removed)
		}
	})

	t.Run("should keep appending to the newest file when the clock goes back", func(t *testing.T) {
		folder := t.TempDir()
		seed(t, folder, "autogenerated-2026-10-28.log", 10)
		r := newTestEngine(t, folder, 1000, &now)

		write(t, r, 2)
		if names, total := logFiles(t, r); len(names) != 1 || total != 12 {
			t.Errorf("expected the newest file with 12 bytes, but got %v with %d", names, total)
		}
	})

	t.Run("should date a new file at local midnight", func(t *testing.T) {
		now := time.Now()
		r := newTestEngine(t, t.TempDir(), 1000, &now)

		write(t, r, 2)
		y, m, d := now.Date()
		if date := r.files[0].date; !date.Equal(time.Date(y, m, d, 0, 0, 0, 0, time.Local)) {
			t.Errorf("expected the file dated at midnight, but got %v", date)
		}
	})

	t.Run("should recreate the newest file removed before a write", func(t *testing.T) {
		folder := t.TempDir()
		seed(t, folder, "autogenerated-2026-10-26.log", 10)
		r := newTestEngine(t, folder, 1000, &now)
		if err := os.Remove(filepath.Join(folder, "autogenerated-2026-10-26.log")); err != nil {
			t.Fatal(err)
		}

		write(t, r, 2)
		if names, _ := logFiles(t, r); len(names) != 1 {
			t.Errorf("expected the file recreated, but got %v", names)
		}
	})

	t.Run("should recreate a removed folder after a second", func(t *testing.T) {
		now := now
		folder := filepath.Join(t.TempDir(), "logs")
		r := newTestEngine(t, folder, 1000, &now)
		write(t, r, 2)
		write(t, r, 2)
		if err := os.RemoveAll(folder); err != nil {
			t.Fatal(err)
		}

		now = now.Add(500 * time.Millisecond)
		write(t, r, 2)
		if _, err := os.Stat(folder); !os.IsNotExist(err) {
			t.Errorf("expected no check within a second, but got %v", err)
		}

		now = now.Add(500 * time.Millisecond)
		write(t, r, 2)
		if names, total := logFiles(t, r); len(names) != 1 || total != 2 {
			t.Errorf("expected one file with 2 bytes, but got %v with %d", names, total)
		}
	})

	t.Run("should remove the oldest days first", func(t *testing.T) {
		folder := t.TempDir()
		for _, day := range []string{"23", "24", "25"} {
			seed(t, folder, "autogenerated-2026-10-"+day+".log", 300)
		}
		r := newTestEngine(t, folder, 1000, &now)

		write(t, r, 100)
		write(t, r, 100)
		want := []string{"autogenerated-2026-10-24.log", "autogenerated-2026-10-25.log", "autogenerated-2026-10-26.log", "autogenerated-2026-10-26.1.log"}
		if names, total := logFiles(t, r); !slices.Equal(names, want) || total != 800 {
			t.Errorf("expected %v with 800 bytes, but got %v with %d", want, names, total)
		}
	})

	t.Run("should list parts in order", func(t *testing.T) {
		folder := t.TempDir()
		for _, name := range []string{"autogenerated-2026-10-26.10.log", "autogenerated-2026-10-26.2.log", "autogenerated-2026-10-26.log", "autogenerated-2026-10-25.1.log"} {
			seed(t, folder, name, 1)
		}
		r := newTestEngine(t, folder, 1000, &now)

		want := []string{"autogenerated-2026-10-25.1.log", "autogenerated-2026-10-26.log", "autogenerated-2026-10-26.2.log", "autogenerated-2026-10-26.10.log"}
		if names, _ := logFiles(t, r); !slices.Equal(names, want) {
			t.Errorf("expected %v, but got %v", want, names)
		}
	})

	t.Run("should reopen a file replaced while open", func(t *testing.T) {
		now := now
		folder := t.TempDir()
		r := newTestEngine(t, folder, 1000, &now)
		write(t, r, 2)

		name := filepath.Join(folder, "autogenerated-2026-10-26.log")
		if err := os.Remove(name); err != nil {
			t.Fatal(err)
		}
		seed(t, folder, "autogenerated-2026-10-26.log", 3)
		now = now.Add(time.Second)
		write(t, r, 2)
		if info, err := os.Stat(name); err != nil || info.Size() != 5 {
			t.Errorf("expected the write appended to the new file, but got %v, %v", info, err)
		}
	})

	t.Run("should reload after a failed open", func(t *testing.T) {
		r := newTestEngine(t, t.TempDir(), 1000, &now)
		r.OpenFile = func(string, int, os.FileMode) (*os.File, error) { return nil, os.ErrPermission }

		if _, err := r.Write([]byte("a")); err == nil || r.ready {
			t.Fatalf("expected an error and a reload, but got %v", err)
		}
		r.OpenFile = os.OpenFile
		write(t, r, 2)
		if names, total := logFiles(t, r); len(names) != 1 || total != 2 {
			t.Errorf("expected one file with 2 bytes, but got %v with %d", names, total)
		}
	})
}
