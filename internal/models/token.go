package models

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TODO: move to env variable
var jwtKey = []byte("verysecretkey")

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

type JWTModel struct {
}

type Token struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
	IssuedAt    time.Time `json:"issued_at"`
}

func (m *JWTModel) GenerateJWT(userID int64, duration time.Duration) (*Token, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, err
	}

	return &Token{
		AccessToken: tokenStr,
		ExpiresAt:   claims.ExpiresAt.Time,
		IssuedAt:    claims.IssuedAt.Time,
	}, nil
}

func (m *JWTModel) ParseJWT(tokenStr string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		return jwtKey, nil
	}, jwt.WithLeeway(5*time.Second))
	if err != nil {
		return nil, fmt.Errorf("could not parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
