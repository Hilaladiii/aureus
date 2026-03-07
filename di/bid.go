package di

import (
	"github.com/Hilaladiii/aureus/internal/repository"
	"github.com/google/wire"
)

var BidSet = wire.NewSet(
	repository.NewBidRepo,
	wire.Bind(new(repository.BidRepoItf), new(*repository.BidRepo)),
)
