// Package memory provides a way to keep track of the log history.
// It allows custom logging modes, such as logOnce and logOnChange.
package memory

import (
	"log/slog"
	"sync"

	"github.com/cespare/xxhash"
)

type RecordMode uint8

type RecordUnity struct {
	MsgHash uint64
	Mode    RecordMode
}

type RecordHistory struct {
	Records map[uint64]*RecordUnity
	mu      sync.Mutex
}

type IRecordHistory interface {
	AddRecord(id uint64, msg string, mode RecordMode) error
	RemoveRecord(id uint64)
	GetRecord(id uint64) *RecordUnity
}

const (
	LogOnce RecordMode = iota
	LogOnChange
)

func NewRecordHistory() IRecordHistory {
	return &RecordHistory{
		Records: make(map[uint64]*RecordUnity),
	}
}

func GenHash(s string) uint64 {
	return xxhash.Sum64String(s)
}

func (r *RecordHistory) AddRecord(id uint64, msg string, mode RecordMode) error {
	if r.GetRecord(id) != nil {
		return ErrRecordIDCollision
	}
	r.writeRecord(
		id,
		&RecordUnity{
			MsgHash: GenHash(msg),
			Mode:    mode,
		},
	)
	return nil
}

func (r *RecordHistory) RemoveRecord(id uint64) {
	if r.GetRecord(id) == nil {
		slog.Debug("Trying to remove non-existing record")
	}
	r.deleteRecord(id)
}

func (r *RecordHistory) GetRecord(id uint64) *RecordUnity {
	record := r.readRecord(id)
	return record
}

func (r *RecordHistory) readRecord(id uint64) *RecordUnity {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.Records[id]
}

func (r *RecordHistory) writeRecord(id uint64, req *RecordUnity) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Records[id] = req
}

func (r *RecordHistory) deleteRecord(id uint64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.Records, id)
}
