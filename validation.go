package apimaker

import "github.com/go-playground/validator/v10"

type CustomValidator struct {
	Validator *validator.Validate
}

// Validate performs validation using the underlying validator instance
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		return err
	}
	return nil
}
