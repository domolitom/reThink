package models

import (
	"time"
)

type MarketStatus string

const (
	MarketOpen     MarketStatus = "open"
	MarketClosed   MarketStatus = "closed"
	MarketResolved MarketStatus = "resolved"
)

type Market struct {
	ID          uint         `json:"id" gorm:"primaryKey"`
	Title       string       `json:"title" gorm:"not null"`
	Description string       `json:"description"`
	CreatorID   uint         `json:"creator_id"`
	Creator     User         `json:"creator" gorm:"foreignKey:CreatorID"`
	CloseDate   time.Time    `json:"close_date"`
	ResolveDate time.Time    `json:"resolve_date"`
	Status      MarketStatus `json:"status" gorm:"default:'open'"`
	Outcome     *bool        `json:"outcome"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}
