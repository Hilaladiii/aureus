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

type UserRegisterRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UserUpdateRequest struct {
	Username string `json:"username,omitempty" validate:"omitempty"`
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
	Password string `json:"password,omitempty" validate:"omitempty,min=8"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Role      Role      `json:"role"`
}
