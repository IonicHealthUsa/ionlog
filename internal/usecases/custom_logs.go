package usecases

import (
	"fmt"
	"os"

	"github.com/IonicHealthUsa/ionlog/internal/infrastructure/memory"
)

func LogOnce(logsMemory memory.IRecordMemory, msg string, args ...string) bool {
	id := memory.GenHash(fmt.Sprintf("%s%s%s", args[0], args[1], args[2]))

	rec := logsMemory.GetRecord(id)
	if rec == nil {
		err := logsMemory.AddRecord(id, msg)
		if err != nil {
			fmt.Fprint(os.Stderr, "Failed to add record to memory\n")
		}
		return true
	}

	msgHash := memory.GenHash(msg)

	if rec.GetMsgHash() != msgHash {
		rec.SetMsgHash(msgHash)
		return true
	}

	return false
}
