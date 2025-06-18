package usecases

import (
	"fmt"
	"testing"

	"github.com/IonicHealthUsa/ionlog/internal/infrastructure/memory"
)

func TestLogOnce(t *testing.T) {
	t.Run("should return true when log a message", func(t *testing.T) {
		r := memory.NewRecordMemory()

		msgHash := memory.GenHash("pkg")

		if !LogOnce(r, "pkg", "function", "file", "msg") {
			t.Errorf("LogOnce() failed")
		}

		id := memory.GenHash(fmt.Sprintf("%s%s%s", "function", "file", "msg"))
		rec := r.GetRecord(id)

		if rec.GetMsgHash() != msgHash {
			t.Errorf("expected hash message to be %q, but got %q", msgHash, rec.GetMsgHash())
		}
	})

	t.Run("should return false when two logs check same logs menssage for the same hash", func(t *testing.T) {
		r := memory.NewRecordMemory()

		msgHash := memory.GenHash("pkg")

		if !LogOnce(r, "pkg", "function", "file", "msg") {
			t.Error("expected the return to be 'true', but got 'false'")
		}

		if LogOnce(r, "pkg", "function", "file", "msg") {
			t.Error("expected the return to be 'false', but got 'true'")
		}

		id := memory.GenHash(fmt.Sprintf("%s%s%s", "function", "file", "msg"))
		rec := r.GetRecord(id)

		if rec.GetMsgHash() != msgHash {
			t.Errorf("expected hash message to be %q, but got %q", msgHash, rec.GetMsgHash())
		}
	})

	t.Run("should return true when set two different logs with same message", func(t *testing.T) {
		r := memory.NewRecordMemory()

		msgHash := memory.GenHash("pkg")

		if !LogOnce(r, "pkg", "function", "file", "msg") {
			t.Error("expected the return to be 'true', but got 'false'")
		}
		id1 := memory.GenHash(fmt.Sprintf("%s%s%s", "function", "file", "msg"))

		if !LogOnce(r, "pkg", "function", "file", "New Msg") {
			t.Error("expected the return to be 'true', but got 'false'")
		}
		id2 := memory.GenHash(fmt.Sprintf("%s%s%s", "function", "file", "New Msg"))

		rec1 := r.GetRecord(id1)

		if rec1.GetMsgHash() != msgHash {
			t.Errorf("expected hash message to be %q, but got %q", msgHash, rec1.GetMsgHash())
		}

		rec2 := r.GetRecord(id2)

		if rec2.GetMsgHash() != msgHash {
			t.Errorf("expected hash message to be %q, but got %q", msgHash, rec2.GetMsgHash())
		}
	})

	t.Run("should set the new message on the ready hash", func(t *testing.T) {
		r := memory.NewRecordMemory()

		if !LogOnce(r, "pkg", "function", "file", "msg") {
			t.Error("expected the return to be 'true', but got 'false'")
		}

		id := memory.GenHash(fmt.Sprintf("%s%s%s", "function", "file", "msg"))
		rec := r.GetRecord(id)

		msgHash := memory.GenHash("pkg")
		if msgHash != rec.GetMsgHash() {
			t.Errorf("expected hash to be %q, but got %q", msgHash, rec.GetMsgHash())
		}

		if !LogOnce(r, "gkp", "function", "file", "msg") {
			t.Error("expected the return to be 'true', but got 'false'")
		}

		msgHash = memory.GenHash("gkp")
		if msgHash != rec.GetMsgHash() {
			t.Errorf("expected hash to be %q, but got %q", msgHash, rec.GetMsgHash())
		}
	})
}
