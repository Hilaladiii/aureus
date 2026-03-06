package model

import (
	"time"
)

type Wallet struct {
	ID            string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	ActiveBalance float64   `gorm:"type:decimal(15,2);default:0"`
	HeldBalance   float64   `gorm:"type:decimal(15,2);default:0"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoCreateTime;autoUpdateTime"`
	UserID        string    `gorm:"type:uuid;uniqueIndex"`
	User          *User     `gorm:"foreignKey:UserID;references:ID"`
}

type WalletCreateRequest struct {
	ActiveBalance float64 `form:"activeBalance" validate:"required"`
}

type WalletUpdateRequest struct {
	ActiveBalance float64 `form:"activeBalance,omitempty" validate:"omitempty"`
	HeldBalance   float64 `form:"heldBalance,omitempty" validate:"omitempty"`
}

type WalletResource struct {
	ID            string  `json:"id"`
	ActiveBalance float64 `json:"activeBalance"`
	HeldBalance   float64 `json:"heldBalance"`
	UserID        string  `json:"userId"`
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
