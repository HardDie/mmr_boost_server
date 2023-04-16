package dto

type UserUpdatePasswordRequest struct {
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required,nefield=OldPassword"`
}

type UserUpdateSteamIDRequest struct {
	SteamID string `json:"steamId" validate:"required,numeric"`
}
