package entity

import "time"

type ApplicationPublic struct {
	ID           int32      `json:"id"`
	UserID       int32      `json:"userId"`
	StatusID     int32      `json:"statusId"`
	TypeID       int32      `json:"typeId"`
	CurrentMMR   int32      `json:"currentMmr"`
	TargetMMR    int32      `json:"targetMmr"`
	TgContact    string     `json:"tgContact"`
	IsPrivateSet bool       `json:"isPrivateSet"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	DeletedAt    *time.Time `json:"deletedAt"`
}

type ApplicationPrivate struct {
	ID            int32      `json:"id"`
	SteamLogin    *string    `json:"steamLogin"`
	SteamPassword *string    `json:"steamPassword"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	DeletedAt     *time.Time `json:"deletedAt"`
}
