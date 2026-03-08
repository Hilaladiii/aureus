package di

import (
	"github.com/Hilaladiii/aureus/internal/delivery/handler"
	"github.com/Hilaladiii/aureus/internal/repository"
	"github.com/Hilaladiii/aureus/internal/usecase"
	"github.com/google/wire"
)

var AuctionSet = wire.NewSet(
	repository.NewAuctionRepo,
	wire.Bind(new(repository.AuctionRepoItf), new(*repository.AuctionRepo)),
	repository.NewAuctionCacheRepo,
	wire.Bind(new(repository.AuctionCacheRepoItf), new(*repository.AuctionCacheRepo)),
	usecase.NewAuctionUsecase,
	wire.Bind(new(usecase.AuctionUsecaseItf), new(*usecase.AuctionUsecase)),
	handler.NewAuctionHandler,
)
