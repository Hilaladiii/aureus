package router

import (
	"github.com/Hilaladiii/aureus/internal/delivery/handler"
	"github.com/Hilaladiii/aureus/internal/delivery/middleware"

	"github.com/gofiber/fiber/v3"
)

type Router struct {
	UserHandler     *handler.UserHandler
	CategoryHandler *handler.CategoryHandler
	Middleware      middleware.MiddlewareItf
}

func NewRouter(userHandler *handler.UserHandler, categoryHandler *handler.CategoryHandler, middleware *middleware.Middleware) *Router {
	return &Router{
		UserHandler:     userHandler,
		CategoryHandler: categoryHandler,
		Middleware:      middleware,
	}
}

func (r *Router) Setup(app *fiber.App) {
	api := app.Group("/api/v1")

	// users route
	user := api.Group("/users")
	user.Post("/register", r.UserHandler.Register)
	user.Post("/login", r.UserHandler.Login)
	user.Use(r.Middleware.JwtMiddleware())
	user.Get("/me", r.UserHandler.GetProfile)

	// categories route
	category := api.Group("/categories")
	category.Use(r.Middleware.JwtMiddleware())
	category.Post("", r.CategoryHandler.CreateCategory)
	category.Put("/:categoryId", r.CategoryHandler.UpdateCategory)
	category.Delete("/:categoryId", r.CategoryHandler.DeleteCategory)
	category.Get("", r.CategoryHandler.GetAll)
	category.Get("/:categoryId", r.CategoryHandler.GetByID)
}
