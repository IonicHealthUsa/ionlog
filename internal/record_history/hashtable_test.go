package recordhistory

import "testing"

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
	t.Run("TestAddRecord", func(t *testing.T) {
		r := NewRecordHistory()
		err := r.AddRecord(1, "test", logOnce)
		if err != nil {
			t.Errorf("AddRecord() failed")
		}
	})

	t.Run("TestAddRecord collision check", func(t *testing.T) {
		r := NewRecordHistory()
		err := r.AddRecord(1, "test", logOnce)
		if err != nil {
			t.Errorf("AddRecord() failed")
		}

		err = r.AddRecord(1, "test", logOnce)
		if err != ErrRecordIDCollision {
			t.Errorf("AddRecord() failed; Expected collision error")
		}
	})
}

func TestRemoveRecord(t *testing.T) {
	id := uint64(1)
	r := NewRecordHistory()
	r.AddRecord(id, "test", logOnce)

	if r.GetRecord(id) == nil {
		t.Errorf("Test preset failed")
	}

	r.RemoveRecord(id)

	if r.GetRecord(id) != nil {
		t.Errorf("RemoveRecord() failed")
	}
}

func TestGetRecord(t *testing.T) {
	t.Run("GetRecord", func(t *testing.T) {
		r := NewRecordHistory()
		r.AddRecord(1, "", logOnce)
		if r.GetRecord(1) == nil {
			t.Errorf("GetRecord() failed")
		}
	})

	t.Run("GetRecord instance check (pointer)", func(t *testing.T) {
		r := NewRecordHistory()
		r.AddRecord(1, "", logOnce)
		rec := r.GetRecord(1)
		rec.Mode = logOnChange
		if r.GetRecord(1).Mode != logOnChange {
			t.Errorf("GetRecord() failed; Expected pointer instance")
		}
	})
}
