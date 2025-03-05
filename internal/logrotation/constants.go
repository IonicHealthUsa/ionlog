package logrotation

type PeriodicRotation int

const (
	NoAutoRotate PeriodicRotation = iota
	Daily
	Weekly
	Monthly
)

const NoMaxFolderSize uint = 0
