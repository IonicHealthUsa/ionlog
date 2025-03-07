package usecases

import (
	"fmt"

	"github.com/IonicHealthUsa/ionlog/internal/infrastructure/memory"
)

func LogOnce(logsMemory memory.IRecordMemory, msg string, args ...string) bool {
	id := memory.GenHash(fmt.Sprintf("%s%s%s", args[0], args[1], args[2]))

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
