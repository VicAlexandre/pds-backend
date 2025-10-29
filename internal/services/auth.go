package services

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/VicAlexandre/pds-backend/internal/models"
)

const TokenDuration = 60 * time.Minute

type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthService struct {
	UserModel  *models.UserModel
	TokenModel *models.JWTModel
}

func NewAuthService(userModel *models.UserModel, tokenModel *models.JWTModel) *AuthService {
	return &AuthService{
		UserModel:  userModel,
		TokenModel: tokenModel,
	}
}

func (s *AuthService) Register(ctx context.Context, input RegisterInput) (*models.Token, error) {
	if input.Email == "" || input.Password == "" || input.Name == "" {
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

	token, err := s.TokenModel.GenerateJWT(user.ID, TokenDuration)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (*models.Token, error) {
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

	token, err := s.TokenModel.GenerateJWT(user.ID, TokenDuration)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *AuthService) Logout(ctx context.Context) error {
	// TODO: implement refresh token invalidation or blacklist
	return nil
}
