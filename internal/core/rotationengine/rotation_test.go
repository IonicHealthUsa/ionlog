package rotationengine

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
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

		n, err := r.Write(msg)
		if n != 0 {
			t.Errorf("expected length of message to be %q, but got %q", 0, n)
		}
		if err != ErrLogFileNotSet {
			t.Errorf("expected error to be %q, but got %q", ErrLogFileNotSet, err)
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

		_r.createNewFile()
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

		_r.createNewFile()
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

		_r.createNewFile()
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

type mockWriteCloser struct {
	*bufio.Writer
}

func (m *mockWriteCloser) Write(p []byte) (int, error) {
	return 0, nil
}

func (m *mockWriteCloser) Close() error {
	return nil
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

		file := &mockWriteCloser{}

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

		fileName, err := _r.getMostRecentLogFile()
		if err != nil {
			t.Errorf("expected no error, but got %q", err)
		}

		err = os.RemoveAll(filepath.Join(_r.folder, fileName))
		if err != nil {
			t.Errorf("expected no error, but got %q", err)
		}

		_r.logFile = nil

		_r.autoRotate()

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

			_r.logFile = nil

			_r.autoRotate()

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

		_r.logFile = nil

		_r.autoRotate()

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

		_r.autoRotate()

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

		_r.logFile = nil

		_r.autoRotate()

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

		_r.autoCheckFolderSize()

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

		_r.autoCheckFolderSize()

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

		_r.autoCheckFolderSize()

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

		_r.autoCheckFolderSize()

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

		_r.logFile = nil

		_r.autoCheckFolderSize()

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
	})
}
