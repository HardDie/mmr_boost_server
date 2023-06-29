package dto

type ApplicationListRequest struct {
	UserID   *int32
	StatusID *int32
}

type ApplicationItemRequest struct {
	UserID        *int32
	ApplicationID int32
}

type ApplicationUpdateStatusRequest struct {
	ApplicationID int32
	StatusID      int32
}

type ApplicationUpdateRequest struct {
	ApplicationID int32
	CurrentMMR    int32
	TargetMMR     int32
	TgContact     *string
	Price         *float64
}
