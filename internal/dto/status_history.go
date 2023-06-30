package dto

type StatusHistoryNewEventRequest struct {
	UserID        int32 `json:"userId"`
	ApplicationID int32 `json:"applicationId"`
	NewStatusID   int32 `json:"newStatusId"`
}

type StatusHistoryListRequest struct {
	ApplicationID int32 `json:"applicationId" validation:"required,gt=0"`
}
