package di

import (
	"github.com/Hilaladiii/aureus/internal/delivery/handler"
	"github.com/Hilaladiii/aureus/internal/repository"
	"github.com/Hilaladiii/aureus/internal/usecase"

	"github.com/google/wire"
)

var UserSet = wire.NewSet(
	repository.NewUserRepo,
	wire.Bind(new(repository.UserRepoItf), new(*repository.UserRepo)),
	usecase.NewUserUsecase,
	wire.Bind(new(usecase.UserUsecaseItf), new(*usecase.UserUsecase)),
	handler.NewUserHandler,
)
