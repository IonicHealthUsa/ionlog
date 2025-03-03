package usecases

import (
	"fmt"

	"github.com/IonicHealthUsa/ionlog/internal/infrastructure/memory"
)

func LogOnce(logHistory memory.IRecordHistory, pkg string, function string, file string, line int, msg string) bool {
	id := memory.GenHash(fmt.Sprintf("%s%s%s", pkg, function, file))

	rec := logHistory.GetRecord(id)
	if rec == nil {
		logHistory.AddRecord(id, msg)
		return true
	}

	msgHash := memory.GenHash(msg)

	if rec.GetMsgHash() != msgHash {
		rec.SetMsgHash(msgHash)
		return true
	}

	return false
}
