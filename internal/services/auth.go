package services

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/VicAlexandre/pds-backend/internal/models"
)

type RegisterInput struct {
	Name     string `json:"name"`
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

type AuthService struct {
	UserModel *models.UserModel
	// TODO: TokenService (JWT or similar)
}

func NewAuthService(userModel *models.UserModel) *AuthService {
	return &AuthService{UserModel: userModel}
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*models.User, error) {
	if input.Name == "" || input.Email == "" || input.Password == "" {
		return nil, errors.New("name, email and password required")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := s.UserModel.Insert(ctx, input.Name, input.Email, string(hashed))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (*Tokens, error) {
	if input.Email == "" || input.Password == "" {
		return nil, errors.New("invalid credentials")
	}

	user, err := s.UserModel.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// TODO: replace with proper JWT generation
	tokens := &Tokens{
		AccessToken:  "fake-access-token",
		RefreshToken: "fake-refresh-token",
		ExpiresAt:    time.Now().Add(15 * time.Minute),
	}
	return tokens, nil
}

func (s *AuthService) Logout(ctx context.Context) error {
	// TODO: invalidate refresh token/session
	return nil
}
