package apimaker

import (

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

func CreateSuccessResponse(code int, message string, data interface{}) Response {
	return Response{
		SuccessMessage: message,
		Code:           code,
		Data:           data,
	}

}

func CreateErrorResponse(code int, errorMessage string) Response {
	return Response{
		ErrorMessage: errorMessage,
		Code:         code,
	}
}

func NewMetaData(pagination Pagination, totalCounts, totalPages int) MetaData {
	return MetaData{
		Limit:       pagination.Limit,
		TotalCounts: totalCounts,
		TotalPages:  totalPages,
		CurrentPage: pagination.Page,
		NextPage:    pagination.Page + 1,
		Sort:        pagination.Sort,
	}
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
