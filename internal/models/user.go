package models

import (
	"time"
)

type User struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	Username        string    `json:"username" gorm:"unique;not null"`
	Email           string    `json:"email" gorm:"unique;not null"`
	PasswordHash    string    `json:"-" gorm:"not null"`
	Bio             string    `json:"bio"`
	PredictionScore float64   `json:"prediction_score" gorm:"default:0"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
