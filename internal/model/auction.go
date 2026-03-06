package model

import (
	"time"
)

type AuctionStatus string

type AuctionItem struct {
	ID           string  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name         string  `gorm:"type:varchar(50)"`
	Description  *string `gorm:"type:text"`
	StartPrice   float64 `gorm:"type:decimal(15,2)"`
	CurrentPrice float64 `gorm:"type:decimal(15,2)"`
	StartTime    time.Time
	EndTime      time.Time
	CategoryID   string       `gorm:"type:uuid;"`
	Category     *Category    `gorm:"foreignKey:CategoryID;references:ID"`
	BidHistory   []BidHistory `gorm:"foreignKey:ItemID;references:ID"`
}
