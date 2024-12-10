package recordhistory

import "errors"

var (
	ErrRecordIDCollision = errors.New("Record ID collision detected")
)
