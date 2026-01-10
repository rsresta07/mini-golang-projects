package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("CHANGE_THIS_SECRET")

// GenerateToken generates a JWT token for a given user ID. The token
// will expire in 24 hours.
func GenerateToken(userID uint) (string, error) {
	// Create a new JWT token
	// Set the expiration time to 24 hours from now
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken takes a JWT token string and returns the associated user ID if the token is valid.
// If the token is invalid or has expired, it returns an error.
func ParseToken(tokenStr string) (uint, error) {
	// Parse the JWT token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return 0, err
	}

	// Extract the user ID from the token
	claims := token.Claims.(jwt.MapClaims)
	return uint(claims["user_id"].(float64)), nil
}
