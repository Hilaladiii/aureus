package router

import (
	"github.com/Hilaladiii/aureus/internal/delivery/handler"
	"github.com/Hilaladiii/aureus/internal/delivery/middleware"

	"github.com/gofiber/fiber/v3"
)

type Router struct {
	UserHandler *handler.UserHandler
	Middleware  middleware.MiddlewareItf
}

func NewRouter(userHandler *handler.UserHandler, middleware *middleware.Middleware) *Router {
	return &Router{
		UserHandler: userHandler,
		Middleware:  middleware,
	}
}

func (r *Router) Setup(app *fiber.App) {
	api := app.Group("/api/v1")

	user := api.Group("/users")
	user.Post("/register", r.UserHandler.Register)
	user.Post("/login", r.UserHandler.Login)

	// private route
	user.Use(r.Middleware.JwtMiddleware())
	user.Get("/me", r.UserHandler.GetProfile)
}
