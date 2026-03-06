package di

import (
	"github.com/Hilaladiii/aureus/internal/delivery/handler"
	"github.com/Hilaladiii/aureus/internal/repository"
	"github.com/Hilaladiii/aureus/internal/usecase"
	"github.com/google/wire"
)

var WalletSet = wire.NewSet(
	repository.NewWalletRepo,
	wire.Bind(new(repository.WalletRepoItf), new(*repository.WalletRepo)),
	usecase.NewWalletUsecase,
	wire.Bind(new(usecase.WalletUsecaseItf), new(*usecase.WalletUsecase)),
	handler.NewWalletHandler,
)
