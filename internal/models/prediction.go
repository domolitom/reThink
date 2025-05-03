package models

import (
	"time"
)

type Prediction struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id"`
	User       User      `json:"user" gorm:"foreignKey:UserID"`
	MarketID   uint      `json:"market_id"`
	Market     Market    `json:"market" gorm:"foreignKey:MarketID"`
	Prediction bool      `json:"prediction"`
	Confidence float64   `json:"confidence" gorm:"default:50"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
