package apimaker

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

type Pagination struct {
	Limit     int    `query:"limit"`
	Page      int    `query:"page"`
	Sort      string `query:"sort"`
	Unlimited string `query:"unlimited"`
}

func SetPagination(c echo.Context) (Pagination, error) {
	pag := new(Pagination)
	if err := c.Bind(pag); err != nil {
		return *pag, err
	}

	if pag.Limit < 1 || pag.Limit > 100 {
		pag.Limit = 10
	}

	if ok, _ := strconv.ParseBool(pag.Unlimited); ok {
		pag.Limit = -1
	}

	if pag.Page < 1 {
		pag.Page = 1
	}

	if pag.Sort == "" {
		pag.Sort = ""
	}

	return *pag, nil
}
