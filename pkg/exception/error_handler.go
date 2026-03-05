package exception

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/Hilaladiii/aureus/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

func ErrorHandler(ctx fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	status := http.StatusText(code)
	var errMessages any = err.Error()

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		code = fiberErr.Code
		status = http.StatusText(code)
		errMessages = fiberErr.Message
	}

	var valErrs validator.ValidationErrors
	if errors.As(err, &valErrs) {
		code = fiber.StatusBadRequest
		status = "BAD_REQUEST"

		var errList []string
		for _, e := range valErrs {
			errMessage := fmt.Sprintf("Field '%s' is invalid (Reason: %s %s)", e.Field(), e.Tag(), e.Param())
			errList = append(errList, errMessage)
		}
		errMessages = errList
	}

	return ctx.Status(code).JSON(response.WebResponse{
		Code:   code,
		Status: status,
		Errors: errMessages,
	})
}
