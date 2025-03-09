package logbuilder

import (
	"strings"
)

type logBuilder struct {
	buf strings.Builder
}

type ILogBuilder interface {
	AddFields(args ...string)
	Compile() []byte
}

// NewLogBuilder creates a new logBody with initialized fields map
func NewLogBuilder() ILogBuilder {
	lb := &logBuilder{}
	lb.resetBuff()
	return lb
}

func (l *logBuilder) resetBuff() {
	l.buf.Reset()
	l.buf.Grow(256)
	l.buf.WriteByte('{')
}

// AddFields adds a single field
func (l *logBuilder) AddFields(args ...string) {
	if len(args)%2 != 0 {
		return
	}
	for i := 0; i < len(args); i += 2 {
		if len(l.buf.String()) > 1 {
			l.buf.WriteByte(',')
		}
		l.buf.WriteByte('"')
		l.buf.WriteString(args[i])
		l.buf.WriteByte('"')
		l.buf.WriteByte(':')

		l.buf.WriteByte('"')
		l.buf.WriteString(args[i+1])
		l.buf.WriteByte('"')
	}
}

func (l *logBuilder) Compile() []byte {
	defer l.resetBuff()
	l.buf.WriteString("}\n")

	return []byte(l.buf.String())
}
