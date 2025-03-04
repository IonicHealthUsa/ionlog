package ionlog

import (
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

const DefaultLogFolder = "logs"

var DefaultOutput = ioncore.DefaultOutput
