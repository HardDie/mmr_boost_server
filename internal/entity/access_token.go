package entity

import "time"

type AccessToken struct {
	ID        int32      `json:"id"`
	UserID    int32      `json:"userId"`
	TokenHash string     `json:"tokenHash"`
	ExpiredAt time.Time  `json:"expiredAt"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}
