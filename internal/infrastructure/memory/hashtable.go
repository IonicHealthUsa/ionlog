// Package memory provides a way to keep track of the log history.
// It allows custom logging modes, such as logOnce and logOnChange.
package memory

import (
	"log/slog"
	"sync"

	"github.com/cespare/xxhash"
)

type recordUnity struct {
	MsgHash uint64
}

type recordHistory struct {
	records map[uint64]*recordUnity
	mu      sync.Mutex
}

type IRecordUnity interface {
	GetMsgHash() uint64
	SetMsgHash(msg uint64)
}

type IRecordHistory interface {
	AddRecord(id uint64, msg string) error
	RemoveRecord(id uint64)
	GetRecord(id uint64) IRecordUnity
}

func NewRecordHistory() IRecordHistory {
	return &recordHistory{
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

func (r *recordHistory) AddRecord(id uint64, msg string) error {
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

func (r *recordHistory) RemoveRecord(id uint64) {
	if r.GetRecord(id) == nil {
		slog.Debug("Trying to remove non-existing record")
		return
	}
	r.deleteRecord(id)
}

func (r *recordHistory) GetRecord(id uint64) IRecordUnity {
	record := r.readRecord(id)
	if record == nil {
		return nil // yeah, it have to be like this.
	}
	return record
}

func (r *recordHistory) readRecord(id uint64) *recordUnity {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.records[id]
}

func (r *recordHistory) writeRecord(id uint64, req *recordUnity) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.records[id] = req
}

func (r *recordHistory) deleteRecord(id uint64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.records, id)
}
