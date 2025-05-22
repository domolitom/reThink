package models

import (
	"time"
)

// Result represents the outcome of a prediction after its end date
type Result struct {
	ID           int       `json:"id" db:"id"`
	PredictionID int       `json:"prediction_id" db:"prediction_id"`
	Outcome      bool      `json:"outcome" db:"outcome"` // true = came true, false = didn't come true
	EvidenceURL  string    `json:"evidence_url" db:"evidence_url"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// ResultRequest represents the data needed to submit a prediction result
type ResultRequest struct {
	Outcome     bool   `json:"outcome" binding:"required"`
	EvidenceURL string `json:"evidence_url" binding:"required"`
}
