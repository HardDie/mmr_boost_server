package postgres

import "time"

func (r *User) SetTimeNow(now func() time.Time) {
	r.timeNow = now
}

func (r *Password) SetTimeNow(now func() time.Time) {
	r.timeNow = now
}

func (r *AccessToken) SetTimeNow(now func() time.Time) {
	r.timeNow = now
}

func (r *Application) SetTimeNow(now func() time.Time) {
	r.timeNow = now
}
