package model

import "time"

type Role string

const (
	ADMIN Role = "ADMIN"
	USER  Role = "USER"
)

type User struct {
	ID         string       `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Email      string       `gorm:"unique;type:varchar(255)"`
	Username   string       `gorm:"unique;type:varchar(50)"`
	Password   string       `gorm:"type:varchar(255)"`
	CreatedAt  time.Time    `gorm:"autoCreateTime"`
	UpdatedAt  time.Time    `gorm:"autoCreateTime;autoUpdateTime"`
	Role       Role         `gorm:"type:varchar(10);default:'USER'"`
	Wallet     Wallet       `gorm:"foreignKey:UserID;references:ID"`
	BidHistory []BidHistory `gorm:"foreignKey:UserID;references:ID"`
}
