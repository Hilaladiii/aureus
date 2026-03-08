package repository

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type LeaderboardEntry struct {
	UserID    string  `json:"userId"`
	BidAmount float64 `json:"bidAmount"`
}

type AuctionCacheRepoItf interface {
	AddBidToLeaderboard(ctx context.Context, auctionID string, userID string, bidAmount float64) error
	GetTopBidders(ctx context.Context, auctionID string) ([]LeaderboardEntry, error)
}

type AuctionCacheRepo struct {
	redis *redis.Client
}

func NewAuctionCacheRepo(redis *redis.Client) *AuctionCacheRepo {
	return &AuctionCacheRepo{redis}
}

func (r *AuctionCacheRepo) AddBidToLeaderboard(ctx context.Context, auctionID string, userID string, bidAmount float64) error {
	key := fmt.Sprintf("leaderboard:auction:%s", auctionID)
	return r.redis.ZAdd(ctx, key, redis.Z{
		Score:  bidAmount,
		Member: userID,
	}).Err()
}

func (r *AuctionCacheRepo) GetTopBidders(ctx context.Context, auctionID string) ([]LeaderboardEntry, error) {
	key := fmt.Sprintf("leaderboard:auction:%s", auctionID)
	results, err := r.redis.ZRevRangeWithScores(ctx, key, 0, 9).Result()
	if err != nil {
		return nil, err
	}

	var leaderboard []LeaderboardEntry
	for _, z := range results {
		leaderboard = append(leaderboard, LeaderboardEntry{
			UserID:    z.Member.(string),
			BidAmount: z.Score,
		})
	}

	return leaderboard, nil
}
