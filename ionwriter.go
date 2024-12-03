package ionlog

import (
	"io"
	"os"
)

type ionWriter struct {
	writeTargets []io.Writer
}

// Write writes the contents of p to all writeTargets
// if any writeTarget returns an error, the error is returned with the number of bytes written
// if no writeTarget returns an error, this function returns no err = nil and n = 0
func (w *ionWriter) Write(p []byte) (n int, err error) {
	for _, target := range w.writeTargets {

		if target == nil {
			// TODO: O que fazer? tirar? continuar?
		}

		// It will save the latest failure error while continue writing to other writeTargets
		// latter, it will return the latest failure error
		_n, _err := target.Write(p)
		if _err != nil {
			n = _n
			err = _err

			// TODO: Log this errors
		}
	}
	return
}

func Stdout() io.Writer {
	return os.Stdout
}
