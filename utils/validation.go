package utils

import (
	"errors"
	"regexp"
	"strconv"
)

// Common errors
var (
	ErrNotFound  = errors.New("resource not found")
	ErrDuplicate = errors.New("duplicate resource")
)

// ValidateEmail checks if an email follows a valid format
func ValidateEmail(email string) bool {
	// Basic regex pattern for email validation
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(email)
}

// ParsePaginationParams parses and validates page and limit parameters for pagination
func ParsePaginationParams(pageStr, limitStr string) (int, int) {
	// Default values
	page := 1
	limit := 10
	maxLimit := 100

	// Parse page
	if pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	// Parse limit
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
			if limit > maxLimit {
				limit = maxLimit
			}
		}
	}

	return page, limit
}
