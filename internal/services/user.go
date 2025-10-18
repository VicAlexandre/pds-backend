package services

import (
	"context"

	"github.com/VicAlexandre/pds-backend/internal/models"
)

type UserService struct {
	UserModel  *models.UserModel
	TokenModel *models.JWTModel
}

func NewUserService(userModel *models.UserModel) *UserService {
	return &UserService{
		UserModel: userModel,
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
