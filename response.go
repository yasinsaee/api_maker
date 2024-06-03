package apimaker

import (
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
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
func ErrorResponse(c echo.Context, code int, err error, message string) error {
	resp := &Response{
		Code:         code,
		ErrorMessage: message,
	}
	if err != nil {
		resp.ErrorMessage = err.Error()
	}
	return c.JSON(code, resp)
}

func MongoErrorHandler(err error) (int, string) {
	switch err {
	case mongo.ErrNoDocuments:
		return 404, "not found"
	default:
		return 500, "internal server error"
	}
}

func ToMap(key string, value interface{}) (res []map[string]interface{}) {
	switch v := value.(type) {
	case string:
		res = append(res, map[string]interface{}{
			key: v,
		})
	case []string:
		for _, v := range v {
			res = append(res, map[string]interface{}{
				key: v,
			})
		}
	case []interface{}:
		for _, v := range v {
			res = append(res, map[string]interface{}{
				key: v,
			})
		}
	}
	return res
}
