// Package api defines structures and methods for handling API services with
// functionalities such as creating and editing resources, along with handling
// authentication, authorization, and validation.

package apimaker

import (
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
func (createService CreateServiceRequest) Create(a APIService) error {
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
func (updateService UpdateServiceRequest) Edit(a APIService) error {
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
// It performs the following steps:
// 1. Extract ID: Retrieves the ID of the model to be viewed from the context parameters.
// 2. Authentication: If an authenticator is provided, it checks if the request is authenticated.
// 3. Authorization: If an authorizer is provided, it checks if the request is authorized.
// 4. Retrieve Model: It retrieves the model from the database using its ID.
// 5. After Find Hook: It calls an optional after find function to perform any post-find operations.
// 6. Success Response: It returns a success response with the retrieved model.
//
// Parameters:
// - viewService: A ViewServiceRequest struct containing the context, model, security handlers, and after find hook.
// - a: An APIService interface providing methods for error and success responses.
//
// Returns:
// - error: An error if any step fails; otherwise, nil.
func (viewService ViewServiceRequest) View(a APIService) error {
	var (
		err error
	)

	// Step 1: Extract ID
	id := viewService.Context.Param("id")

	// Step 2: Authentication
	if viewService.Security.Authenticator != nil {
		if authenticated, err := viewService.Security.Authenticator(viewService.Context); err != nil || !authenticated {
			return a.ErrorResponse(viewService.Context, http.StatusUnauthorized, err, "authentication failed")
		}
	}

	// Step 3: Authorization
	if viewService.Security.Authorizer != nil {
		if authorized, err := viewService.Security.Authorizer(viewService.Context); err != nil || !authorized {
			return a.ErrorResponse(viewService.Context, http.StatusForbidden, err, "authorization failed")
		}
	}

	// Step 4: Retrieve Model
	if err = viewService.Model.GetOne(id); err != nil {
		return a.ErrorResponse(viewService.Context, http.StatusBadRequest, err, fmt.Sprintf("cannot find any %s", a.Name))
	}

	// Step 5: After Find Hook
	if viewService.AfterFind.Function != nil {
		if err = viewService.AfterFind.Function(viewService.Model, viewService.AfterFind.Params...); err != nil {
			return a.ErrorResponse(viewService.Context, http.StatusBadRequest, err, fmt.Sprintf("cannot use function after find, error : %s ", err.Error()))
		}
	}

	// Step 6: Success Response
	return SuccessResponse(viewService.Context, http.StatusOK, fmt.Sprintf("successfully loaded %s", a.Name), echo.Map{a.Name: viewService.Model}, MetaData{})
}

// List handles listing models with pagination and filtering.
func (listService ListServiceRequest) List(a APIService) error {
	var (
		err error
	)

	pfilter, _ := SetPagination(listService.Context)

	if listService.Security.Authenticator != nil {
		if authenticated, err := listService.Security.Authenticator(listService.Context); err != nil || !authenticated {
			return a.ErrorResponse(listService.Context, http.StatusUnauthorized, err, "authentication failed")
		}
	}

	// Step 3: Authorization
	if listService.Security.Authorizer != nil {
		if authorized, err := listService.Security.Authorizer(listService.Context); err != nil || !authorized {
			return a.ErrorResponse(listService.Context, http.StatusForbidden, err, "authorization failed")
		}
	}

	if err := listService.Context.Bind(listService.Filters); err != nil {
		return a.ErrorResponse(listService.Context, http.StatusBadRequest, err, fmt.Sprintf("cannot bind %s filter", a.Name))
	}

	if listService.BeforeGetList.Function != nil {
		if err = listService.BeforeGetList.Function(listService.Model, listService.BeforeGetList.Params...); err != nil {
			return a.ErrorResponse(listService.Context, http.StatusBadRequest, err, fmt.Sprintf("cannot use function before get list, error : %s ", err.Error()))
		}
	}

	totalCounts, totalPages, list, err := listService.Model.List(listService.Filters, pfilter)
	if err != nil {
		return a.ErrorResponse(listService.Context, http.StatusBadRequest, err, fmt.Sprintf("cannot find any %s", a.Name))
	}

	if listService.AfterGetList.Function != nil {
		if err = listService.AfterGetList.Function(listService.Model, listService.AfterGetList.Params...); err != nil {
			return a.ErrorResponse(listService.Context, http.StatusBadRequest, err, fmt.Sprintf("cannot use function after get list, error : %s ", err.Error()))
		}
	}

	return SuccessResponse(listService.Context, http.StatusOK, fmt.Sprintf("successfully loaded %s list", a.Name), echo.Map{
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
// It performs the following steps:
// 1. Extract ID: Retrieves the ID of the model to be deleted from the context parameters.
// 2. Authentication: If an authenticator is provided, it checks if the request is authenticated.
// 3. Authorization: If an authorizer is provided, it checks if the request is authorized.
// 4. Before Remove Hook: It calls an optional before remove function to perform any pre-remove operations.
// 5. Remove Model: It removes the model from the database.
// 6. After Remove Hook: It calls an optional after remove function to perform any post-remove operations.
// 7. Success Response: It returns a success response if the model is successfully removed.
//
// Parameters:
// - deleteService: A DeleteServiceRequest struct containing the context, model, security handlers, and hooks.
// - a: An APIService interface providing methods for error and success responses.
//
// Returns:
// - error: An error if any step fails; otherwise, nil.
func (deleteService DeleteServiceRequest) Delete(a APIService) error {
	var (
		err error
	)

	// Step 1: Extract ID
	id := deleteService.Context.Param("id")

	// Step 2: Authentication
	if deleteService.Security.Authenticator != nil {
		if authenticated, err := deleteService.Security.Authenticator(deleteService.Context); err != nil || !authenticated {
			return a.ErrorResponse(deleteService.Context, http.StatusUnauthorized, err, "authentication failed")
		}
	}

	// Step 3: Authorization
	if deleteService.Security.Authorizer != nil {
		if authorized, err := deleteService.Security.Authorizer(deleteService.Context); err != nil || !authorized {
			return a.ErrorResponse(deleteService.Context, http.StatusForbidden, err, "authorization failed")
		}
	}

	// Step 4: Before Remove Hook
	if deleteService.BeforeRemove.Function != nil {
		if err = deleteService.BeforeRemove.Function(deleteService.Model, deleteService.BeforeRemove.Params...); err != nil {
			return a.ErrorResponse(deleteService.Context, http.StatusBadRequest, err, fmt.Sprintf("cannot use function before remove, error : %s ", err.Error()))
		}
	}

	// Step 5: Remove Model
	if err = deleteService.Model.Remove(id); err != nil {
		return a.ErrorResponse(deleteService.Context, http.StatusBadRequest, err, fmt.Sprintf("cannot find any %s", a.Name))
	}

	// Step 6: After Remove Hook
	if deleteService.AfterRemove.Function != nil {
		if err = deleteService.AfterRemove.Function(deleteService.Model, deleteService.AfterRemove.Params...); err != nil {
			return a.ErrorResponse(deleteService.Context, http.StatusBadRequest, err, fmt.Sprintf("cannot use function after remove, error : %s ", err.Error()))
		}
	}

	// Step 7: Success Response
	return SuccessResponse(deleteService.Context, http.StatusOK, "successfully removed", nil, MetaData{})
}
