// Package memory provides a way to keep track of the log history.
package memory

import (
	"log/slog"
	"sync"

	"github.com/cespare/xxhash"
)

type recordUnity struct {
	MsgHash uint64
}

type recordMemory struct {
	records map[uint64]*recordUnity
	mu      sync.Mutex
}

type IRecordUnity interface {
	GetMsgHash() uint64
	SetMsgHash(msg uint64)
}

type IRecordMemory interface {
	AddRecord(id uint64, msg string) error
	RemoveRecord(id uint64)
	GetRecord(id uint64) IRecordUnity
}

func NewRecordMemory() IRecordMemory {
	return &recordMemory{
		records: make(map[uint64]*recordUnity),
	}
}

func (r recordUnity) GetMsgHash() uint64 {
	return r.MsgHash
}

func (r *recordUnity) SetMsgHash(msg uint64) {
	r.MsgHash = msg
}

func GenHash(s string) uint64 {
	return xxhash.Sum64String(s)
}

func (r *recordMemory) AddRecord(id uint64, msg string) error {
	if r.readRecord(id) != nil {
		return ErrRecordIDCollision
	}
	r.writeRecord(
		id,
		&recordUnity{
			MsgHash: GenHash(msg),
		},
	)
	return nil
}

func (r *recordMemory) RemoveRecord(id uint64) {
	if r.GetRecord(id) == nil {
		slog.Debug("Trying to remove non-existing record")
		return
	}
	r.deleteRecord(id)
}

func (r *recordMemory) GetRecord(id uint64) IRecordUnity {
	record := r.readRecord(id)
	if record == nil {
		// accessed in 21/03/2025:
		// check: https://trstringer.com/go-nil-interface-and-interface-with-nil-concrete-value/
		return nil // yeah, it have to be like this.
	}
	return record
}

func (r *recordMemory) readRecord(id uint64) *recordUnity {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.records[id]
}

func (r *recordMemory) writeRecord(id uint64, req *recordUnity) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.records[id] = req
}

func (r *recordMemory) deleteRecord(id uint64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.records, id)
}
