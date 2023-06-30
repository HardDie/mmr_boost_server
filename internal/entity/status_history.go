package entity

import "time"

type StatusHistory struct {
	ID            int32     `json:"id"`
	UserID        int32     `json:"userId"`
	ApplicationID int32     `json:"applicationId"`
	NewStatusID   int32     `json:"newStatusId"`
	CreatedAt     time.Time `json:"createdAt"`
}
