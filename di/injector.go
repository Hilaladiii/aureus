//go:build wireinject
// +build wireinject

package di

import (
	"github.com/Hilaladiii/aureus/internal/delivery/middleware"
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
		WalletSet,
		middleware.NewMiddleware,
		config.NewRouter,
		config.NewFiberServer,
	)

	return nil, nil
}
