package util

import (
	"errors"

	"github.com/gofiber/fiber/v3"
)

func GetJwtClaimLocals(c fiber.Ctx) (string, error) {
	userIDLocals := c.Locals("user_id")

	if userIDLocals == "" {
		return "", errors.New("userId not found")
	}

	userID, ok := userIDLocals.(string)
	if !ok {
		return "", errors.New("invalid userId")
	}

	return userID, nil
}
