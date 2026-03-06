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
