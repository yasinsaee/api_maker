package apimaker

import (
	"encoding/json"
	"errors"

	"github.com/labstack/echo/v4"
)

func BindStruct(g echo.Context, form interface{}, model interface{}) error {
	if err := g.Bind(form); err != nil {
		return errors.New("error in bind form")
	}

	if err := g.Validate(form); err != nil {
		return errors.New("error in validation form, error : " + err.Error())
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
