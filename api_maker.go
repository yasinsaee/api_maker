package apimaker

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Form interface {
	Bind(Model) error
}

type Model interface {
	Save() error
	GetOne(id interface{}) error
	//filter for filter, page filtering - totalCounts, totalPages, list , error
	List(filter Filter, pfilter Pagination) (int, int, interface{}, error)
	Remove(id interface{}) error
}

type Filter interface {
	GetFilters() map[string]interface{}
}

type APIService struct {
	Name      string
	GroupName string
}

func (a APIService) Create(c echo.Context, model Model, form Form) error {
	var (
		err error
	)

	resp := new(Response)

	if err = BindStruct(c, form, model); err != nil {
		resp.ErrorMessage = err.Error()
		resp.Code = http.StatusBadRequest
		return c.JSON(http.StatusBadRequest, resp)
	}

	if err = form.Bind(model); err != nil {
		resp.ErrorMessage = err.Error()
		resp.Code = http.StatusBadRequest
		return c.JSON(http.StatusBadRequest, resp)
	}

	if err := model.Save(); err != nil {
		resp.Code = http.StatusInternalServerError
		resp.ErrorMessage = fmt.Sprintf("Can not add %s", a.Name)
		return c.JSON(http.StatusInternalServerError, resp)
	}

	resp.SuccessMessage = fmt.Sprintf("Successfully added %s", a.Name)
	resp.Code = http.StatusOK
	resp.Data = echo.Map{
		a.Name: model,
	}
	return c.JSON(http.StatusOK, resp)
}

func (a APIService) View(c echo.Context, model Model) error {
	var (
		err error
	)

	resp := new(Response)

	id := c.Param("id")

	if err = model.GetOne(id); err != nil {
		resp.Code = http.StatusBadRequest
		resp.ErrorMessage = fmt.Sprintf("can not find any %s", a.Name)
		return c.JSON(http.StatusBadRequest, resp)
	}

	resp.SuccessMessage = fmt.Sprintf("Successfully load %s", a.Name)
	resp.Code = http.StatusOK
	resp.Data = echo.Map{
		a.Name: model,
	}
	return c.JSON(http.StatusOK, resp)
}

func (a APIService) List(c echo.Context, model Model, filter Filter) error {
	var (
		err                     error
		totalCounts, totalPages int
		list                    interface{}
	)

	resp := new(Response)

	pfilter, _ := SetPagination(c, false)

	if err = c.Bind(filter); err != nil {
		resp.Code = http.StatusBadRequest
		resp.ErrorMessage = fmt.Sprintf("can not bind %s filter", a.Name)
		return c.JSON(http.StatusBadRequest, resp)
	}

	if totalCounts, totalPages, list, err = model.List(filter, pfilter); err != nil {
		resp.Code = http.StatusBadRequest
		resp.ErrorMessage = fmt.Sprintf("can not find any %s", a.Name)
		return c.JSON(http.StatusBadRequest, resp)
	}
	resp.SuccessMessage = fmt.Sprintf("Successfully loaded %s list", a.Name)
	resp.Code = http.StatusOK
	resp.Data = echo.Map{
		a.Name + "s":   list,
		"total_counts": totalCounts,
		"total_pages":  totalPages,
	}
	resp.MetaData = MetaData{
		Limit:       pfilter.Limit,
		CurrentPage: pfilter.Page,
		TotalCounts: totalCounts,
		TotalPages:  totalPages,
		Sort:        pfilter.Sort,
	}

	return c.JSON(http.StatusOK, resp)
}

func (a APIService) Delete(c echo.Context, model Model) error {
	var (
		err error
	)

	resp := new(Response)

	id := c.Param("id")

	if err = model.Remove(id); err != nil {
		resp.Code = http.StatusBadRequest
		resp.ErrorMessage = fmt.Sprintf("can not find any %s", a.Name)
		return c.JSON(http.StatusBadRequest, resp)
	}

	resp.SuccessMessage = "Successfully remove"
	resp.Code = http.StatusOK
	return c.JSON(http.StatusOK, resp)
}
