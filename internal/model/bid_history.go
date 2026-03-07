package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Status string

const (
	WINNING Status = "WINNING"
	OUTBID  Status = "OUTBID"
	WON     Status = "WON"
)

type BidHistory struct {
	ID          string          `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	BidAmount   decimal.Decimal `gorm:"type:decimal(20,2)"`
	Status      Status          `gorm:"type:varchar(10)check:status IN ('WINNING', 'OUTBID','WON')"`
	CreatedAt   time.Time       `gorm:"autoCreateTime"`
	User        *User           `gorm:"foreignKey:UserID;references:ID"`
	AuctionItem *AuctionItem    `gorm:"foreignKey:ItemID;references:ID"`
	ItemID      string
	UserID      string
}

type BidHistoryCreateRequest struct {
	BidAmount decimal.Decimal `form:"bidAmount" validate:"required"`
	Status    Status          `form:"status" validate:"omitempty"`
	UserID    string          `form:"userID" validate:"required"`
	ItemID    string          `form:"itemID" validate:"required"`
}

type BidHistoryUpdateRequest struct {
	BidAmount decimal.Decimal `form:"bidAmount" validate:"omitempty"`
	Status    Status          `form:"status" validate:"omitempty"`
}

type BidHistoryResource struct {
	ID        string          `json:"id"`
	BidAmount decimal.Decimal `json:"bidAmount"`
	Status    Status          `json:"status"`
	CreatedAt time.Time       `json:"createdAt"`
	ItemID    string          `json:"itemId"`
	UserID    string          `json:"userId"`
}

func (b *BidHistory) Resource() BidHistoryResource {
	return BidHistoryResource{
		ID:        b.ID,
		BidAmount: b.BidAmount,
		Status:    b.Status,
		CreatedAt: b.CreatedAt,
		ItemID:    b.ItemID,
		UserID:    b.UserID,
	}
}

func BidResources(bids []BidHistory) []BidHistoryResource {
	if len(bids) == 0 {
		return []BidHistoryResource{}
	}

	resources := make([]BidHistoryResource, 0, len(bids))
	for i := range bids {
		resources = append(resources, bids[i].Resource())
	}

	return resources
}
