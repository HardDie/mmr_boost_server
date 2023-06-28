package dto

type ApplicationCreateRequest struct {
	UserID     int32  `json:"userId" validate:"required,gt=0"`
	TypeID     int32  `json:"typeId" validate:"required,gt=0"`
	CurrentMMR int32  `json:"currentMmr" validate:"required,gte=0,ltfield=TargetMMR"`
	TargetMMR  int32  `json:"targetMmr" validate:"required,gtfield=CurrentMMR,lte=7000"`
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

type ApplicationManagementListRequest struct {
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

type ApplicationManagementItemRequest struct {
	ApplicationID int32 `json:"applicationId" validate:"required,gt=0"`
}

type ApplicationItemDeleteRequest struct {
	ApplicationID int32 `json:"applicationId" validate:"required,gt=0"`

	// autofill
	UserID int32
}

type ApplicationUpdateStatusRequest struct {
	ApplicationID int32
	StatusID      int32
}
