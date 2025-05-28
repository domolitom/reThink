package models

import (
	"errors"
	"time"

	"github.com/domolitom/reThink/utils"
)

// User represents a user in the system
type User struct {
	ID        uint      `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"` // never expose in JSON
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// RegisterRequest represents the data needed to register a new user
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginRequest represents the data needed for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the data returned after successful login
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// Validate performs validation on the user model
func (u *User) Validate() error {
	// Name validation
	if u.Name == "" {
		return errors.New("name: name is required")
	}

	// Email validation
	if u.Email == "" {
		return errors.New("email: email is required")
	}
	if !utils.ValidateEmail(u.Email) {
		return errors.New("email: invalid email format")
	}

	// Password validation
	if u.Password == "" {
		return errors.New("password: password is required")
	}
	if len(u.Password) < 8 {
		return errors.New("password: password must be at least 8 characters long")
	}

	return nil
}

// UserProfile represents public user information
type UserProfile struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// ToProfile converts a User to a UserProfile
func (u *User) ToProfile() UserProfile {
	return UserProfile{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}
