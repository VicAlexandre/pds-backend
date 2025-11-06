package services

import (
	"context"
	"fmt"

	"github.com/VicAlexandre/pds-backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserModel  *models.UserModel
	TokenModel *models.JWTModel
}

func NewUserService(userModel *models.UserModel) *UserService {
	return &UserService{
		UserModel:  userModel,
		TokenModel: &models.JWTModel{},
	}
}

func (s *UserService) GetUserByID(ctx context.Context, token string) (*models.User, error) {
	claims, err := s.TokenModel.ParseJWT(token)
	if err != nil {
		return nil, err
	}

	userData, err := s.UserModel.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	return userData, nil
}

type ChangePasswordInput struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

func (s *UserService) ChangePassword(ctx context.Context, token string, input ChangePasswordInput) error {
	claims, err := s.TokenModel.ParseJWT(token)
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	// Buscar usuário atual para verificar senha
	user, err := s.UserModel.FindByID(ctx, claims.UserID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Verificar senha atual
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.CurrentPassword))
	if err != nil {
		return fmt.Errorf("current password is incorrect")
	}

	// Hash da nova senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// Atualizar senha
	err = s.UserModel.UpdatePassword(ctx, claims.UserID, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

func (s *UserService) DeleteAccount(ctx context.Context, token string) error {
	claims, err := s.TokenModel.ParseJWT(token)
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	err = s.UserModel.Delete(ctx, claims.UserID)
	if err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}

	return nil
}
