package service

import (
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/IonicHealthUsa/ionlog/internal/core/rotationengine"
)

func TestNewRotationService(t *testing.T) {
	folderName := "new_rotation_service"
	maxFolderSizes := []uint{
		rotationengine.NoMaxFolderSize,
		rotationengine.KB,
		rotationengine.MB,
		rotationengine.GB,
	}
	rotations := []rotationengine.PeriodicRotation{
		rotationengine.NoAutoRotate,
		rotationengine.Daily,
		rotationengine.Weekly,
		rotationengine.Monthly,
	}

	t.Run("should return the rotation service", func(t *testing.T) {
		for _, rotation := range rotations {
			for _, maxFolderSize := range maxFolderSizes {
				rs := NewRotationService(folderName, maxFolderSize, rotation)
				if rs == nil {
					t.Error("expected a interface of rotaion service")
				}
				if reflect.ValueOf(rs).IsNil() {
					t.Error("expected a interface of rotation service")
				}

				_rs, ok := rs.(*rotationService)
				if !ok {
					t.Error("expected a instance of rotation service")
				}

				if _rs.rotationEngine == nil {
					t.Error("expected a interface of rotation engine")
				}
				if reflect.ValueOf(_rs.rotationEngine).IsNil() {
					t.Error("expected a interface of rotaion engine")
				}

				if _, err := os.Stat(folderName); err != nil {
					t.Errorf("expected a direcory %q, but did not exist", folderName)
				}
				if err := os.RemoveAll(folderName); err != nil {
					t.Error("expected remove all file and the directory")
				}
			}
		}
	})
}

func TestRotationService(t *testing.T) {
	folderName := "rotaion_service"
	maxFolderSize := rotationengine.GB
	rotation := rotationengine.Daily

	t.Run("should return the interface of rotation service", func(t *testing.T) {
		rs := NewRotationService(folderName, maxFolderSize, rotation)
		if rs == nil {
			t.Error("expected a interface of rotaion service")
		}
		if reflect.ValueOf(rs).IsNil() {
			t.Error("expected a interface of rotation service")
		}

		re := rs.RotationEngine()
		if re == nil {
			t.Error("expected a interface of rotation engine")
		}
		if reflect.ValueOf(re).IsNil() {
			t.Error("expected a interface of rotation engine")
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a direcory %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})
}

func TestStart_rotatation(t *testing.T) {
	folderName := "rotaion_start"
	maxFolderSize := rotationengine.GB
	rotation := rotationengine.Daily

	t.Run("should start the rotate service", func(t *testing.T) {
		rs := NewRotationService(folderName, maxFolderSize, rotation)
		_rs, ok := rs.(*rotationService)
		if !ok {
			t.Error("expected a instance of rotation of rotation service")
		}

		startSync := &sync.WaitGroup{}
		startSync.Add(1)
		go rs.Start(startSync)
		startSync.Wait()

		_rs.serviceStatusLock.Lock()
		if _rs.serviceStatus != Running {
			t.Errorf("expected rotate service status to be %q, but got %q", Running, _rs.serviceStatus)
		}
		_rs.serviceStatusLock.Unlock()

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a direcory %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})
}

func TestStop_rotation(t *testing.T) {
	folderName := "rotaion_stop"
	maxFolderSize := rotationengine.GB
	rotation := rotationengine.Daily

	t.Run("should stop the rotate service", func(t *testing.T) {
		rs := NewRotationService(folderName, maxFolderSize, rotation)
		_rs, ok := rs.(*rotationService)
		if !ok {
			t.Error("expected a instance of rotation of rotation service")
		}

		startSync := &sync.WaitGroup{}
		startSync.Add(1)
		go rs.Start(startSync)
		startSync.Wait()

		_rs.serviceStatusLock.Lock()
		if _rs.serviceStatus != Running {
			t.Errorf("expected rotate service status to be %q, but got %q", Running, _rs.serviceStatus)
		}
		_rs.serviceStatusLock.Unlock()

		rs.Stop()
		time.Sleep(time.Millisecond)

		if _rs.serviceStatus != Stopped {
			t.Errorf("expected rotate service status to be %q, but got %q", Stopped, _rs.serviceStatus)
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a direcory %q, but did not exist", folderName)
		}
		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})
}

func TestStatusRotationService(t *testing.T) {
	folderName := "status"
	maxFolderSize := rotationengine.GB
	rotation := rotationengine.Daily

	t.Run("should return running status", func(t *testing.T) {
		rs := NewRotationService(folderName, maxFolderSize, rotation)
		_rs, ok := rs.(*rotationService)
		if !ok {
			t.Error("expected a instance of rotation of rotation service")
		}

		_rs.serviceStatus = Running

		if rs.Status() != Running {
			t.Errorf("expected staus to be %q, but got %q", Running, rs.Status())
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a direcory %q, but did not exist", folderName)
		}

		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})

	t.Run("should return stopped status", func(t *testing.T) {
		rs := NewRotationService(folderName, maxFolderSize, rotation)

		_rs, ok := rs.(*rotationService)
		if !ok {
			t.Error("expected a instance of rotation of rotation service")
		}

		_rs.serviceStatus = Stopped

		if rs.Status() != Stopped {
			t.Errorf("expected staus to be %q, but got %q", Stopped, rs.Status())
		}

		if _, err := os.Stat(folderName); err != nil {
			t.Errorf("expected a direcory %q, but did not exist", folderName)
		}

		if err := os.RemoveAll(folderName); err != nil {
			t.Error("expected remove all file and the directory")
		}
	})
}
