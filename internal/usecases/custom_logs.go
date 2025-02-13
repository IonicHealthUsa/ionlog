package usecases

import (
	"fmt"

	"github.com/IonicHealthUsa/ionlog/internal/infrastructure/memory"
)

// LogOnce allows logging only once per application execution.
// It returns true if it is the first time the message is logged.
// Otherwise, it returns false.
func LogOnce(
	logHistory memory.IRecordHistory,
	pkg string,
	function string,
	file string,
	line int,
	msg string,
) bool {
	id := memory.GenHash(fmt.Sprintf("%s%s%s%d", pkg, function, file, line))

	if logHistory.GetRecord(id) != nil {
		return false
	}

	logHistory.AddRecord(id, msg, memory.LogOnce)
	return true
}

// LogOnChange allows logging only when the message changes.
// It returns true if the message has changed. Otherwise, it returns false.
func LogOnChange(
	logHistory memory.IRecordHistory,
	pkg string,
	function string,
	file string,
	line int,
	msg string,
) bool {
	id := memory.GenHash(fmt.Sprintf("%s%s%s", pkg, function, file))

	rec := logHistory.GetRecord(id)
	if rec == nil {
		logHistory.AddRecord(id, msg, memory.LogOnChange)
		return true
	}

	msgHash := memory.GenHash(msg)

	if rec.MsgHash != msgHash {
		rec.MsgHash = msgHash
		return true
	}

	return false
}
