package memory

import (
	"testing"
)

func TestNewRecordHistory(t *testing.T) {
	r := NewRecordHistory()
	if r == nil {
		t.Errorf("NewRecordHistory() failed")
	}
	if _, ok := r.(*recordHistory); !ok {
		t.Errorf("NewRecordHistory() failed")
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
	t.Run("Simple Add", func(t *testing.T) {
		r := NewRecordHistory()
		err := r.AddRecord(1, "test")
		if err != nil {
			t.Errorf("AddRecord() failed")
		}
	})

	t.Run("Collision Check", func(t *testing.T) {
		r := NewRecordHistory()
		err := r.AddRecord(1, "test")
		if err != nil {
			t.Errorf("AddRecord() failed")
		}

		err = r.AddRecord(1, "test")
		if err != ErrRecordIDCollision {
			t.Errorf("AddRecord() failed; Expected collision error")
		}
	})
}

func TestRemoveRecord(t *testing.T) {
	t.Run("Simple Remove", func(t *testing.T) {
		id := uint64(1)
		r := NewRecordHistory()
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
	t.Run("GetRecord", func(t *testing.T) {
		r := NewRecordHistory()
		r.AddRecord(1, "")
		if r.GetRecord(1) == nil {
			t.Errorf("GetRecord() failed")
		}
	})
}

func TestReadRecord(t *testing.T) {
	t.Run("ReadRecord", func(t *testing.T) {
		_r := NewRecordHistory()
		r := _r.(*recordHistory)
		r.AddRecord(1, "")
		if r.readRecord(1) == nil {
			t.Errorf("readRecord() failed")
		}
	})
}

func TestWriteRecord(t *testing.T) {
	t.Run("WriteRecord", func(t *testing.T) {
		_r := NewRecordHistory()
		r := _r.(*recordHistory)
		r.writeRecord(1, &recordUnity{})
		if r.records[1] == nil {
			t.Errorf("writeRecord() failed")
		}
	})
}

func TestDeleteRecord(t *testing.T) {
	t.Run("DeleteRecord", func(t *testing.T) {
		_r := NewRecordHistory()
		r := _r.(*recordHistory)
		r.AddRecord(1, "")
		r.deleteRecord(1)
		if r.records[1] != nil {
			t.Errorf("deleteRecord() failed")
		}
	})
}
