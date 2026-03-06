package response

import "github.com/gofiber/utils/v2"

type WebResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   any    `json:"data,omitempty"`
	Errors any    `json:"errors,omitempty"`
}

func SuccessResponse(code int, data any) WebResponse {
	return WebResponse{
		Code:   code,
		Status: utils.StatusMessage(code),
		Data:   data,
	}
}

func ErrorResponse(code int, err any) WebResponse {
	return WebResponse{
		Code:   code,
		Status: utils.StatusMessage(code),
		Errors: err,
	}
}
