package logcore

import (
	"fmt"
	"io"
	"os"
)

type ionWriter struct {
	targets []io.Writer
}

var DefaultOutput = os.Stdout

// Write writes the contents of p to all writeTargets
// This function returns no error nor the number of bytes written
func (i *ionWriter) Write(p []byte) (int, error) {
	for index, t := range i.targets {
		if t == nil {
			fmt.Fprintf(os.Stderr, "Expected the %v° target to be not nil\n", index+1)
			continue
		}

		_, err := t.Write(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write to in the %v° target, error: %v\n", index+1, err)
			continue
		}
	}

	return 0, nil
}

func (i *ionWriter) SetTargets(writers ...io.Writer) {
	i.targets = writers
}

func (i *ionWriter) AddTarget(writer io.Writer) {
	i.targets = append(i.targets, writer)
}
