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

	u := new(custom.User)

	proGP.POST("/create", func(c echo.Context) error {
		pro := new(product.Product)
		form := new(product.AddProductForm)

		proType := c.QueryParam("type")

		afterSave := apimaker.CreateFunc{
			Function: u.MyCustomAfterSaveFunction,
			Params: []apimaker.Params{
				{
					Key:   "username",
					Value: true,
				},
			},
		}

		beforeSave := apimaker.CreateFunc{
			Function: custom.MyCustomBeforeSaveFunction,
			Params: []apimaker.Params{
				{
					Key:   "proType",
					Value: proType,
				},
			},
		}

		createServiceReq := apimaker.CreateServiceRequest{
			BaseServiceRequest: apimaker.BaseServiceRequest{
				Context:  c,
				Model:    pro,
				Security: apimaker.Security{},
			},
			Form:       form,
			AfterSave:  afterSave,
			BeforeSave: beforeSave,
		}

		return createServiceReq.Create(*apiService)
	})

	proGP.PUT("/edit/:id", func(c echo.Context) error {
		pro := new(product.Product)
		form := new(product.AddProductForm)

		updateServiceReq := apimaker.UpdateServiceRequest{
			BaseServiceRequest: apimaker.BaseServiceRequest{
				Context:  c,
				Model:    pro,
				Security: apimaker.Security{},
			},
			Form: form,
		}

		if err := updateServiceReq.Edit(*apiService); err != nil {
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

		viewServiceReq := apimaker.ViewServiceRequest{
			BaseServiceRequest: apimaker.BaseServiceRequest{
				Context:  c,
				Model:    pro,
				Security: apimaker.Security{},
			},
		}.View(*apiService)

		if err := viewServiceReq; err != nil {
			return err
		}

		return nil
	})

	proGP.DELETE("/delete/:id", func(c echo.Context) error {
		pro := new(product.Product)

		deleteServiceReq := apimaker.DeleteServiceRequest{
			BaseServiceRequest: apimaker.BaseServiceRequest{
				Context:  c,
				Model:    pro,
				Security: apimaker.Security{}, //optional
			},
			BeforeRemove: apimaker.CreateFunc{}, //optional
			AfterRemove:  apimaker.CreateFunc{}, // optional
		}.Delete(*apiService)

		if err := deleteServiceReq; err != nil {
			return err
		}
		return nil
	})

	ec.Logger.Fatal(ec.Start(":1111"))
}
