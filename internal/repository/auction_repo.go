package repository

import (
	"context"

	"github.com/Hilaladiii/aureus/internal/model"
	tx "github.com/Hilaladiii/aureus/pkg/config"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AuctionRepoItf interface {
	Create(ctx context.Context, auction *model.AuctionItem) error
	Update(ctx context.Context, auction *model.AuctionItem) error
	Delete(ctx context.Context, auctionID string) error
	GetAll(ctx context.Context) ([]model.AuctionItem, error)
	GetByID(ctx context.Context, auctionID string) (*model.AuctionItem, error)
	GetByIDWithLock(ctx context.Context, auctionID string) (*model.AuctionItem, error)
	GetByAuctioneerID(ctx context.Context, auctioneerID string) ([]model.AuctionItem, error)
}

type AuctionRepo struct {
	db *gorm.DB
}

func NewAuctionRepo(db *gorm.DB) *AuctionRepo {
	return &AuctionRepo{db}
}

func (r *AuctionRepo) Create(ctx context.Context, auction *model.AuctionItem) error {
	db := tx.ExtractTx(ctx, r.db)
	err := db.WithContext(ctx).Create(auction).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *AuctionRepo) Update(ctx context.Context, auction *model.AuctionItem) error {
	db := tx.ExtractTx(ctx, r.db)
	err := db.WithContext(ctx).Save(auction).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *AuctionRepo) Delete(ctx context.Context, auctionID string) error {
	db := tx.ExtractTx(ctx, r.db)
	err := db.WithContext(ctx).Delete(&model.AuctionItem{}, "id = ?", auctionID).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *AuctionRepo) GetAll(ctx context.Context) ([]model.AuctionItem, error) {
	var auctions []model.AuctionItem
	db := tx.ExtractTx(ctx, r.db)
	err := db.WithContext(ctx).Preload("AuctionImage").Joins("Category").Find(&auctions).Error
	if err != nil {
		return nil, err
	}

	return auctions, nil
}

func (r *AuctionRepo) GetByID(ctx context.Context, auctionID string) (*model.AuctionItem, error) {
	var auction model.AuctionItem
	db := tx.ExtractTx(ctx, r.db)
	err := db.WithContext(ctx).Preload("AuctionImage").Joins("Category").First(&auction, "auction_items.id = ?", auctionID).Error
	if err != nil {
		return nil, err
	}
	return &auction, nil
}

func (r *AuctionRepo) GetByIDWithLock(ctx context.Context, auctionID string) (*model.AuctionItem, error) {
	var auction model.AuctionItem
	db := tx.ExtractTx(ctx, r.db)
	err := db.WithContext(ctx).
		Clauses(clause.Locking{
			Strength: "UPDATE",
			Table:    clause.Table{Name: "auction_items"},
		}).
		Preload("AuctionImage").Joins("Category").
		First(&auction, "auction_items.id = ?", auctionID).Error
	if err != nil {
		return nil, err
	}
	return &auction, nil
}

func (r *AuctionRepo) GetByAuctioneerID(ctx context.Context, auctioneerID string) ([]model.AuctionItem, error) {
	var auctions []model.AuctionItem
	err := r.db.WithContext(ctx).Find(&auctions, "auction_items.auctioneerId = ?", auctioneerID).Error
	if err != nil {
		return nil, err
	}

	return auctions, nil
}
