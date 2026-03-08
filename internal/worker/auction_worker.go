package worker

import (
	"context"
	"log"
	"strings"

	"github.com/Hilaladiii/aureus/internal/usecase"
	"github.com/redis/go-redis/v9"
)

type AuctionWorker struct {
	redisClient *redis.Client
	usecase     usecase.AuctionUsecaseItf
}

func NewAuctionWorker(redisClient *redis.Client, usecase usecase.AuctionUsecaseItf) *AuctionWorker {
	return &AuctionWorker{redisClient, usecase}
}

func (w *AuctionWorker) Start(ctx context.Context) {
	pubsub := w.redisClient.Subscribe(ctx, "__keyevent@0__:expired")

	go func() {
		defer pubsub.Close()
		ch := pubsub.Channel()

		for msg := range ch {
			if strings.HasPrefix(msg.Payload, "auction:expire:") {
				auctionID := strings.TrimPrefix(msg.Payload, "auction:expire:")
				err := w.usecase.FinalizeAuction(ctx, auctionID)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}()
}
