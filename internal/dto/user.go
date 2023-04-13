package dto

type UserUpdatePasswordRequest struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,nefield=OldPassword"`
}
