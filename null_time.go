package xorm

import (
	"time"
)

type NullTime struct {
	Time  time.Time
	Valid bool
}

func (nt NullTime) Ptr() *time.Time {
	if !nt.Valid {
		return nil
	}
	return &nt.Time
}

func (nt NullTime) ValueOrZero() time.Time {
	if !nt.Valid {
		return time.Time{}
	}
	return nt.Time
}

func (nt NullTime) IsNil() bool {
	return !nt.Valid
}
