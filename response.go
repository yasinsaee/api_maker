package apimaker

import (
	"github.com/labstack/echo/v4"
)

type (
	Response struct {
		Code           int         `json:"code"`
		SuccessMessage string      `json:"success_message"`
		ErrorMessage   string      `json:"error_message"`
		Data           interface{} `json:"data"`
		MetaData       MetaData    `json:"metadata"`
	}

	MetaData struct {
		Limit       int    `json:"limit"`
		TotalCounts int    `json:"total_counts"`
		TotalPages  int    `json:"total_pages"`
		CurrentPage int    `json:"current_page"`
		NextPage    int    `json:"next_page"`
		Sort        string `json:"sort"`
	}
)

// SuccessResponse handles sending success responses.
func SuccessResponse(c echo.Context, code int, message string, data echo.Map, metaData MetaData) error {
	resp := &Response{
		Code:           code,
		SuccessMessage: message,
		Data:           data,
		MetaData:       metaData,
	}
	return c.JSON(code, resp)
}

// ErrorResponse handles sending error responses.
func (a *APIService) ErrorResponse(c echo.Context, code int, err error, message string) error {

	a.Logger.Errorf("%s: %v", message, err)

	resp := &Response{
		Code:         code,
		ErrorMessage: message,
	}

	if err != nil {
		resp.ErrorMessage = err.Error()
	}

	return c.JSON(code, resp)
}
