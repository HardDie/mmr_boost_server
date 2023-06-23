package dto

type PriceRequest struct {
	TypeID     int32 `json:"typeId" validate:"required,gt=0"`
	CurrentMmr int32 `json:"currentMmr" validate:"omitempty,gte=0,ltfield=TargetMmr"`
	TargetMmr  int32 `json:"targetMmr" validate:"omitempty,gtfield=CurrentMmr"`
}
