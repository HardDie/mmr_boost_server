package entity

import "time"

type User struct {
	ID          int32      `json:"id"`
	Email       string     `json:"email"`
	Username    string     `json:"username"`
	RoleID      int32      `json:"roleId"`
	SteamID     string     `json:"steamId"`
	IsActivated bool       `json:"isActivated"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt"`
}
