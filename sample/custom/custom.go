package custom

import (
	"fmt"

	"github.com/labstack/echo/v4"
	apimaker "github.com/yasinsaee/api_maker"
	"github.com/yasinsaee/api_maker/sample/product"
)

type User struct {
	Name string `json:"name"`
}

func (u User) MyCustomAfterSaveFunction(model apimaker.Model, params ...apimaker.Params) error {
	_, ok := model.(*product.Product)
	if !ok {
		return fmt.Errorf("invalid model type")
	}

	return nil
}

func MyCustomBeforeSaveFunction(model apimaker.Model, params ...apimaker.Params) error {
	_, ok := model.(*product.Product)
	if !ok {
		return fmt.Errorf("invalid model type")
	}

	proType := params[0].Value.(string)

	// Custom logic here
	if proType != "" {
		fmt.Println(proType)
	}

	return nil
}

func CheckLogin(c echo.Context, model apimaker.Model) error {
	_, ok := model.(*product.Product)
	if !ok {
		return fmt.Errorf("invalid model type")
	}

	// Custom logic here
	return nil
}
