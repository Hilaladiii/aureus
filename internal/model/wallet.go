package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Wallet struct {
	ID            string          `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	ActiveBalance decimal.Decimal `gorm:"type:decimal(20,2);default:0"`
	HeldBalance   decimal.Decimal `gorm:"type:decimal(20,2);default:0"`
	CreatedAt     time.Time       `gorm:"autoCreateTime"`
	UpdatedAt     time.Time       `gorm:"autoCreateTime;autoUpdateTime"`
	UserID        string          `gorm:"type:uuid;uniqueIndex"`
	User          *User           `gorm:"foreignKey:UserID;references:ID"`
}

type WalletCreateRequest struct {
	ActiveBalance decimal.Decimal `form:"activeBalance" validate:"required"`
}

type WalletUpdateRequest struct {
	ActiveBalance decimal.Decimal `form:"activeBalance,omitempty" validate:"omitempty"`
	HeldBalance   decimal.Decimal `form:"heldBalance,omitempty" validate:"omitempty"`
}

type WalletResource struct {
	ID            string          `json:"id"`
	ActiveBalance decimal.Decimal `json:"activeBalance"`
	HeldBalance   decimal.Decimal `json:"heldBalance"`
	UserID        string          `json:"userId"`
}

func (w *Wallet) Resource() WalletResource {
	return WalletResource{
		ID:            w.ID,
		ActiveBalance: w.ActiveBalance,
		HeldBalance:   w.HeldBalance,
		UserID:        w.UserID,
	}
}

func WalletResources(wallets []Wallet) []WalletResource {
	if len(wallets) == 0 {
		return []WalletResource{}
	}

	responses := make([]WalletResource, 0, len(wallets))
	for i := range wallets {
		responses = append(responses, wallets[i].Resource())
	}
	return responses
}
