package utils

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Custom errors
var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
)

// JWT secret key loaded from environment or using a default for development
var jwtSecret = getJWTSecret()

func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// Default secret for development - should be set properly in production
		secret = "rethink_default_jwt_secret_change_in_production"
	}
	return []byte(secret)
}

// JWTClaims represents the claims in the JWT token
type JWTClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken creates a JWT token for a user
func GenerateToken(userID int) (string, error) {
	// Default token duration is 24 hours
	expirationTime := time.Now().Add(24 * time.Hour)

	if expStr := os.Getenv("JWT_EXPIRATION_HOURS"); expStr != "" {
		if expHours, err := strconv.Atoi(expStr); err == nil {
			expirationTime = time.Now().Add(time.Duration(expHours) * time.Hour)
		}
	}

	claims := &JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// VerifyToken validates a JWT token and returns the user ID
func VerifyToken(tokenString string) (int, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		},
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return 0, ErrExpiredToken
		}
		return 0, ErrInvalidToken
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return 0, ErrInvalidToken
	}

	return claims.UserID, nil
}
