package middleware

import (
	"strings"

	"github.com/Hilaladiii/aureus/internal/model"
	"github.com/Hilaladiii/aureus/pkg/jwt"

	"github.com/gofiber/fiber/v3"
)

type MiddlewareItf interface {
	JwtMiddleware() fiber.Handler
	RoleMiddleware(role model.Role) fiber.Handler
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

		claims, err := m.jwt.VerifyToken(tokenString)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized)
		}
		c.Locals("user_id", claims.UserID)
		c.Locals("role", claims.Role)
		return c.Next()
	}
}

func (m *Middleware) RoleMiddleware(requiredRole model.Role) fiber.Handler {
	return func(c fiber.Ctx) error {
		role := c.Locals("role")
		userRole, ok := role.(model.Role)
		if !ok {
			return fiber.NewError(fiber.StatusUnauthorized)
		}

		if requiredRole != userRole {
			return fiber.NewError(fiber.StatusForbidden)
		}

		return c.Next()
	}
}
