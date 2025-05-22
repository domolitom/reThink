// utils_test.go
package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordHashing(t *testing.T) {
	password := "securepassword123"

	// Test hash generation
	hashedPassword, err := HashPassword(password)
	assert.NoError(t, err)
	assert.NotEqual(t, password, hashedPassword)
	assert.NotEmpty(t, hashedPassword)

	// Test password verification
	valid := CheckPasswordHash(password, hashedPassword)
	assert.True(t, valid)

	// Test incorrect password
	invalid := CheckPasswordHash("wrongpassword", hashedPassword)
	assert.False(t, invalid)
}

func TestJWTTokenGeneration(t *testing.T) {
	userID := 123

	// Generate token
	token, err := GenerateToken(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify token
	extractedID, err := VerifyToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, extractedID)

	// Test invalid token
	_, err = VerifyToken("invalid.token.string")
	assert.Error(t, err)

	// Test expired token (if we can manipulate token expiry for testing)
	// This would require modifying the token generation function to accept a custom expiry
	// for testing purposes, or mocking time.Now()

	// Example (if implemented):
	// expiredToken, _ := GenerateTokenWithExpiry(userID, time.Now().Add(-24 * time.Hour))
	// _, err = VerifyToken(expiredToken)
	// assert.Error(t, err)
	// assert.Contains(t, err.Error(), "expired")
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email    string
		expected bool
	}{
		{"test@example.com", true},
		{"test.name@example.com", true},
		{"test+123@example.com", true},
		{"test@subdomain.example.com", true},
		{"test@example.co.uk", true},
		{"test", false},
		{"test@", false},
		{"@example.com", false},
		{"test@example", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.email, func(t *testing.T) {
			result := ValidateEmail(tt.email)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParsePaginationParams(t *testing.T) {
	tests := []struct {
		name          string
		pageStr       string
		limitStr      string
		expectedPage  int
		expectedLimit int
	}{
		{
			name:          "Valid Values",
			pageStr:       "3",
			limitStr:      "20",
			expectedPage:  3,
			expectedLimit: 20,
		},
		{
			name:          "Default Values",
			pageStr:       "",
			limitStr:      "",
			expectedPage:  1,
			expectedLimit: 10,
		},
		{
			name:          "Negative Page",
			pageStr:       "-1",
			limitStr:      "10",
			expectedPage:  1,
			expectedLimit: 10,
		},
		{
			name:          "Zero Page",
			pageStr:       "0",
			limitStr:      "10",
			expectedPage:  1,
			expectedLimit: 10,
		},
		{
			name:          "Negative Limit",
			pageStr:       "1",
			limitStr:      "-5",
			expectedPage:  1,
			expectedLimit: 10,
		},
		{
			name:          "Zero Limit",
			pageStr:       "1",
			limitStr:      "0",
			expectedPage:  1,
			expectedLimit: 10,
		},
		{
			name:          "Exceeding Max Limit",
			pageStr:       "1",
			limitStr:      "101",
			expectedPage:  1,
			expectedLimit: 100,
		},
		{
			name:          "Invalid Page Format",
			pageStr:       "abc",
			limitStr:      "10",
			expectedPage:  1,
			expectedLimit: 10,
		},
		{
			name:          "Invalid Limit Format",
			pageStr:       "1",
			limitStr:      "abc",
			expectedPage:  1,
			expectedLimit: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page, limit := ParsePaginationParams(tt.pageStr, tt.limitStr)
			assert.Equal(t, tt.expectedPage, page)
			assert.Equal(t, tt.expectedLimit, limit)
		})
	}
}
