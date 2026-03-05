package middleware

import (
	"strings"

	"github.com/Hilaladiii/aureus/pkg/jwt"

	"github.com/gofiber/fiber/v3"
)

type MiddlewareItf interface {
	JwtMiddleware() fiber.Handler
}

type Middleware struct {
	jwt jwt.JwtItf
}

func NewMiddleware(jwt jwt.JwtItf) *Middleware {
	return &Middleware{jwt: jwt}
}

func (m *Middleware) JwtMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return fiber.NewError(fiber.StatusUnauthorized)
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		userID, err := m.jwt.VerifyToken(tokenString)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized)
		}
		c.Locals("user_id", userID)
		return c.Next()
	}
}
