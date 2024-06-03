package main

import (
	"github.com/labstack/echo/v4"
	apimaker "github.com/yasinsaee/api_maker"
	"github.com/yasinsaee/api_maker/sample/product"
)

func main() {
	ec := echo.New()

	proGP := ec.Group("/product")
	apiService := apimaker.NewAPIService("product", proGP, ec.Validator, ec.Logger)

	proGP.POST("/create", func(c echo.Context) error {
		pro := new(product.Product)
		form := new(product.AddProductForm)

		if err := apiService.Create(c, pro, form); err != nil {
			return err
		}
		return nil
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
