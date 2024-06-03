package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	apimaker "github.com/yasinsaee/api_maker"
	"github.com/yasinsaee/api_maker/sample/custom"
	"github.com/yasinsaee/api_maker/sample/product"
)

func main() {
	ec := echo.New()
	ec.Validator = &apimaker.CustomValidator{Validator: validator.New()}
	// Use Echo's default middleware
	ec.Use(middleware.Logger())
	ec.Use(middleware.Recover())

	proGP := ec.Group("/product")
	apiService := apimaker.NewAPIService("product", proGP, ec.Validator, ec.Logger)

	proGP.POST("/create", func(c echo.Context) error {
		pro := new(product.Product)
		form := new(product.AddProductForm)

		createServiceReq := apimaker.ServiceRequest{
			Context:    c,
			Model:      pro,
			Form:       form,
			AfterSave:  custom.MyCustomAfterSaveFunction,
			BeforeSave: custom.MyCustomBeforeSaveFunction,
		}

		return apiService.RequestService(apimaker.CreateServiceRequest, createServiceReq)
	})

	proGP.PUT("/edit/:id", func(c echo.Context) error {
		pro := new(product.Product)
		form := new(product.AddProductForm)

		if err := apiService.Edit(c, pro, form); err != nil {
			return err
		}
		return nil
	})

	proGP.GET("/list", func(c echo.Context) error {
		pro := new(product.Product)
		filter := new(product.ProductFilter)

		if err := apiService.List(c, pro, filter); err != nil {
			return err
		}
		return nil
	})

	proGP.GET("/view/:id", func(c echo.Context) error {
		pro := new(product.Product)

		if err := apiService.View(c, pro); err != nil {
			return err
		}
		return nil
	})

	proGP.DELETE("/delete/:id", func(c echo.Context) error {
		pro := new(product.Product)
		if err := apiService.Delete(c, pro); err != nil {
			return err
		}
		return nil
	})

	ec.Logger.Fatal(ec.Start(":1111"))
}
