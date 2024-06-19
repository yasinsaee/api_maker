package apimaker

import "github.com/labstack/echo/v4"

func CreateApi(apiService APIService, model Model, form Form) error {
	apiService.Group.POST("/create", func(c echo.Context) error {

		createService := CreateServiceRequest{
			BaseServiceRequest: BaseServiceRequest{
				Context: c,
				Model:   model,
			},
			Form: form,
		}.Create(apiService)

		return createService
	})

	return nil
}

func UpdateApi(apiService APIService, model Model, form Form) error {
	apiService.Group.PUT("/update/:id", func(c echo.Context) error {

		updateService := UpdateServiceRequest{
			BaseServiceRequest: BaseServiceRequest{
				Context: c,
				Model:   model,
			},
			Form: form,
		}.Edit(apiService)

		return updateService
	})

	return nil
}

func ListApi(apiService APIService, model Model, filter Filter) error {
	apiService.Group.GET("/list", func(c echo.Context) error {

		listService := ListServiceRequest{
			BaseServiceRequest: BaseServiceRequest{
				Context: c,
				Model:   model,
			},
			Filters: filter,
		}.List(apiService)

		return listService
	})

	return nil
}

func ViewApi(apiService APIService, model Model, filter Filter) error {
	apiService.Group.GET("/view/:id", func(c echo.Context) error {

		viewService := ViewServiceRequest{
			BaseServiceRequest: BaseServiceRequest{
				Context: c,
				Model:   model,
			},
		}.View(apiService)

		return viewService
	})

	return nil
}

func DeleteApi(apiService APIService, model Model, filter Filter) error {
	apiService.Group.DELETE("/delete/:id", func(c echo.Context) error {

		deleteService := DeleteServiceRequest{
			BaseServiceRequest: BaseServiceRequest{
				Context: c,
				Model:   model,
			},
		}.Delete(apiService)

		return deleteService
	})

	return nil
}
