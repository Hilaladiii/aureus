package server

import (
	"github.com/Hilaladiii/aureus/internal/delivery/router"
	"github.com/Hilaladiii/aureus/pkg/exception"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func NewFiberServer(r *router.Router) *fiber.App {
	app := fiber.New(fiber.Config{
		JSONEncoder:  sonic.Marshal,
		JSONDecoder:  sonic.Unmarshal,
		ErrorHandler: exception.ErrorHandler,
	})

	app.Use(recover.New())

	r.Setup(app)

	return app
}
