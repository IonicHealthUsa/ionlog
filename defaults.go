package ionlog

import (
	"io"

	ioncore "github.com/IonicHealthUsa/ionlog/internal/logcore"
	"github.com/IonicHealthUsa/ionlog/internal/logrotation"
)

const (
	Daily   = logrotation.Daily
	Weekly  = logrotation.Weekly
	Monthly = logrotation.Monthly
)

const (
	Kibibyte = 1024
	Mebibyte = 1024 * Kibibyte
	Gibibyte = 1024 * Mebibyte
)

// DefaultOutput returns the standard output (stdout)
func DefaultOutput() io.Writer {
	return ioncore.DefaultOutput
}
