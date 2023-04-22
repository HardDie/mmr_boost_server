package dto

type ApplicationCreateRequest struct {
	TypeID     int32  `json:"typeId" validate:"required,gt=0"`
	CurrentMMR int32  `json:"currentMmr" validate:"required,gt=0,lt=10000"`
	TargetMMR  int32  `json:"targetMmr" validate:"required,gtfield=CurrentMMR,lt=10000"`
	TgContact  string `json:"tgContact" validate:"required,gte=5"`
}

type ApplicationUserListRequest struct {
	UserID   int32  `json:"userId" validate:"required"`
	StatusID *int32 `json:"statusId" validate:"omitempty,gte=1,lte=7"`
}
