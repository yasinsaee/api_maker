package apimaker

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type APIService struct {
	Name      string
	Group     *echo.Group
	Validator echo.Validator
	Logger    echo.Logger
}

func NewAPIService(name string, group *echo.Group, validator echo.Validator, logger echo.Logger) *APIService {
	return &APIService{
		Name:      name,
		Group:     group,
		Validator: validator,
		Logger:    logger,
	}
}

type ServiceType struct {
	Name string
}

type ServiceRequest struct {
	Context    echo.Context
	Model      Model
	Form       Form
	AfterSave  AfterSave
	BeforeSave BeforeSave
	Security   Security
}

func (a APIService) RequestService(serviceType string, service ServiceRequest) error {
	if serviceType == CreateServiceRequest {
		return service.Create(a)
	} else if serviceType == UpdateServiceRequest {
		// return service.Create(a)
	}
	return errors.New("please select a service type")
}

func (createService ServiceRequest) Create(a APIService) error {
	var (
		err error
	)

	if createService.Security.Authenticator != nil {
        if authenticated, err := createService.Security.Authenticator(createService.Context); err != nil || !authenticated {
            return a.ErrorResponse(createService.Context, http.StatusUnauthorized, err, "authentication failed")
        }
    }

    if createService.Security.Authorizer != nil {
        if authorized, err := createService.Security.Authorizer(createService.Context); err != nil || !authorized {
            return a.ErrorResponse(createService.Context, http.StatusForbidden, err, "authorization failed")
        }
    }

	if err = BindStruct(createService.Context, createService.Form, createService.Model); err != nil {
		return a.ErrorResponse(createService.Context, http.StatusBadRequest, err, "failed to bind form")
	}

	if err = createService.Form.Bind(createService.Model); err != nil {
		return a.ErrorResponse(createService.Context, http.StatusBadRequest, err, "failed to bind form data")
	}

	if createService.BeforeSave.Function != nil {
		if err = createService.BeforeSave.Function(createService.Model, createService.BeforeSave.Params...); err != nil {
			return a.ErrorResponse(createService.Context, http.StatusBadRequest, err, fmt.Sprintf("cannot use function beforesave, error : %s ", err.Error()))
		}
	}

	if err = createService.Model.Save(); err != nil {
		return a.ErrorResponse(createService.Context, http.StatusInternalServerError, err, fmt.Sprintf("cannot add %s", a.Name))
	}

	if createService.AfterSave.Function != nil {
		if err = createService.AfterSave.Function(createService.Model, createService.AfterSave.Params...); err != nil {
			return a.ErrorResponse(createService.Context, http.StatusBadRequest, err, fmt.Sprintf("cannot use function aftersave, error : %s ", err.Error()))
		}
	}

	return SuccessResponse(
		createService.Context,
		http.StatusOK,
		fmt.Sprintf("successfully added %s", a.Name),
		echo.Map{a.Name: createService.Model},
		MetaData{},
	)
}

func (a APIService) Edit(c echo.Context, model Model, form Form) error {
	var (
		err error
	)

	id := c.Param("id")

	if err = model.GetOne(id); err != nil {
		return a.ErrorResponse(c, http.StatusBadRequest, err, fmt.Sprintf("cannot find any %s", a.Name))
	}

	if err = BindStruct(c, form, model); err != nil {
		return a.ErrorResponse(c, http.StatusBadRequest, err, "failed to bind form")
	}

	if err = form.Bind(model); err != nil {
		return a.ErrorResponse(c, http.StatusBadRequest, err, "failed to bind form data")
	}

	if err = model.Save(); err != nil {
		return a.ErrorResponse(c, http.StatusInternalServerError, err, fmt.Sprintf("cannot edit %s", a.Name))
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
		return a.ErrorResponse(c, http.StatusBadRequest, err, fmt.Sprintf("cannot find any %s", a.Name))
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
		return a.ErrorResponse(c, http.StatusBadRequest, err, fmt.Sprintf("cannot bind %s filter", a.Name))
	}

	totalCounts, totalPages, list, err := model.List(filter, pfilter)
	if err != nil {
		return a.ErrorResponse(c, http.StatusBadRequest, err, fmt.Sprintf("cannot find any %s", a.Name))
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
		return a.ErrorResponse(c, http.StatusBadRequest, err, fmt.Sprintf("cannot find any %s", a.Name))
	}

	return SuccessResponse(c, http.StatusOK, "successfully removed", nil, MetaData{})
}
