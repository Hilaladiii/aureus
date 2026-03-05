package model

import (
	"time"
)

type Category struct {
	ID          string        `gorm:"primaryKey;json:id;type:uuid;default:uuid_generate_v4()"`
	Name        string        `gorm:"type:varchar(50);uniqueIndex"`
	Description *string       `gorm:"type:text"`
	CreatedAt   time.Time     `gorm:"autoCreateTime"`
	AuctionItem []AuctionItem `gorm:"foreignKey:CategoryID;references:ID"`
}
