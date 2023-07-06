package postgres

import "time"

func (r *User) SetTimeNow(now func() time.Time) {
	r.timeNow = now
}

func (r *Password) SetTimeNow(now func() time.Time) {
	r.timeNow = now
}
