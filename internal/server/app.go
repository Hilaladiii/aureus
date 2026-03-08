package server

import (
	"github.com/Hilaladiii/aureus/internal/worker"
	"github.com/gofiber/fiber/v3"
)

type App struct {
	Web    *fiber.App
	Worker *worker.AuctionWorker
}

func NewApp(web *fiber.App, worker *worker.AuctionWorker) *App {
	return &App{
		Web:    web,
		Worker: worker,
	}
}
