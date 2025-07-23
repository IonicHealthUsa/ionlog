package logbuilder

import (
	"fmt"
	"os"
)

const bufsize = 1024
const maxBufsize = bufsize * 512 // 1/2 MB

type logBuilder struct {
	buf []byte
	p   uint
}

type ILogBuilder interface {
	AddFields(args ...string)
	Compile() []byte
}

// NewLogBuilder creates a new logBody with initialized fields map
func NewLogBuilder() ILogBuilder {
	lb := &logBuilder{}
	lb.buf = make([]byte, bufsize)
	lb.resetBuff()
	return lb
}

func (l *logBuilder) writeByte(b byte) {
	if len(l.buf) >= maxBufsize {
		fmt.Fprintf(os.Stderr, "logBuilder buffer is full, cannot handle more strings for this log entry.\n")
		return
	}
	if l.p == uint(len(l.buf)) {
		newBuf := make([]byte, len(l.buf)+bufsize)
		copy(newBuf, l.buf)
		l.buf = newBuf
	}
	l.buf[l.p] = b
	l.p++
}

func (l *logBuilder) writeString(str string) {
	for _, s := range []byte(str) {
		l.writeByte(s)
	}
}

func (l *logBuilder) resetBuff() {
	l.p = 0
	l.writeByte('{')
}

// AddFields adds a single field
func (l *logBuilder) AddFields(args ...string) {
	if len(args)%2 != 0 {
		return
	}
	for i := 0; i < len(args); i += 2 {
		if l.p > 1 {
			l.writeByte(',')
		}
		l.writeByte('"')
		l.writeString(args[i])
		l.writeByte('"')
		l.writeByte(':')

		l.writeByte('"')
		l.writeString(args[i+1])
		l.writeByte('"')
	}
}

func (l *logBuilder) Compile() []byte {
	defer l.resetBuff()
	l.writeString("}\n")

	return l.buf[:l.p]
}
