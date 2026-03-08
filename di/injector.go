//go:build wireinject
// +build wireinject

package di

import (
	"github.com/Hilaladiii/aureus/internal/delivery/middleware"
	"github.com/Hilaladiii/aureus/internal/server"
	"github.com/Hilaladiii/aureus/internal/worker"
	"github.com/Hilaladiii/aureus/pkg/config"
	"github.com/Hilaladiii/aureus/pkg/jwt"
	"github.com/google/wire"
)

func InitializeApp() (*server.App, error) {
	wire.Build(
		config.LoadEnv,
		config.NewDB,
		jwt.NewJwt,
		config.NewValidator,
		config.NewSeaweedFSStorage,
		config.NewTxManager,
		config.NewRedisClient,
		wire.Bind(new(config.SeaweedFSStorageItf), new(*config.SeaweedFSStorage)),
		UserSet,
		CategorySet,
		WalletSet,
		AuctionSet,
		BidSet,
		worker.NewAuctionWorker,
		middleware.NewMiddleware,
		server.NewRouter,
		server.NewFiberServer,
		server.NewApp,
	)

	return nil, nil
}
