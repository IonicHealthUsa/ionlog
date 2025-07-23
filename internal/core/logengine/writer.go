package logengine

import (
	"fmt"
	"io"
	"os"
	"slices"
	"sync"
)

type ionWriter struct {
	writeLock sync.Mutex
	writers   []io.Writer
}

type IWriter interface {
	io.Writer
	AddWriter(writer ...io.Writer)
	DeleteWriter(writer ...io.Writer)
}

func NewWriter() IWriter {
	return &ionWriter{}
}

// Write writes the contents of p to all writeTargets
// This function returns no error nor the number of bytes written
func (i *ionWriter) Write(p []byte) (int, error) {
	i.writeLock.Lock()
	defer i.writeLock.Unlock()

	for index, w := range i.writers {
		if w == nil {
			fmt.Fprintf(os.Stderr, "Expected the %v° target to be not nil\n", index+1)
			continue
		}

		_, err := w.Write(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write to in the %v° target, error: %v\n", index+1, err)
			continue
		}
	}

	return 0, nil
}

func (i *ionWriter) AddWriter(writer ...io.Writer) {
	i.writeLock.Lock()
	defer i.writeLock.Unlock()
	for _, w := range writer {
		if slices.Contains(i.writers, w) {
			fmt.Fprintf(os.Stderr, "writer with the pointer %p already exists in the list of writers\n", w)
			continue
		}
		i.writers = append(i.writers, w)
	}
}

func (i *ionWriter) DeleteWriter(writer ...io.Writer) {
	i.writeLock.Lock()
	defer i.writeLock.Unlock()

	for _, wd := range writer {
		isFind := false
		for index, w := range i.writers {
			if wd == w {
				isFind = true
				i.writers = slices.Delete(i.writers, index, index+1)
				break
			}
		}
		if !isFind {
			fmt.Fprintf(os.Stderr, "writer with the pointer %q does not find", wd)
		}
	}
}
