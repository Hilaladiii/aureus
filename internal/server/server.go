package server

import (
	"github.com/Hilaladiii/aureus/pkg/exception"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func NewFiberServer(r *Router) *fiber.App {
	app := fiber.New(fiber.Config{
		JSONEncoder:  sonic.Marshal,
		JSONDecoder:  sonic.Unmarshal,
		ErrorHandler: exception.ErrorHandler,
	})

	app.Use(recover.New())
	app.Use(cors.New())

	r.Setup(app)

	return app
}
