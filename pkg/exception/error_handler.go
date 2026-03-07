package exception

import (
	"fmt"
	"log"

	"github.com/Hilaladiii/aureus/pkg/response"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

func ErrorHandler(c fiber.Ctx, err error) error {
	if e, ok := err.(*AppError); ok {
		return c.Status(e.Code).JSON(response.ErrorResponse(e.Code, e.Message))
	}

	if e, ok := err.(*fiber.Error); ok {
		return c.Status(e.Code).JSON(response.ErrorResponse(e.Code, e.Message))
	}

	if e, ok := err.(validator.ValidationErrors); ok {
		var errMessages []string
		for _, fieldErr := range e {
			msg := fmt.Sprintf("Column '%s' failed validation rule: '%s'", fieldErr.Field(), fieldErr.Tag())
			errMessages = append(errMessages, msg)
		}
		if len(errMessages) > 1 {
			return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(fiber.StatusBadRequest, errMessages))
		}
		return c.Status(fiber.StatusBadRequest).JSON(response.ErrorResponse(fiber.StatusBadRequest, errMessages[0]))
	}

	log.Printf("ERROR 500 %s %s - %v\n", c.Method(), c.Path(), err)
	return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponse(fiber.StatusInternalServerError, "an error occurred in the system"))
}
