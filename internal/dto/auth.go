package dto

type AuthRegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
