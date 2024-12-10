package recordhistory

import (
	"fmt"
)

// LogOnce allows logging only once per application execution.
// It returns true if it is the first time the message is logged.
// Otherwise, it returns false.
func LogOnce(
	logHistory *RecordHistory,
	pkg string,
	function string,
	file string,
	line int,
	msg string,
) bool {
	id := GenHash(fmt.Sprintf("%s%s%s%d", pkg, function, file, line))

	if logHistory.GetRecord(id) != nil {
		return false
	}

	logHistory.AddRecord(id, msg, logOnce)
	return true
}

// LogOnChange allows logging only when the message changes.
// It returns true if the message has changed. Otherwise, it returns false.
func LogOnChange(
	logHistory *RecordHistory,
	pkg string,
	function string,
	file string,
	line int,
	msg string,
) bool {
	id := GenHash(fmt.Sprintf("%s%s%s%d", pkg, function, file, line))

	rec := logHistory.GetRecord(id)
	if rec == nil {
		logHistory.AddRecord(id, msg, logOnChange)
		return true
	}

	msgHash := GenHash(msg)

	if rec.MsgHash != msgHash {
		rec.MsgHash = msgHash
		return true
	}

	return false
}
