package rotationengine

type PeriodicRotation int

const (
	NoAutoRotate PeriodicRotation = iota
	Daily
	Weekly
	Monthly
)

const (
	NoMaxFolderSize uint = 0
	KB              uint = 1024
	MB              uint = 1024 * KB
	GB              uint = 1024 * MB
)
