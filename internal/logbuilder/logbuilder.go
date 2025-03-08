package logbuilder

import (
	"strings"
)

type LogBuilder struct {
	fields map[string]string
	base   string
}

// NewLogBuilder creates a new logBody with initialized fields map
func NewLogBuilder() *LogBuilder {
	return &LogBuilder{
		fields: make(map[string]string),
	}
}

// AddField adds a single field
func (l *LogBuilder) AddField(key, value string) {
	l.fields[key] = value
}

// String generates the final string representation when needed
func (l *LogBuilder) String() string {
	var builder strings.Builder
	builder.WriteByte('{')

	first := true
	for key, value := range l.fields {
		if !first {
			builder.WriteByte(',')
		}
		builder.WriteByte('"')
		builder.WriteString(key)
		builder.WriteByte('"')
		builder.WriteByte(':')

		builder.WriteByte('"')
		builder.WriteString(value)
		builder.WriteByte('"')
		first = false
	}

	builder.WriteByte('}')
	builder.WriteByte('\n')
	return builder.String()
}
