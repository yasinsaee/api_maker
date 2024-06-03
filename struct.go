package apimaker

import (
	"encoding/json"
	"errors"

	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
)

// CustomValidator interface for custom validation logic
type CustomValidator interface {
	Validate() error
}

func Validation(form interface{}) error {
	if _, err := govalidator.ValidateStruct(form); err != nil {
		return err
	}

	if customValidator, ok := form.(CustomValidator); ok {
		if err := customValidator.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func BindStruct(g echo.Context, form interface{}, model interface{}) error {
	if err := g.Bind(form); err != nil {
		return errors.New("error in bind form")
	}

	if err := Validation(form); err != nil {
		return errors.New("error in validation form")
	}

	jsonString, err := json.Marshal(form)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonString, &model)
	if err != nil {
		return err
	}
	return nil
}
