package errs

import (
	"fmt"
	"net/http"
)

var (
	ErrInternalError  = NewError("internal error")
	ErrBadRequest     = NewError("bad request", http.StatusBadRequest)
	ErrUserBlocked    = NewError("user is blocked", http.StatusUnauthorized)
	ErrSessionInvalid = NewError("session invalid", http.StatusUnauthorized)

	ErrEmailValidationCodeExpired  = NewError("validation code expired", http.StatusBadRequest)
	ErrEmailValidationCodeNotExist = NewError("validation code not exist", http.StatusBadRequest)
)

type MmrError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Err     error  `json:"err"`
}

func NewError(message string, code ...int) *MmrError {
	err := &MmrError{
		Message: message,
		Code:    http.StatusInternalServerError,
	}
	if len(code) > 0 {
		err.Code = code[0]
	}
	return err
}

func (e MmrError) Error() string {
	return fmt.Sprintf("HTTP[%d] %s", e.GetCode(), e.GetMessage())
}
func (e MmrError) Unwrap() error {
	return e.Err
}

func (e *MmrError) HTTP(code int) *MmrError {
	return &MmrError{
		Message: e.Message,
		Code:    code,
		Err:     e,
	}
}
func (e *MmrError) AddMessage(message string) *MmrError {
	return &MmrError{
		Message: message,
		Code:    e.Code,
		Err:     e,
	}
}

func (e *MmrError) GetCode() int       { return e.Code }
func (e *MmrError) GetMessage() string { return e.Message }
