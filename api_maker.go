package apimaker

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type APIService struct {
	Name      string
	GroupName echo.Group
}

func NewAPIService(name string, group_name echo.Group) *APIService {
	return &APIService{
		Name:      name,
		GroupName: group_name,
	}
}

func (a APIService) Create(c echo.Context, model Model, form Form) error {
	var (
		err error
	)

	if err = BindStruct(c, form, model); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err, "failed to bind form")
	}

	if err = form.Bind(model); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err, "failed to bind form data")
	}

	if err := model.Save(); err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err, fmt.Sprintf("cannot add %s", a.Name))
	}

	return SuccessResponse(
		c,
		http.StatusOK,
		fmt.Sprintf("successfully added %s", a.Name),
		echo.Map{a.Name: model},
		MetaData{},
	)
}

func (a APIService) Edit(c echo.Context, model Model, form Form) error {
	var (
		err error
	)

	id := c.Param("id")

	if err = model.GetOne(id); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err, fmt.Sprintf("cannot find any %s", a.Name))
	}

	if err = BindStruct(c, form, model); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err, "failed to bind form")
	}

	if err = form.Bind(model); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err, "failed to bind form data")
	}

	if err = model.Save(); err != nil {
		return ErrorResponse(c, http.StatusInternalServerError, err, fmt.Sprintf("cannot edit %s", a.Name))
	}

	return SuccessResponse(c, http.StatusOK, fmt.Sprintf("successfully edited %s", a.Name), echo.Map{a.Name: model}, MetaData{})
}

// View handles retrieving a single model.
func (a APIService) View(c echo.Context, model Model) error {
	var (
		err error
	)

	id := c.Param("id")

	if err = model.GetOne(id); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err, fmt.Sprintf("cannot find any %s", a.Name))
	}

	return SuccessResponse(c, http.StatusOK, fmt.Sprintf("successfully loaded %s", a.Name), echo.Map{a.Name: model}, MetaData{})
}

// List handles listing models with pagination and filtering.
func (a APIService) List(c echo.Context, model Model, filter Filter) error {
	var (
		err error
	)

	pfilter, _ := SetPagination(c, false)

	if err := c.Bind(filter); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err, fmt.Sprintf("cannot bind %s filter", a.Name))
	}

	totalCounts, totalPages, list, err := model.List(filter, pfilter)
	if err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err, fmt.Sprintf("cannot find any %s", a.Name))
	}

	return SuccessResponse(c, http.StatusOK, fmt.Sprintf("successfully loaded %s list", a.Name), echo.Map{
		a.Name + "s":   list,
		"total_counts": totalCounts,
		"total_pages":  totalPages,
	}, MetaData{
		Limit:       pfilter.Limit,
		CurrentPage: pfilter.Page,
		TotalCounts: totalCounts,
		TotalPages:  totalPages,
		Sort:        pfilter.Sort,
	})
}

// Delete handles deleting a model.
func (a APIService) Delete(c echo.Context, model Model) error {
	var (
		err error
	)

	id := c.Param("id")

	if err = model.Remove(id); err != nil {
		return ErrorResponse(c, http.StatusBadRequest, err, fmt.Sprintf("cannot find any %s", a.Name))
	}

	return SuccessResponse(c, http.StatusOK, "successfully removed", nil, MetaData{})
}
