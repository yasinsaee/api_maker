package apimaker

import "github.com/labstack/echo/v4"

// BaseServiceRequest defines common fields for all service requests.
type BaseServiceRequest struct {
	Context  echo.Context
	Model    Model
	Security Security
}

// CreateServiceRequest defines the structure for a service request used for creating resources.
type CreateServiceRequest struct {
	BaseServiceRequest
	Form       Form
	AfterSave  CreateFunc
	BeforeSave CreateFunc
}

// UpdateServiceRequest defines the structure for a service request used for updating resources.
type UpdateServiceRequest struct {
	BaseServiceRequest
	Form       Form
	AfterSave  CreateFunc
	BeforeSave CreateFunc
}

// ListServiceRequest defines the structure for a service request used for listing resources.
type ListServiceRequest struct {
	BaseServiceRequest
	Pagination    Pagination
	Filters       Filter
	BeforeGetList CreateFunc
	AfterGetList  CreateFunc
}

// ViewServiceRequest defines the structure for a service request used for viewing a single resource.
type ViewServiceRequest struct {
	BaseServiceRequest
	AfterFind CreateFunc
}

// DeleteServiceRequest defines the structure for a service request used for deleting a resource.
type DeleteServiceRequest struct {
	BaseServiceRequest
	BeforeRemove CreateFunc
	AfterRemove  CreateFunc
}
