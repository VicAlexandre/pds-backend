package services

import (
	"context"
	"errors"
	"time"
)

type RegisterInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Tokens struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type User struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

type AuthService struct {
	// TODO: add dependencies like DB, cache, etc.
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*User, error) {
	// TODO: validate, hash password, insert into DB
	if input.Email == "" || input.Password == "" {
		return nil, errors.New("email and password required")
	}

	user := &User{
		ID:    1,
		Email: input.Email,
	}
	return user, nil
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (*Tokens, error) {
	//TODO: validate credentials, check against DB and generate tokens
	if input.Email == "" || input.Password == "" {
		return nil, errors.New("invalid credentials")
	}

	tokens := &Tokens{
		AccessToken:  "fake-access-token",
		RefreshToken: "fake-refresh-token",
		ExpiresAt:    time.Now().Add(15 * time.Minute),
	}
	return tokens, nil
}

func (s *AuthService) Logout(ctx context.Context) error {
	// TODO: properly handle logout, e.g., invalidate tokens
	return nil
}
