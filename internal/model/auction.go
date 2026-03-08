package model

import (
	"mime/multipart"
	"time"

	"github.com/shopspring/decimal"
)

type AuctionStatus string

type AuctionImage struct {
	ID          string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	ItemID      string `gorm:"type:uuid"`
	ImageUrl    string
	CreatedAt   time.Time   `gorm:"autoCreateTime"`
	AuctionItem AuctionItem `gorm:"foreignKey:ItemID;references:ID"`
}

type AuctionImageResource struct {
	ID        string `json:"id"`
	ItemID    string `json:"itemId"`
	ImageUrl  string `json:"imageUrl"`
	CreatedAt string `json:"createdAt"`
}

func (a *AuctionImage) Resource() AuctionImageResource {
	if a == nil {
		return AuctionImageResource{}
	}
	return AuctionImageResource{
		ID:        a.ID,
		ItemID:    a.ItemID,
		ImageUrl:  a.ImageUrl,
		CreatedAt: a.CreatedAt.Format(time.RFC3339),
	}
}

func AuctionImageResources(images []AuctionImage) []AuctionImageResource {
	if len(images) == 0 {
		return []AuctionImageResource{}
	}

	resources := make([]AuctionImageResource, 0, len(images))
	for i := range images {
		resources = append(resources, images[i].Resource())
	}
	return resources
}

type AuctionItem struct {
	ID           string          `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name         string          `gorm:"type:varchar(50)"`
	Description  *string         `gorm:"type:text"`
	StartPrice   decimal.Decimal `gorm:"type:decimal(20,2)"`
	CurrentPrice decimal.Decimal `gorm:"type:decimal(20,2)"`
	BidIncrement decimal.Decimal `gorm:"type:decimal(20,2)"`
	StartTime    time.Time
	EndTime      time.Time
	AuctionImage []AuctionImage `gorm:"foreignKey:ItemID;references:ID"`
	CategoryID   string         `gorm:"type:uuid;"`
	Category     *Category      `gorm:"foreignKey:CategoryID;references:ID"`
	BidHistory   []BidHistory   `gorm:"foreignKey:ItemID;references:ID"`
	AuctioneerID string         `gorm:"type:uuid"`
	User         User           `gorm:"foreignKey:AuctioneerID;references:ID"`
}

type LeaderboardResource struct {
	CensoredName string  `json:"name"`
	BidAmount    float64 `json:"bidAmount"`
}

type AuctionCreateRequest struct {
	Name         string                  `form:"name" validate:"required"`
	Description  *string                 `form:"description,omitempty" validate:"omitempty"`
	StartPrice   decimal.Decimal         `form:"startPrice" validate:"required"`
	BidIncrement decimal.Decimal         `form:"bidIncrement" validate:"required"`
	StartTime    time.Time               `form:"startTime" validate:"required"`
	EndTime      time.Time               `form:"endTime" validate:"required"`
	CategoryID   string                  `form:"categoryId" validate:"required"`
	Images       []*multipart.FileHeader `form:"images" validate:"required"`
}

type AuctionBidRequest struct {
	CurrentPrice decimal.Decimal `form:"currentPrice" validate:"required"`
}

type AuctionUpdateRequest struct {
	Name         string          `form:"name,omitempty" validate:"omitempty"`
	Description  *string         `form:"description,omitempty" validate:"omitempty"`
	StartPrice   decimal.Decimal `form:"startPrice,omitempty" validate:"omitempty"`
	CurrentPrice decimal.Decimal `form:"currentPrice,omitempty" validate:"omitempty"`
	StartTime    time.Time       `form:"startTime,omitempty" validate:"omitempty"`
	EndTime      time.Time       `form:"endTime,omitempty" validate:"omitempty"`
	CategoryID   string          `form:"categoryId" validate:"required"`
}

type AuctionResource struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Description  *string                `json:"description"`
	StartPrice   decimal.Decimal        `json:"startPrice"`
	BidIncrement decimal.Decimal        `json:"bidIncrement"`
	CurrentPrice decimal.Decimal        `json:"currentPrice"`
	StartTime    time.Time              `json:"startTime"`
	EndTime      time.Time              `json:"endTime"`
	Category     *Category              `json:"category"`
	Images       []AuctionImageResource `json:"images"`
	AuctioneerID string                 `json:"auctioneerId"`
}

func (a *AuctionItem) Resource() AuctionResource {
	return AuctionResource{
		ID:           a.ID,
		Name:         a.Name,
		Description:  a.Description,
		StartPrice:   a.StartPrice,
		CurrentPrice: a.CurrentPrice,
		BidIncrement: a.BidIncrement,
		StartTime:    a.StartTime,
		EndTime:      a.EndTime,
		Category:     a.Category,
		AuctioneerID: a.AuctioneerID,
		Images:       AuctionImageResources(a.AuctionImage),
	}
}

func AuctionResources(auctions []AuctionItem) []AuctionResource {
	if len(auctions) == 0 {
		return []AuctionResource{}
	}

	resources := make([]AuctionResource, 0, len(auctions))
	for i := range auctions {
		resources = append(resources, auctions[i].Resource())
	}

	return resources
}
