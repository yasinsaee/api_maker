// Package api defines structures and methods for handling API services with
// functionalities such as creating and editing resources, along with handling
// authentication, authorization, and validation.

package apimaker

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// APIService defines the structure for an API service, containing necessary
// components such as name, group, validator, and logger.
type APIService struct {
	Name      string
	Group     *echo.Group
	Validator echo.Validator
	Logger    echo.Logger
}

// NewAPIService creates a new instance of APIService with the given parameters.
//
// Parameters:
// - name: The name of the API service.
// - group: The echo group associated with the API service.
// - validator: The validator to be used for request validation.
// - logger: The logger to be used for logging.
//
// Returns:
// - *APIService: A pointer to the newly created APIService instance.
func NewAPIService(name string, group *echo.Group, validator echo.Validator, logger echo.Logger) *APIService {
	return &APIService{
		Name:      name,
		Group:     group,
		Validator: validator,
		Logger:    logger,
	}
}

// ServiceType defines the structure for a service type, containing its name.
type ServiceType struct {
	Name string
}

// ServiceRequest defines the structure for a service request, containing
// context, model, form, after save hook, before save hook, and security.
type ServiceRequest struct {
	Context    echo.Context
	Model      Model
	Form       Form
	AfterSave  AfterSave
	BeforeSave BeforeSave
	Security   Security
}

// RequestService handles a service request based on the provided service type.
//
// Parameters:
// - serviceType: The type of service request (e.g., CreateServiceRequest, UpdateServiceRequest).
// - service: The service request containing context, model, form, hooks, and security.
//
// Returns:
// - error: An error if the service type is not recognized or if the service request fails.
func (a APIService) RequestService(serviceType string, service ServiceRequest) error {
	if serviceType == CreateServiceRequest {
		return service.Create(a)
	} else if serviceType == UpdateServiceRequest {
		return service.Edit(a)
	}
	return errors.New("please select a service type")
}

// Create handles the creation of a new resource in the API service.
// It performs the following steps:
// 1. Authentication: If an authenticator is provided, it checks if the request is authenticated.
// 2. Authorization: If an authorizer is provided, it checks if the request is authorized.
// 3. Data Binding: It binds the request data to the provided model.
// 4. Before Save Hook: It calls an optional before save function to perform any pre-save operations.
// 5. Save: It saves the model to the database.
// 6. After Save Hook: It calls an optional after save function to perform any post-save operations.
// 7. Success Response: It returns a success response if all steps are completed without errors.
//
// Parameters:
// - createService: A ServiceRequest struct containing the context, security handlers, form, model, and hooks.
// - a: An APIService interface providing methods for error and success responses.
//
// Returns:
// - error: An error if any step fails; otherwise, nil.
func (createService ServiceRequest) Create(a APIService) error {
	var (
		err error
	)

	// Step 1: Authentication
	if createService.Security.Authenticator != nil {
		if authenticated, err := createService.Security.Authenticator(createService.Context); err != nil || !authenticated {
			return a.ErrorResponse(createService.Context, http.StatusUnauthorized, err, "authentication failed")
		}
	}

	// Step 2: Authorization
	if createService.Security.Authorizer != nil {
		if authorized, err := createService.Security.Authorizer(createService.Context); err != nil || !authorized {
			return a.ErrorResponse(createService.Context, http.StatusForbidden, err, "authorization failed")
		}
	}

	// Step 3: Data Binding
	if err = BindStruct(createService.Context, createService.Form, createService.Model); err != nil {
		return a.ErrorResponse(createService.Context, http.StatusBadRequest, err, "failed to bind form")
	}

	if err = createService.Form.Bind(createService.Model); err != nil {
		return a.ErrorResponse(createService.Context, http.StatusBadRequest, err, "failed to bind form data")
	}

	// Step 4: Before Save Hook
	if createService.BeforeSave.Function != nil {
		if err = createService.BeforeSave.Function(createService.Model, createService.BeforeSave.Params...); err != nil {
			return a.ErrorResponse(createService.Context, http.StatusBadRequest, err, fmt.Sprintf("cannot use function beforesave, error : %s ", err.Error()))
		}
	}

	// Step 5: Save the Model
	if err = createService.Model.Save(); err != nil {
		return a.ErrorResponse(createService.Context, http.StatusInternalServerError, err, fmt.Sprintf("cannot add %s", a.Name))
	}

	// Step 6: After Save Hook
	if createService.AfterSave.Function != nil {
		if err = createService.AfterSave.Function(createService.Model, createService.AfterSave.Params...); err != nil {
			return a.ErrorResponse(createService.Context, http.StatusBadRequest, err, fmt.Sprintf("cannot use function aftersave, error : %s ", err.Error()))
		}
	}

	// Step 7: Success Response
	return SuccessResponse(
		createService.Context,
		http.StatusOK,
		fmt.Sprintf("successfully added %s", a.Name),
		echo.Map{a.Name: createService.Model},
		MetaData{},
	)
}

// Edit handles the editing of an existing resource in the API service.
// It performs the following steps:
// 1. Extract ID: Retrieves the ID of the resource to be edited from the context parameters.
// 2. Authentication: If an authenticator is provided, it checks if the request is authenticated.
// 3. Authorization: If an authorizer is provided, it checks if the request is authorized.
// 4. Fetch Resource: Retrieves the existing resource by its ID.
// 5. Data Binding: It binds the request data to the fetched model.
// 6. Before Save Hook: It calls an optional before save function to perform any pre-save operations.
// 7. Save: It updates the model in the database.
// 8. After Save Hook: It calls an optional after save function to perform any post-save operations.
// 9. Success Response: It returns a success response if all steps are completed without errors.
//
// Parameters:
// - updateService: A ServiceRequest struct containing the context, security handlers, form, model, and hooks.
// - a: An APIService interface providing methods for error and success responses.
//
// Returns:
// - error: An error if any step fails; otherwise, nil.
func (updateService ServiceRequest) Edit(a APIService) error {
	var (
		err error
	)

	// Step 1: Extract ID
	id := updateService.Context.Param("id")

	// Step 2: Authentication
	if updateService.Security.Authenticator != nil {
		if authenticated, err := updateService.Security.Authenticator(updateService.Context); err != nil || !authenticated {
			return a.ErrorResponse(updateService.Context, http.StatusUnauthorized, err, "authentication failed")
		}
	}

	// Step 3: Authorization
	if updateService.Security.Authorizer != nil {
		if authorized, err := updateService.Security.Authorizer(updateService.Context); err != nil || !authorized {
			return a.ErrorResponse(updateService.Context, http.StatusForbidden, err, "authorization failed")
		}
	}

	// Step 4: Fetch Resource
	if err = updateService.Model.GetOne(id); err != nil {
		return a.ErrorResponse(updateService.Context, http.StatusBadRequest, err, fmt.Sprintf("cannot find any %s", a.Name))
	}

	// Step 5: Data Binding
	if err = BindStruct(updateService.Context, updateService.Form, updateService.Model); err != nil {
		return a.ErrorResponse(updateService.Context, http.StatusBadRequest, err, "failed to bind form")
	}

	if err = updateService.Form.Bind(updateService.Model); err != nil {
		return a.ErrorResponse(updateService.Context, http.StatusBadRequest, err, "failed to bind form data")
	}

	// Step 6: Before Save Hook
	if updateService.BeforeSave.Function != nil {
		if err = updateService.BeforeSave.Function(updateService.Model, updateService.BeforeSave.Params...); err != nil {
			return a.ErrorResponse(updateService.Context, http.StatusBadRequest, err, fmt.Sprintf("cannot use function beforesave, error : %s ", err.Error()))
		}
	}

	// Step 7: Save the Model
	if err = updateService.Model.Save(); err != nil {
		return a.ErrorResponse(updateService.Context, http.StatusInternalServerError, err, fmt.Sprintf("cannot edit %s", a.Name))
	}

	// Step 8: After Save Hook
	if updateService.AfterSave.Function != nil {
		if err = updateService.AfterSave.Function(updateService.Model, updateService.AfterSave.Params...); err != nil {
			return a.ErrorResponse(updateService.Context, http.StatusBadRequest, err, fmt.Sprintf("cannot use function aftersave, error : %s ", err.Error()))
		}
	}

	// Step 9: Success Response
	return SuccessResponse(updateService.Context, http.StatusOK, fmt.Sprintf("successfully edited %s", a.Name), echo.Map{a.Name: updateService.Model}, MetaData{})
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
