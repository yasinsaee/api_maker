package apimaker

import (
	"github.com/labstack/echo/v4"
)

type Model interface {
	Save() error
	GetOne(id interface{}) error
	//filter for filter, page filtering - totalCounts, totalPages, list , error
	List(filter Filter, pfilter Pagination) (int, int, interface{}, error)
	Remove(id interface{}) error
}

type Models []Model

type Params struct {
	Key   string
	Value interface{}
}

type CreateFunc struct {
	Function func(model Model, params ...Params) error
	Params   []Params
}

type Security struct {
	Authenticator func(c echo.Context) (bool, error)
	Authorizer    func(c echo.Context) (bool, error)
}
