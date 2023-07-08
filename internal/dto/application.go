package dto

type ApplicationCreateRequest struct {
	UserID     int32  `json:"userId" validate:"required,gt=0"`
	TypeID     int32  `json:"typeId" validate:"required,gt=0"`
	CurrentMMR int32  `json:"currentMmr" validate:"omitempty,gte=0,ltfield=TargetMMR"`
	TargetMMR  int32  `json:"targetMmr" validate:"omitempty,lte=7000,gtfield=CurrentMMR"`
	TgContact  string `json:"tgContact" validate:"required,gte=5"`
	Price      int32  `json:"price"`
}

type ApplicationUserListRequest struct {
	UserID   int32  `json:"userId" validate:"required,gt=0"`
	StatusID *int32 `json:"statusId" validate:"omitempty,gte=1,lte=8"`
}

type ApplicationManagementListRequest struct {
	UserID   *int32 `json:"userId" validate:"omitempty,gt=0"`
	StatusID *int32 `json:"statusId" validate:"omitempty,gte=1,lte=8"`
}

type ApplicationUserItemRequest struct {
	UserID        int32 `json:"userId" validate:"required,gt=0"`
	ApplicationID int32 `json:"applicationId" validate:"required,gt=0"`
}

type ApplicationManagementItemRequest struct {
	ApplicationID int32 `json:"applicationId" validate:"required,gt=0"`
}

type ApplicationManagementPrivateItemRequest struct {
	ApplicationID int32 `json:"applicationId" validate:"required,gt=0"`

	// autofill
	UserID int32
}

type ApplicationItemDeleteRequest struct {
	ApplicationID int32 `json:"applicationId" validate:"required,gt=0"`

	// autofill
	UserID int32
}

type ApplicationManagementUpdateStatusRequest struct {
	ApplicationID int32 `json:"applicationId" validate:"required,gt=0"`
	StatusID      int32 `json:"statusId" validate:"required,gte=1,lte=8"`

	// autofill
	UserID int32
}

type ApplicationManagementUpdateItemRequest struct {
	ApplicationID int32 `json:"applicationId" validate:"required,gt=0"`
	CurrentMMR    int32 `json:"currentMmr" validate:"omitempty,gte=0,ltfield=TargetMMR"`
	TargetMMR     int32 `json:"targetMmr" validate:"omitempty,lte=7000,gtfield=CurrentMMR"`
	Price         int32 `json:"price" validate:"required,gt=0"`

	// autofill
	UserID int32
}

type ApplicationManagementUpdatePrivateRequest struct {
	ApplicationID int32  `json:"applicationId" validate:"required,gt=0"`
	SteamLogin    string `json:"steamLogin" validate:"required"`
	SteamPassword string `json:"steamPassword" validate:"required"`

	// autofill
	UserID int32
}
