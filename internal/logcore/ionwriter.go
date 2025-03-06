package logcore

import (
	"fmt"
	"io"
	"os"
)

type ionWriter struct {
	writeTargets []io.Writer
}

var DefaultOutput = os.Stdout

// Write writes the contents of p to all writeTargets
// This function returns no error nor the number of bytes written
func (w *ionWriter) Write(p []byte) (int, error) {
	for i, target := range w.writeTargets {
		if target == nil {
			fmt.Fprintf(os.Stderr, "Expected the %v° target to be not nil\n", i+1)
			continue
		}

		_, err := target.Write(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write to in the %v° target, error: %v\n", i+1, err)
		}
	}

	return 0, nil
}

func (w *ionWriter) SetTargets(targets ...io.Writer) {
	w.writeTargets = targets
}
