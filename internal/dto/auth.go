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

type AuthValidateEmailRequest struct {
	Code string `json:"code" validate:"required,uuid"`
}

type AuthSendValidationEmailRequest struct {
	Username string `json:"username" validate:"required"`
}
