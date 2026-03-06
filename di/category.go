package di

import (
	"github.com/Hilaladiii/aureus/internal/delivery/handler"
	"github.com/Hilaladiii/aureus/internal/repository"
	"github.com/Hilaladiii/aureus/internal/usecase"
	"github.com/google/wire"
)

var CategorySet = wire.NewSet(
	repository.NewCategoryRepo,
	wire.Bind(new(repository.CategoryRepoItf), new(*repository.CategoryRepo)),
	usecase.NewCategoryUsecase,
	wire.Bind(new(usecase.CategoryUsecaseItf), new(*usecase.CategoryUsecase)),
	handler.NewCategoryHandler)
