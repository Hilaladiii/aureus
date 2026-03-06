//go:build wireinject
// +build wireinject

package di

import (
	"github.com/Hilaladiii/aureus/internal/delivery/middleware"
	"github.com/Hilaladiii/aureus/internal/delivery/router"
	"github.com/Hilaladiii/aureus/internal/delivery/server"
	"github.com/Hilaladiii/aureus/pkg/config"
	"github.com/Hilaladiii/aureus/pkg/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/google/wire"
)

func InitializeApp() (*fiber.App, error) {
	wire.Build(
		config.LoadEnv,
		config.NewDB,
		jwt.NewJwt,
		config.NewValidator,
		UserSet,
		CategorySet,
		middleware.NewMiddleware,
		router.NewRouter,
		server.NewFiberServer,
	)

	return nil, nil
}
