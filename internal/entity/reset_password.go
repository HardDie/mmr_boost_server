package entity

import "time"

type ResetPassword struct {
	ID        int32     `json:"id"`
	UserID    int32     `json:"userId"`
	Code      string    `json:"code"`
	ExpiredAt time.Time `json:"expiredAt"`
	CreatedAt time.Time `json:"createdAt"`
}
