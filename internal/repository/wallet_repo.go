package repository

import (
	"context"

	"github.com/Hilaladiii/aureus/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletRepoItf interface {
	Create(ctx context.Context, wallet *model.Wallet) error
	Update(ctx context.Context, wallet *model.Wallet, walletID string) error
	AddBalance(ctx context.Context, amount float64, walletID string) error
	GetBalance(ctx context.Context, walletID string) (*model.Wallet, error)
	GetAll(ctx context.Context) ([]model.Wallet, error)
	GetByID(ctx context.Context, walletID string) (*model.Wallet, error)
	GetByUserID(ctx context.Context, userID string) (*model.Wallet, error)
}

type WalletRepo struct {
	db *gorm.DB
}

func NewWalletRepo(db *gorm.DB) *WalletRepo {
	return &WalletRepo{db: db}
}

func (r *WalletRepo) Create(ctx context.Context, wallet *model.Wallet) error {
	err := r.db.WithContext(ctx).Create(&wallet).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *WalletRepo) Update(ctx context.Context, wallet *model.Wallet, walletID string) error {
	err := r.db.WithContext(ctx).Save(wallet).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *WalletRepo) AddBalance(ctx context.Context, amount float64, walletID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var wallet model.Wallet

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", walletID).
			First(&wallet).Error; err != nil {
			return err
		}

		wallet.ActiveBalance += amount

		if err := tx.Save(&wallet).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *WalletRepo) GetBalance(ctx context.Context, walletID string) (*model.Wallet, error) {
	var wallet model.Wallet
	err := r.db.WithContext(ctx).Select("active_balance", "held_balance").First(&wallet, "id = ?", walletID).Error
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (r *WalletRepo) GetAll(ctx context.Context) ([]model.Wallet, error) {
	var wallets []model.Wallet
	err := r.db.WithContext(ctx).Find(&wallets).Error
	if err != nil {
		return nil, err
	}

	return wallets, err
}

func (r *WalletRepo) GetByID(ctx context.Context, walletID string) (*model.Wallet, error) {
	var wallet model.Wallet
	err := r.db.WithContext(ctx).First(&wallet, "id = ?", walletID).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *WalletRepo) GetByUserID(ctx context.Context, userID string) (*model.Wallet, error) {
	var wallet model.Wallet
	err := r.db.WithContext(ctx).First(&wallet, "user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}
