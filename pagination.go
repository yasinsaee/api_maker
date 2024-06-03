package apimaker

import "github.com/labstack/echo/v4"

type Pagination struct {
	Limit     int    `query:"limit"`
	Page      int    `query:"page"`
	Sort      string `query:"sort"`
	Unlimited bool
}

func SetPagination(c echo.Context, unlimited bool) (Pagination, error) {
	pag := new(Pagination)
	if err := c.Bind(pag); err != nil {
		return *pag, err
	}

	if pag.Limit < 1 || pag.Limit > 100 {
		pag.Limit = 10
	}

	if unlimited {
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
