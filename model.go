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

type Params struct {
	Key   string
	Value interface{}
}

type AfterSave struct {
	Function func(model Model, params ...Params) error
	Params   []Params
}

type BeforeSave struct {
	Function func(model Model, params ...Params) error
	Params   []Params
}

type (
	CheckLogin func(c echo.Context, model Model) error
)
