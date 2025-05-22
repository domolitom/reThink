// models_test.go
package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUserValidation(t *testing.T) {
	tests := []struct {
		name        string
		user        User
		expectedErr bool
		errField    string
	}{
		{
			name: "Valid User",
			user: User{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "securepassword",
			},
			expectedErr: false,
		},
		{
			name: "Empty Name",
			user: User{
				Name:     "",
				Email:    "test@example.com",
				Password: "securepassword",
			},
			expectedErr: true,
			errField:    "name",
		},
		{
			name: "Invalid Email",
			user: User{
				Name:     "Test User",
				Email:    "notavalidemail",
				Password: "securepassword",
			},
			expectedErr: true,
			errField:    "email",
		},
		{
			name: "Empty Email",
			user: User{
				Name:     "Test User",
				Email:    "",
				Password: "securepassword",
			},
			expectedErr: true,
			errField:    "email",
		},
		{
			name: "Short Password",
			user: User{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "short",
			},
			expectedErr: true,
			errField:    "password",
		},
		{
			name: "Empty Password",
			user: User{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "",
			},
			expectedErr: true,
			errField:    "password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			
			if tt.expectedErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errField)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPredictionValidation(t *testing.T) {
	now := time.Now()
	
	tests := []struct {
		name        string
		prediction  Prediction
		expectedErr bool
		errField    string
	}{
		{
			name: "Valid Prediction",
			prediction: Prediction{
				UserID:      1,
				Title:       "Stock Market Prediction",
				Description: "I predict the S&P 500 will reach 5000 by end of year",
				Category:    "Finance",
				EndDate:     now.AddDate(0, 3, 0), // 3 months in the future
			},
			expectedErr: false,
		},
		{
			name: "Missing Title",
			prediction: Prediction{
				UserID:      1,
				Title:       "",
				Description: "Some description",
				Category:    "Finance",
				EndDate:     now.AddDate(0, 3, 0),
			},
			expectedErr: true,
			errField:    "title",
		},
		{
			name: "Title Too Long",
			prediction: Prediction{
				UserID:      1,
				Title:       string(make([]byte, 201)), // 201 characters
				Description: "Some description",
				Category:    "Finance",
				EndDate:     now.AddDate(0, 3, 0),
			},
			expectedErr: true,
			errField:    "title",
		},
		{
			name: "Missing Description",
			prediction: Prediction{
				UserID:      1,
				Title:       "Valid Title",
				Description: "",
				Category:    "Finance",
				EndDate:     now.AddDate(0, 3, 0),
			},
			expectedErr: true,
			errField:    "description",
		},
		{
			name: "Missing Category",
			prediction: Prediction{
				UserID:      1,
				Title:       "Valid Title",
				Description: "Valid description",
				Category:    "",
				EndDate:     now.AddDate(0, 3, 0),
			},
			expectedErr: true,
			errField:    "category",
		},
		{
			name: "End Date in Past",
			prediction: Prediction{
				UserID:      1,
				Title:       "Valid Title",
				Description: "Valid description",
				Category:    "Finance",
				EndDate:     now.AddDate(0, 0, -1), // Yesterday
			},
			expectedErr: true,
			errField:    "end_date",
		},
		{
			name: "End Date Too Far in Future",
			prediction: Prediction{
				UserID:      1,
				Title:       "Valid Title",
				Description: "Valid description",
				Category:    "Finance",
				EndDate:     now.AddDate(11, 0, 0), // 11 years in future
			},
			expectedErr: true,
			errField:    "end_date",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.prediction.Validate()
			
			if tt.expectedErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errField)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestVoteValidation(t *testing.T) {
	tests := []struct {
		name        string
		vote        Vote
		expectedErr bool
		errField    string
	}{
		{
			name: "Valid Vote - Agree",
			vote: Vote{
				UserID:       1,
				PredictionID: 1,
				Value:        true,
			},
			expectedErr: false,
		},
		{
			name: "Valid Vote - Disagree",
			vote: Vote{
				UserID:       1,
				PredictionID: 1,
				Value:        false,
			},
			expectedErr: false,
		},
		{
			name: "Invalid User ID",
			vote: Vote{
				UserID:       0,
				PredictionID: 1,
				Value:        true,
			},
			expectedErr: true,
			errField:    "user_id",
		},
		{
			name: "Invalid Prediction ID",
			vote: Vote{
				UserID:       1,
				PredictionID: 0,
				Value:        true,
			},
			expectedErr: true,
			errField:    "prediction_id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.vote.Validate()
			
			if tt.expectedErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errField)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}