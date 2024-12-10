package recordhistory

import (
	"github.com/cespare/xxhash"
)

type RecordMode int

type RecordUnity struct {
	ID      uint64
	MsgHash uint64
	Mode    RecordMode
}

type RecordHistory struct {
	Records []RecordUnity
}

const (
	logOnce RecordMode = iota
	logOnChange
)

func NewRecordHistory() *RecordHistory {
	return &RecordHistory{}
}

func GenHash(s string) uint64 {
	return xxhash.Sum64String(s)
}

func (r *RecordHistory) AddRecord(id uint64, msg string, mode RecordMode) error {
	if r.GetRecord(id) != nil {
		return ErrRecordIDCollision
	}

	r.Records = append(r.Records, RecordUnity{
		ID:      id,
		MsgHash: GenHash(msg),
		Mode:    mode,
	})
	return nil
}

func (r *RecordHistory) RemoveRecord(id uint64) {
	for i, rec := range r.Records {
		if rec.ID == id {
			r.Records = append(r.Records[:i], r.Records[i+1:]...)
			break
		}
	}
}

func (r *RecordHistory) GetRecord(id uint64) *RecordUnity {
	for i := 0; i < len(r.Records); i++ {
		rec := &r.Records[i]
		if rec.ID == id {
			return rec
		}
	}
	return nil
}
