package dto

type ApplicationCreateRequest struct {
	UserID     int32  `json:"userId" validate:"required,gt=0"`
	TypeID     int32  `json:"typeId" validate:"required,gt=0"`
	CurrentMMR int32  `json:"currentMmr" validate:"required,gt=0,lt=10000"`
	TargetMMR  int32  `json:"targetMmr" validate:"required,gtfield=CurrentMMR,lt=10000"`
	TgContact  string `json:"tgContact" validate:"required,gte=5"`
}

type ApplicationListRequest struct {
	UserID   *int32
	StatusID *int32
}

type ApplicationUserListRequest struct {
	UserID   int32  `json:"userId" validate:"required,gt=0"`
	StatusID *int32 `json:"statusId" validate:"omitempty,gte=1,lte=7"`
}

type ApplicationManagementUserListRequest struct {
	UserID   *int32 `json:"userId" validate:"omitempty,gt=0"`
	StatusID *int32 `json:"statusId" validate:"omitempty,gte=1,lte=7"`
}

type ApplicationItemRequest struct {
	UserID        *int32
	ApplicationID int32
}

type ApplicationUserItemRequest struct {
	UserID        int32 `json:"userId" validate:"required,gt=0"`
	ApplicationID int32 `json:"applicationId" validate:"required,gt=0"`
}

type ApplicationManagementUserItemRequest struct {
	ApplicationID int32 `json:"applicationId" validate:"required,gt=0"`
}
