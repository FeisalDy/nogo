package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the claims in the JWT token
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var (
	JWTSecret       = []byte("your-secret-key-change-this-in-production")
	TokenExpiration = 24 * time.Hour
)

// GenerateToken generates a new JWT token for a user
func GenerateToken(userID uint, email, username string) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		Email:    email,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (any, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken generates a new token with extended expiration
func RefreshToken(oldTokenString string) (string, error) {
	claims, err := ValidateToken(oldTokenString)
	if err != nil {
		return "", err
	}

	// Generate new token with the same user info
	return GenerateToken(claims.UserID, claims.Email, claims.Username)
}
