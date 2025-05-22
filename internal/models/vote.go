package models

import (
	"errors"
	"time"
)

// Vote represents a user's vote on a prediction
type Vote struct {
	ID           int       `json:"id" db:"id"`
	UserID       int       `json:"user_id" db:"user_id"`
	PredictionID int       `json:"prediction_id" db:"prediction_id"`
	Value        bool      `json:"value" db:"value"` // true = agree, false = disagree
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// VoteRequest represents the data needed to vote on a prediction
type VoteRequest struct {
	Value bool `json:"value" binding:"required"`
}

// Validate performs validation on the vote model
func (v *Vote) Validate() error {
	// User ID validation
	if v.UserID <= 0 {
		return errors.New("user_id: user id must be positive")
	}

	// Prediction ID validation
	if v.PredictionID <= 0 {
		return errors.New("prediction_id: prediction id must be positive")
	}

	return nil
}
