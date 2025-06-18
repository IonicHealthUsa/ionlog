package rotationengine

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"
)

func TestGetFileDate(t *testing.T) {
	folderName := "utils_getfiledate"
	maxFolderSize := GB
	rotation := Daily

	t.Run("should return a error when the file can not parse the time", func(t *testing.T) {
		r := NewRotationEngine(folderName, maxFolderSize, rotation)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instance of rotation engine")
		}

		fileName := fmt.Sprintf(logFilePattern, "helloworld")

		_, err := _r.getFileDate(fileName)
		if err == nil {
			t.Error("expected a error, but got nil")
		}
	})

	t.Run("should return the correct time", func(t *testing.T) {
		r := NewRotationEngine(folderName, maxFolderSize, rotation)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instance of rotation engine")
		}

		expectedTime := time.Now().Format(time.DateOnly)
		fileName := fmt.Sprintf(logFilePattern, expectedTime)

		gotTime, err := _r.getFileDate(fileName)
		if err != nil {
			t.Errorf("expected no error, but got %q", err)
		}
		if gotTime.Format(time.DateOnly) != expectedTime {
			t.Errorf("expected time to be %q, but got %q", expectedTime, gotTime.Format(time.DateOnly))
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a direcory %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})
}

func TestGetAllFiles(t *testing.T) {
	folderName := "utils_getallfiles"
	maxFolderSize := GB
	rotation := Daily

	t.Run("should failure when the folder not set", func(t *testing.T) {
		r := NewRotationEngine(folderName, maxFolderSize, rotation)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instance of rotation engine")
		}

		_r.folder = ""

		files, err := _r.getAllfiles()

		if err == nil {
			t.Error("expected a error, but got nil")
		}
		if files != nil {
			t.Errorf("expected no files, but got %q", files)
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a direcory %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})

	t.Run("should return only the log file names", func(t *testing.T) {
		r := NewRotationEngine(folderName, maxFolderSize, rotation)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instance of rotation engine")
		}

		wrongLogFile := "helloWorld.log"
		_, err := os.OpenFile(filepath.Join(_r.folder, wrongLogFile), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			t.Errorf("expected no error, but got %q", err)
		}

		otherFile := "ionlog.txt"
		_, err = _r.OpenFile(filepath.Join(_r.folder, otherFile), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			t.Errorf("expected no error, but got %q", err)
		}

		otherFolder := "ionic_health"
		_ = os.Mkdir(filepath.Join(_r.folder, otherFolder), 0666)

		files, err := _r.getAllfiles()

		if err != nil {
			t.Errorf("expected no error, but got %q", err)
		}

		expectedTime := time.Now().Format(time.DateOnly)
		expectFileName := fmt.Sprintf(logFilePattern, expectedTime)

		if !slices.Contains(files, expectFileName) {
			t.Errorf("expected the folder contains %q, but got %q", expectFileName, files)
		}
		if slices.Contains(files, wrongLogFile) {
			t.Errorf("expected the function no return %q, but got", wrongLogFile)
		}
		if slices.Contains(files, otherFile) {
			t.Errorf("expected the function no return %q, but got", otherFile)
		}
		if slices.Contains(files, otherFolder) {
			t.Errorf("expected the function no return %q, but got", otherFile)
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a direcory %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})
}

func TestGetMostRecentLogFile(t *testing.T) {
	folderName := "utils_getmostrecentlogfile"
	maxFolderSize := GB
	rotation := Daily

	t.Run("should failure when the folder not set", func(t *testing.T) {
		r := NewRotationEngine(folderName, maxFolderSize, rotation)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instance of rotation engine")
		}

		_r.folder = ""

		file, err := _r.getMostRecentLogFile()
		if err == nil {
			t.Error("expected a error, but got nil")
		}
		if file != "" {
			t.Errorf("expected no file, but got %q", file)
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a direcory %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})

	t.Run("should return a error when the folder do not have log files", func(t *testing.T) {
		r := NewRotationEngine(folderName, maxFolderSize, rotation)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instance of rotation engine")
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

		file, err := _r.getMostRecentLogFile()
		if err != ErrNoLogFileFound {
			t.Errorf("expected error to be %q, but got %q", ErrNoLogFileFound, err)
		}
		if file != "" {
			t.Errorf("expected no file, but got %q", file)
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a direcory %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})

	t.Run("should return the most recent log file", func(t *testing.T) {
		r := NewRotationEngine(folderName, maxFolderSize, rotation)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instance of rotation engine")
		}

		fileName := fmt.Sprintf(logFilePattern, "2000-01-01")
		_ = os.NewFile(0666, filepath.Join(_r.folder, fileName))
		fileName = fmt.Sprintf(logFilePattern, "helloworld")
		_ = os.NewFile(0666, filepath.Join(_r.folder, fileName))

		filename := fmt.Sprintf(logFilePattern, time.Now().Format(time.DateOnly))

		file, err := _r.getMostRecentLogFile()
		if err != nil {
			t.Errorf("expected no error, but got %q", err)
		}
		if file != filename {
			t.Errorf("expected most recent log file to be %q, but got %q", filename, file)
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a direcory %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})
}

func TestCreateNewFile(t *testing.T) {
	folderName := "utils_createnewfile"
	maxFolderSize := GB
	rotation := Daily

	t.Run("should failure when the folder not set", func(t *testing.T) {
		r := NewRotationEngine(folderName, maxFolderSize, rotation)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instance of rotation engine")
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

		_r.createNewFile()

		filename := fmt.Sprintf(logFilePattern, time.Now().Format(time.DateOnly))

		files, err = os.ReadDir(folderName)
		if err != nil {
			t.Errorf("expected no error, but got %q", err)
		}

		var filenames = make([]string, 0, len(files))
		for _, file := range files {
			filenames = append(filenames, file.Name())
		}

		if slices.Contains(filenames, filename) {
			t.Errorf("expected the file %q do not exist on the folder", filename)
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a direcory %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})

	t.Run("should create a new log file", func(t *testing.T) {
		r := NewRotationEngine(folderName, maxFolderSize, rotation)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instance of rotation engine")
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

		_r.createNewFile()

		if _r.logFile == nil {
			t.Error("expected log file to be set, but is nil")
		}

		expectedFileName := fmt.Sprintf(logFilePattern, time.Now().Format(time.DateOnly))

		file, err := _r.getMostRecentLogFile()
		if err != nil {
			t.Errorf("expected no error, but got %q", err)
		}
		if file != expectedFileName {
			t.Errorf("expected most recent log file to be %q, but got %q", expectedFileName, file)
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a direcory %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})
}

func TestAssertFolder(t *testing.T) {
	folderName := "utils_assertfolder"
	maxFolderSize := GB
	rotation := Daily

	t.Run("should return nil when the folder exists", func(t *testing.T) {
		r := NewRotationEngine(folderName, maxFolderSize, rotation)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instance of rotation engine")
		}

		if err := _r.assertFolder(); err != nil {
			t.Errorf("expected no error, but got %q", err)
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a direcory %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})

	t.Run("should create a folder when does not exists", func(t *testing.T) {
		r := NewRotationEngine(folderName, maxFolderSize, rotation)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instance of rotation engine")
		}

		if err := _r.assertFolder(); err != nil {
			t.Errorf("expected no error, but got %q", err)
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a folder %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})
}

func TestCheckRotation(t *testing.T) {
	folderName := "utils_checkrotation"
	maxFolderSize := GB

	t.Run("should return true when need rotate for daily rotation", func(t *testing.T) {
		r := NewRotationEngine(folderName, maxFolderSize, Daily)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instance of rotation engine")
		}

		t.Run("different year", func(t *testing.T) {
			dateStr := "2000-01-01"
			fileDate, err := time.Parse(time.DateOnly, dateStr)
			if err != nil {
				t.Errorf("expected no error, but got %q", err)
			}

			if !_r.checkRotation(fileDate) {
				t.Error("expected the return to be true, but got false")
			}
		})

		t.Run("same year but different month", func(t *testing.T) {
			timeNow := time.Now()
			month := timeNow.Month()
			if month == 12 {
				month -= 1
			} else {
				month += 1
			}
			fileDate := time.Date(timeNow.Year(), month, 1, 0, 0, 0, 0, time.UTC)

			if !_r.checkRotation(fileDate) {
				t.Error("expected the return to be true, but got false")
			}
		})

		t.Run("same year and month but different day", func(t *testing.T) {
			timeNow := time.Now()
			day := timeNow.Day()
			if day >= 28 {
				day -= 1
			} else {
				day += 1
			}
			fileDate := time.Date(timeNow.Year(), timeNow.Month(), day, 0, 0, 0, 0, time.UTC)

			if !_r.checkRotation(fileDate) {
				t.Error("expected the return to be true, but got false")
			}
		})

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a folder %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})

	t.Run("should return false when not need rotate for daily rotation", func(t *testing.T) {
		r := NewRotationEngine(folderName, maxFolderSize, Daily)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instance of rotation engine")
		}

		if _r.checkRotation(time.Now()) {
			t.Error("expected the return to be false, but got true")
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a folder %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})

	t.Run("should return true when need rotate for weekly rotation", func(t *testing.T) {
		r := NewRotationEngine(folderName, maxFolderSize, Daily)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instance of rotation engine")
		}

		t.Run("different year", func(t *testing.T) {
			dateStr := "2000-01-01"
			fileDate, err := time.Parse(time.DateOnly, dateStr)
			if err != nil {
				t.Errorf("expected no error, but got %q", err)
			}

			if !_r.checkRotation(fileDate) {
				t.Error("expected the return to be true, but got false")
			}
		})

		t.Run("same year but different week", func(t *testing.T) {
			timeNow := time.Now()
			y, w := timeNow.UTC().ISOWeek()
			if w >= 52 {
				w -= 1
			} else {
				w += 1
			}

			firstDay := time.Date(y, time.January, 1, 0, 0, 0, 0, time.UTC)
			dayWeek := int(firstDay.Weekday())
			if dayWeek == 0 {
				dayWeek = 7
			}
			fitDay := 1 - dayWeek
			firstDay = firstDay.AddDate(0, 0, fitDay)

			fileDate := firstDay.AddDate(0, 0, w*7)

			if !_r.checkRotation(fileDate) {
				t.Error("expected the return to be true, but got false")
			}
		})

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a folder %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})

	t.Run("should return false when not need rotate for weekly rotation", func(t *testing.T) {
		r := NewRotationEngine(folderName, maxFolderSize, Weekly)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instance of rotation engine")
		}

		if _r.checkRotation(time.Now()) {
			t.Error("expected the return to be false, but got true")
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a folder %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})

	t.Run("should return true when need rotate for weekly rotation", func(t *testing.T) {
		r := NewRotationEngine(folderName, maxFolderSize, Monthly)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instance of rotation engine")
		}

		t.Run("different year", func(t *testing.T) {
			dateStr := "2000-01-01"
			fileDate, err := time.Parse(time.DateOnly, dateStr)
			if err != nil {
				t.Errorf("expected no error, but got %q", err)
			}

			if !_r.checkRotation(fileDate) {
				t.Error("expected the return to be true, but got false")
			}
		})

		t.Run("same year but different month", func(t *testing.T) {
			timeNow := time.Now()
			month := timeNow.Month()
			if month == 12 {
				month -= 1
			} else {
				month += 1
			}
			fileDate := time.Date(timeNow.Year(), month, 1, 0, 0, 0, 0, time.UTC)

			if !_r.checkRotation(fileDate) {
				t.Error("expected the return to be true, but got false")
			}
		})

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a folder %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}

	})

	t.Run("should return false when not need rotate for weekly rotation", func(t *testing.T) {
		r := NewRotationEngine(folderName, maxFolderSize, Weekly)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instance of rotation engine")
		}

		t.Run("same year and month but different day", func(t *testing.T) {
			timeNow := time.Now()
			day := timeNow.Day()
			if day >= 28 {
				day -= 1
			} else {
				day += 1
			}
			fileDate := time.Date(timeNow.Year(), timeNow.Month(), day, 0, 0, 0, 0, time.UTC)

			if _r.checkRotation(fileDate) {
				t.Error("expected the return to be true, but got false")
			}
		})

		t.Run("same day", func(t *testing.T) {
			if _r.checkRotation(time.Now()) {
				t.Error("expected the return to be false, but got true")
			}
		})

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a folder %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})
}

func TestGetOldestLogFile(t *testing.T) {
	folderName := "utils_getolderlogfile"
	maxFolderSize := GB
	rotation := Daily

	t.Run("should failure when the folder not exist", func(t *testing.T) {
		r := NewRotationEngine(folderName, maxFolderSize, rotation)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instance of rotation engine")
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a folder %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}

		if _, err := _r.getOldestLogFile(); err == nil {
			t.Error("expected a error, but got nil")
		}
	})

	t.Run("should return the error when no exist log files", func(t *testing.T) {
		r := NewRotationEngine(folderName, maxFolderSize, rotation)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instance of rotation engine")
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

		if _, err := _r.getOldestLogFile(); err != ErrNoLogFileFound {
			t.Errorf("expected error to be %q, but got %q", ErrNoLogFileFound, err)
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a folder %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})

	t.Run("should return the oldest lof file", func(t *testing.T) {
		r := NewRotationEngine(folderName, maxFolderSize, rotation)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instance of rotation engine")
		}

		expectedOldestFile := fmt.Sprintf(logFilePattern, "2000-01-01")
		filePath := filepath.Join(_r.folder, expectedOldestFile)
		_, err := _r.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			t.Errorf("expected no error, but got %q", err)
		}

		oldestFile, err := _r.getOldestLogFile()
		if err != nil {
			t.Errorf("expected no error, but got %q", err)
		}
		if oldestFile != expectedOldestFile {
			t.Errorf("expected the oldest log file to be %q, but got %q", expectedOldestFile, oldestFile)
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a folder %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})
}

func TestGetFolderSize(t *testing.T) {
	folderName := "utils_getfoldersize"
	maxFolderSize := GB
	rotation := Daily

	t.Run("should failure when the folder not exist", func(t *testing.T) {
		r := NewRotationEngine(folderName, maxFolderSize, rotation)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instance of rotation engine")
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a folder %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}

		if _, err := _r.getFolderSize(); err == nil {
			t.Error("expected a error, but got nil")
		}
	})

	t.Run("should return the folder size", func(t *testing.T) {
		r := NewRotationEngine(folderName, maxFolderSize, rotation)
		_r, ok := r.(*rotationEngine)
		if !ok {
			t.Fatal("NewRotationEngine() did not return a instance of rotation engine")
		}

		size, err := _r.getFolderSize()
		if err != nil {
			t.Error("expected a error, but got nil")
		}
		if size != 0 {
			t.Errorf("expected error to be 0, but got %d", size)
		}

		if _, err = _r.logFile.Write([]byte("Hello World")); err != nil {
			t.Errorf("expected no error, but got %q", err)
		}
		size, err = _r.getFolderSize()
		if err != nil {
			t.Error("expected a error, but got nil")
		}
		if size != 11 {
			t.Errorf("expected error to be 11, but got %d", size)
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a folder %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})
}
