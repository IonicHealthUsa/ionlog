package filesystem

import (
	"os"
	"testing"
)

type mockOS struct{}

func (m *mockOS) Stat(string) (os.FileInfo, error) {
	return nil, nil
}

func (m *mockOS) Mkdir(string, os.FileMode) error {
	return nil
}

func (m *mockOS) ReadDir(string) ([]os.DirEntry, error) {
	return nil, nil
}

func (m *mockOS) IsNotExist(error) bool {
	return true
}

func (m *mockOS) OpenFile(name string, flag int, perm os.FileMode) (*os.File, error) {
	return nil, nil
}

func (m *mockOS) RemoveFile(name string) error {
	return nil
}

func TestFileSystem(t *testing.T) {
	t.Run("should received the same functions", func(t *testing.T) {
		m := &mockOS{}

		fs := NewFileSystem(
			m.Stat,
			m.Mkdir,
			m.ReadDir,
			m.IsNotExist,
			m.OpenFile,
			m.RemoveFile,
		)

		if fs.Stat == nil {
			t.Error("expected Stat function to be set, but got nil")
			return
		}

		if file, err := fs.Stat("hello world"); file != nil && err != nil {
			t.Errorf("expected return file=nil and err=nil of Stat function, but got file=%q and err=%q", file, err)
		}

		if fs.Mkdir == nil {
			t.Error("expected Mkdir function to be set, but got nil")
			return
		}

		if err := fs.Mkdir("hello world", 0777); err != nil {
			t.Errorf("expected return err=nil of Mkdir function, but got err=%q", err)
		}

		if fs.ReadDir == nil {
			t.Error("expected ReadDir function to be set, but got nil")
			return
		}

		if dir, err := fs.ReadDir("hello world"); dir != nil && err != nil {
			t.Errorf("expected return dir=nil and err=nil, but got dir=%q and err=%q", dir, err)
		}

		if fs.IsNotExist == nil {
			t.Error("expected IsNotExist function to be set, but got nil")
			return
		}

		if b := fs.IsNotExist(nil); b != true {
			t.Errorf("expected return b=true of IsNotExist function, but got b='%v'", b)
		}

		if fs.OpenFile == nil {
			t.Error("expected OpenFile function to be set, but got nil")
			return
		}

		if file, err := fs.OpenFile("hello world", 0, 0777); file != nil && err != nil {
			t.Errorf("expected return file=nil and err=nil OpenFile, but got file=%v and err=%q", file, err)
		}

		if fs.RemoveFile == nil {
			t.Error("expected RemoveFile function to be set, but got nil")
			return
		}

		if err := fs.RemoveFile("hello world"); err != nil {
			t.Errorf("expected return err=nil of RemoveFile, but got err=%q", err)
		}
	})
}
