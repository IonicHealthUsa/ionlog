package logengine

import "strconv"

type Level int

const (
	Trace Level = -2
	Debug Level = -1
	Info  Level = 0
	Warn  Level = 1
	Error Level = 2
	Panic Level = 3
	Fatal Level = 4
)

func (l Level) String() string {
	switch l {
	case Trace:
		return "TRACE"
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warn:
		return "WARN"
	case Error:
		return "ERROR"
	case Panic:
		return "PANIC"
	case Fatal:
		return "FATAL"
	default:
		return strconv.Itoa(int(l))
	}
}
