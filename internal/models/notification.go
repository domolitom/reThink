package models

import (
	"time"
)

// NotificationType defines the type of notification
type NotificationType string

const (
	NotificationVote       NotificationType = "vote"
	NotificationResult     NotificationType = "result"
	NotificationMention    NotificationType = "mention"
	NotificationEndingSoon NotificationType = "ending_soon"
)

// Notification represents a notification for a user
type Notification struct {
	ID        int              `json:"id" db:"id"`
	UserID    int              `json:"user_id" db:"user_id"`
	Type      NotificationType `json:"type" db:"type"`
	Message   string           `json:"message" db:"message"`
	Link      string           `json:"link" db:"link"`
	Read      bool             `json:"read" db:"read"`
	CreatedAt time.Time        `json:"created_at" db:"created_at"`
}
