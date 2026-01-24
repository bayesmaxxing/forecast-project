package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken   = errors.New("invalid token")
	ErrNotInitialized = errors.New("auth package not initialized")
	secretKey         []byte
)

type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Init initializes the auth package with the JWT secret.
// Must be called before using GenerateToken or ValidateToken.
func Init(secret []byte) error {
	if len(secret) < 32 {
		return fmt.Errorf("JWT secret must be at least 32 characters, got %d", len(secret))
	}
	secretKey = secret
	return nil
}

func GenerateToken(userID int64, username string) (string, error) {
	if secretKey == nil {
		return "", ErrNotInitialized
	}

	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ValidateToken(tokenString string) (*Claims, error) {
	if secretKey == nil {
		return nil, ErrNotInitialized
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return token.Claims.(*Claims), nil
}
