package repository

import (
	"context"

	"github.com/Hilaladiii/aureus/internal/model"
	tx "github.com/Hilaladiii/aureus/pkg/config"
	"gorm.io/gorm"
)

type BidRepoItf interface {
	Create(ctx context.Context, bid *model.BidHistory) error
	Update(ctx context.Context, bid *model.BidHistory, bidID string) error
	Delete(ctx context.Context, bidID string) error
	GetAll(ctx context.Context) ([]model.BidHistory, error)
	GetByID(ctx context.Context, bidID string) (*model.BidHistory, error)
	GetByAuctionID(ctx context.Context, auctionID string) ([]model.BidHistory, error)
	GetHighestBid(ctx context.Context, auctionID string) (*model.BidHistory, error)
}

type BidRepo struct {
	db *gorm.DB
}

func NewBidRepo(db *gorm.DB) *BidRepo {
	return &BidRepo{db}
}

func (r *BidRepo) Create(ctx context.Context, bid *model.BidHistory) error {
	db := tx.ExtractTx(ctx, r.db)
	err := db.WithContext(ctx).Create(bid).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *BidRepo) Update(ctx context.Context, bid *model.BidHistory, bidID string) error {
	db := tx.ExtractTx(ctx, r.db)
	err := db.WithContext(ctx).Save(bid).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *BidRepo) Delete(ctx context.Context, bidID string) error {
	db := tx.ExtractTx(ctx, r.db)
	err := db.WithContext(ctx).Delete(&model.BidHistory{}, "id = ?", bidID).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *BidRepo) GetAll(ctx context.Context) ([]model.BidHistory, error) {
	var bids []model.BidHistory
	db := tx.ExtractTx(ctx, r.db)
	err := db.WithContext(ctx).Find(&bids).Error
	if err != nil {
		return nil, err
	}
	return bids, nil
}

func (r *BidRepo) GetByID(ctx context.Context, bidID string) (*model.BidHistory, error) {
	var bid model.BidHistory
	db := tx.ExtractTx(ctx, r.db)
	err := db.WithContext(ctx).First(bid, "id = ?", bidID).Error
	if err != nil {
		return nil, err
	}
	return &bid, nil
}

func (r *BidRepo) GetByAuctionID(ctx context.Context, auctionID string) ([]model.BidHistory, error) {
	var bidHistories []model.BidHistory
	db := tx.ExtractTx(ctx, r.db)
	err := db.WithContext(ctx).Find(&bidHistories, "bid_histories.item_id = ?", auctionID).Error
	if err != nil {
		return nil, err
	}
	return bidHistories, nil
}

func (r *BidRepo) GetHighestBid(ctx context.Context, auctionID string) (*model.BidHistory, error) {
	var bid model.BidHistory
	db := tx.ExtractTx(ctx, r.db)
	err := db.WithContext(ctx).Where("item_id = ?", auctionID).Order("bid_amount DESC").First(&bid).Error
	if err != nil {
		return nil, err
	}
	return &bid, nil
}
