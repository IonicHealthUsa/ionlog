package memory

import (
	"testing"
	"time"
)

func TestNewRecordMemory(t *testing.T) {
	r := NewRecordMemory()
	if r == nil {
		t.Errorf("NewRecordMemory() failed")
	}
	if _, ok := r.(*recordMemory); !ok {
		t.Errorf("NewRecordMemory() failed")
	}
}

func TestRecordUnityGetMsgHash(t *testing.T) {
	r := recordUnity{
		MsgHash: 1,
	}
	if r.GetMsgHash() != 1 {
		t.Errorf("GetMsgHash() failed")
	}
}

func TestRecordUnitySetMsgHash(t *testing.T) {
	r := recordUnity{}
	r.SetMsgHash(1)
	if r.MsgHash != 1 {
		t.Errorf("SetMsgHash() failed")
	}
}

func TestGenHash(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want uint64
	}{
		{
			name: "TestGenHash",
			s:    "test",
			want: 5754696928334414137,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenHash(tt.s); got != tt.want {
				t.Errorf("GenHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddRecord(t *testing.T) {
	msg := "test"
	msgHash := GenHash(msg)

	t.Run("Simple Add", func(t *testing.T) {
		r := NewRecordMemory()
		_r, ok := r.(*recordMemory)
		if !ok {
			t.Fatal("NewRecordMemory did not returned a instace of record memory")
		}

		if err := r.AddRecord(1, msg); err != nil {
			t.Errorf("AddRecord() failed")
		}

		if _r.records[1] == nil {
			t.Error("exected readRecord return a instace of record unity, but got nil")
			return
		}

		if _r.records[1].MsgHash != msgHash {
			t.Errorf("expected message hash to be %q, but got %q", msgHash, _r.records[1].MsgHash)
		}
	})

	t.Run("Collision Check", func(t *testing.T) {
		r := NewRecordMemory()
		_r, ok := r.(*recordMemory)
		if !ok {
			t.Fatal("NewRecordMemory did not returned a instace of record memory")
		}

		if err := r.AddRecord(1, msg); err != nil {
			t.Errorf("AddRecord() failed")
		}

		if _r.records[1] == nil {
			t.Error("exected readRecord return a instace of record unity, but got nil")
			return
		}

		if _r.records[1].MsgHash != msgHash {
			t.Errorf("expected message hash to be %q, but got %q", msgHash, _r.records[1].MsgHash)
		}

		if err := r.AddRecord(1, msg); err != ErrRecordIDCollision {
			t.Errorf("AddRecord() failed; Expected collision error")
		}
	})
}

func TestRemoveRecord(t *testing.T) {
	t.Run("Simple Remove", func(t *testing.T) {
		id := uint64(1)
		r := NewRecordMemory()
		r.AddRecord(id, "test")

		if r.GetRecord(id) == nil {
			t.Errorf("Test preset failed")
		}

		r.RemoveRecord(id)

		if r.GetRecord(id) != nil {
			t.Errorf("RemoveRecord() failed")
		}
	})
}

func TestGetRecord(t *testing.T) {
	t.Run("should return the interface of record unity", func(t *testing.T) {
		r := NewRecordMemory()
		r.AddRecord(1, "")
		if r.GetRecord(1) == nil {
			t.Error("expected GetRecord() return a interface of record unity, but got nil")
		}
	})

	t.Run("should return nil when unity do not exist", func(t *testing.T) {
		r := NewRecordMemory()

		if r.GetRecord(1) != nil {
			t.Error("expected GetRecord() return nil, but got a interface")
		}
	})
}

func TestReadRecord(t *testing.T) {
	msg := "hello world"
	msgHash := GenHash(msg)

	t.Run("ReadRecord", func(t *testing.T) {
		_r := NewRecordMemory()
		r, ok := _r.(*recordMemory)
		if !ok {
			t.Fatal("NewRecordMemory did not returned a instace of record memory")
		}

		r.records[1] = &recordUnity{MsgHash: msgHash}
		u := r.readRecord(1)

		if u == nil {
			t.Error("exected readRecord return a instace of record unity, but got nil")
			return
		}

		if u.MsgHash != msgHash {
			t.Errorf("expected message hash to be %q, but got %q", msgHash, u.MsgHash)
		}
	})

	t.Run("should timeout when mutex is locked", func(t *testing.T) {
		_r := NewRecordMemory()
		r, ok := _r.(*recordMemory)
		if !ok {
			t.Fatal("NewRecordMemory did not returned a instace of record memory")
		}

		r.records[1] = &recordUnity{MsgHash: msgHash}

		r.mu.Lock()
		comm := make(chan *recordUnity, 1)
		go func(chan *recordUnity) {
			comm <- r.readRecord(1)
		}(comm)

		select {
		case <-comm:
			t.Error("expected no receive a instance of record unity")
		case <-time.After(500 * time.Microsecond):
		}

		r.mu.Unlock()
		select {
		case u := <-comm:
			if u == nil {
				t.Error("exected readRecord return a instace of record unity, but got nil")
				return
			}

			if u.MsgHash != msgHash {
				t.Errorf("expected message hash to be %q, but got %q", msgHash, u.MsgHash)
			}
		case <-time.After(500 * time.Microsecond):
			t.Error("expected receive a instance of record unity, but timeout")
		}
	})
}

func TestWriteRecord(t *testing.T) {
	msg := "hello world"
	msgHash := GenHash(msg)

	t.Run("should write the record unity on record", func(t *testing.T) {
		_r := NewRecordMemory()
		r, ok := _r.(*recordMemory)
		if !ok {
			t.Fatal("NewRecordMemory did not returned a instace of record memory")
		}

		r.writeRecord(1, &recordUnity{MsgHash: msgHash})
		u := r.records[1]

		if u == nil {
			t.Error("expected writeRecord() set a record unity on records")
			return
		}

		if u.MsgHash != msgHash {
			t.Errorf("expected message hash to be %q, but got %q", msgHash, u.MsgHash)
		}
	})

	t.Run("should timeout when mutex is locked", func(t *testing.T) {
		_r := NewRecordMemory()
		r, ok := _r.(*recordMemory)
		if !ok {
			t.Fatal("NewRecordMemory did not returned a instace of record memory")
		}

		r.mu.Lock()
		go func() {
			r.writeRecord(1, &recordUnity{MsgHash: msgHash})
		}()

		time.Sleep(500 * time.Microsecond)
		if u := r.records[1]; u != nil {
			t.Error("expected unity record to be nil, but got a instace")
		}

		r.mu.Unlock()
		time.Sleep(500 * time.Microsecond)

		r.mu.Lock()
		u := r.records[1]
		r.mu.Unlock()

		if u == nil {
			t.Error("expected writeRecord() set a record unity on records")
			return
		}

		if u.MsgHash != msgHash {
			t.Errorf("expected message hash to be %q, but got %q", msgHash, u.MsgHash)
		}
	})
}

func TestDeleteRecord(t *testing.T) {
	msg := "hello world"
	msgHash := GenHash(msg)

	t.Run("should delete the record unity on record", func(t *testing.T) {
		_r := NewRecordMemory()
		r, ok := _r.(*recordMemory)
		if !ok {
			t.Fatal("NewRecordMemory did not returned a instace of record memory")
		}

		r.records[1] = &recordUnity{MsgHash: msgHash}

		u := r.records[1]
		if u == nil {
			t.Error("expected set a record unity on records")
			return
		}
		if u.MsgHash != msgHash {
			t.Errorf("expected message hash to be %q, but got %q", msgHash, u.MsgHash)
		}

		r.deleteRecord(1)
		if r.records[1] != nil {
			t.Error("expected the record unity to be nil after delete")
		}
	})

	t.Run("should timeout when mutex is locked", func(t *testing.T) {
		_r := NewRecordMemory()
		r, ok := _r.(*recordMemory)
		if !ok {
			t.Fatal("NewRecordMemory did not returned a instace of record memory")
		}

		r.records[1] = &recordUnity{MsgHash: msgHash}

		u := r.records[1]
		if u == nil {
			t.Error("expected set a record unity on records")
			return
		}
		if u.MsgHash != msgHash {
			t.Errorf("expected message hash to be %q, but got %q", msgHash, u.MsgHash)
		}

		r.mu.Lock()
		go func() {
			r.deleteRecord(1)
		}()

		time.Sleep(10 * time.Microsecond)
		if u := r.records[1]; u == nil {
			t.Error("expected get a record unity instace, but got nil")
		}

		r.mu.Unlock()
		time.Sleep(10 * time.Microsecond)

		r.mu.Lock()
		if r.records[1] != nil {
			t.Error("expected the record unity to be nil after delete")
		}
		r.mu.Unlock()
	})
}
