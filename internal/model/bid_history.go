package model

import "time"

type Status string

const (
	WINNING Status = "WINNING"
	OUTBID  Status = "OUTBID"
	WON     Status = "WON"
)

type BidHistory struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	BidAmount   float64   `gorm:"type:decimal(15,2)"`
	Status      Status    `gorm:"type:varchar(10)check:status IN ('WINNING', 'OUTBID','WON')"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	ItemID      string
	UserID      string
	User        *User        `gorm:"foreignKey:UserID;references:ID"`
	AuctionItem *AuctionItem `gorm:"foreignKey:ItemID;references:ID"`
}
