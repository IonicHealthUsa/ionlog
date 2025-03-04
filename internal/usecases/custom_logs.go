package usecases

import (
	"fmt"

	"github.com/IonicHealthUsa/ionlog/internal/infrastructure/memory"
)

func LogOnce(logsMemory memory.IRecordMemory, pkg string, function string, file string, msg string) bool {
	id := memory.GenHash(fmt.Sprintf("%s%s%s", pkg, function, file))

	rec := logsMemory.GetRecord(id)
	if rec == nil {
		logsMemory.AddRecord(id, msg)
		return true
	}

	msgHash := memory.GenHash(msg)

	if rec.GetMsgHash() != msgHash {
		rec.SetMsgHash(msgHash)
		return true
	}

	return false
}
