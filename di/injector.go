//go:build wireinject
// +build wireinject

package di

import (
	"github.com/Hilaladiii/aureus/internal/delivery/middleware"
	"github.com/Hilaladiii/aureus/internal/delivery/router"
	"github.com/Hilaladiii/aureus/internal/delivery/server"
	"github.com/Hilaladiii/aureus/pkg/config"
	"github.com/Hilaladiii/aureus/pkg/driver/db"
	"github.com/Hilaladiii/aureus/pkg/jwt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/wire"
)

func ProvideValidator() *validator.Validate {
	return validator.New()
}

func InitializeApp() (*fiber.App, error) {
	wire.Build(
		config.LoadEnv,
		db.NewDB,
		jwt.NewJwt,
		ProvideValidator,
		UserSet,
		middleware.NewMiddleware,
		router.NewRouter,
		server.NewFiberServer,
	)

	return nil, nil
}
