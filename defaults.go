package ionlog

import (
	"os"

	"github.com/IonicHealthUsa/ionlog/internal/core/rotationengine"
	"github.com/IonicHealthUsa/ionlog/internal/service"
	"github.com/IonicHealthUsa/ionlog/internal/styles"
)

const (
	Daily   = rotationengine.Daily
	Weekly  = rotationengine.Weekly
	Monthly = rotationengine.Monthly
)

const (
	NoMaxFolderSize uint = rotationengine.NoMaxFolderSize
	Kibibyte        uint = 1024
	Mebibyte        uint = 1024 * Kibibyte
	Gibibyte        uint = 1024 * Mebibyte
)

const DefaultLogFolder = "logs"

var logger = service.NewCoreService()

var DefaultOutput = os.Stdout

var CustomOutput = styles.CustomOutput
