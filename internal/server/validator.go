package server

import "github.com/go-playground/validator/v10"

var (
	v *validator.Validate
)

func getValidator() *validator.Validate {
	if v == nil {
		v = validator.New()
	}
	return v
}
