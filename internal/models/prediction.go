package models

import (
	"errors"
	"time"
)

// Prediction represents a prediction made by a user
type Prediction struct {
	ID            int       `json:"id" db:"id"`
	UserID        int       `json:"user_id" db:"user_id"`
	Title         string    `json:"title" db:"title"`
	Description   string    `json:"description" db:"description"`
	Category      string    `json:"category" db:"category"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	EndDate       time.Time `json:"end_date" db:"end_date"`
	AgreeCount    int       `json:"agree_count" db:"agree_count"`
	DisagreeCount int       `json:"disagree_count" db:"disagree_count"`
	// Optional fields that might be populated with JOIN queries
	UserName string `json:"user_name,omitempty" db:"user_name"`
	UserVote *bool  `json:"user_vote,omitempty" db:"user_vote"`
}

// PredictionRequest represents the data needed to create a new prediction
type PredictionRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Category    string    `json:"category" binding:"required"`
	EndDate     time.Time `json:"end_date" binding:"required"`
}

// Validate performs validation on the prediction model
func (p *Prediction) Validate() error {
	// Title validation
	if p.Title == "" {
		return errors.New("title: title is required")
	}
	if len(p.Title) > 200 {
		return errors.New("title: title must be less than 200 characters")
	}

	// Description validation
	if p.Description == "" {
		return errors.New("description: description is required")
	}

	// Category validation
	if p.Category == "" {
		return errors.New("category: category is required")
	}

	// End date validation
	now := time.Now()
	if p.EndDate.Before(now) {
		return errors.New("end_date: end date must be in the future")
	}

	// Maximum 10 years in the future
	maxDate := now.AddDate(10, 0, 0)
	if p.EndDate.After(maxDate) {
		return errors.New("end_date: end date must be less than 10 years in the future")
	}

	return nil
}
