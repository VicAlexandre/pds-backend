package services

import (
	"context"
	"fmt"
	"log"

	"github.com/VicAlexandre/pds-backend/internal/models"
	"github.com/google/uuid"
)

type AddApostilaInput struct {
	Id string `json:"data"`
}

// EditedApostilaInput receives a 'data' field with the subfield 'id' and 'file' (html content)
type EditedApostilaInput struct {
	Data struct {
		Id   string `json:"id"`
		Html string `json:"file"`
	} `json:"data"`
}

type ApostilaService struct {
	ApostilaModel *models.ApostilaModel
	UserModel     *models.UserModel
	TokenModel    *models.JWTModel
}

func NewApostilaService(apostilaModel *models.ApostilaModel, userModel *models.UserModel, tokenModel *models.JWTModel) *ApostilaService {
	return &ApostilaService{
		ApostilaModel: apostilaModel,
		UserModel:     userModel,
		TokenModel:    tokenModel,
	}
}

func (s *ApostilaService) AddApostila(ctx context.Context, input AddApostilaInput, token string) (*models.Apostila, error) {
	claims, err := s.TokenModel.ParseJWT(token)
	if err != nil {
		log.Println("Error parsing JWT: ", err)
		return nil, err
	}

	u, err := uuid.Parse(input.Id)
	if err != nil {
		fmt.Printf("Error parsing UUID: %v\n", err)
		return nil, err
	}

	apostila, err := s.ApostilaModel.Insert(ctx, u, claims.UserID)
	if err != nil {
		log.Println("Error inserting apostila: ", err)
		return nil, err
	}

	log.Println("Generated:", apostila)

	return apostila, nil
}

func (s *ApostilaService) GetEditedApostilaHTML(ctx context.Context, id string, token string) (*models.EditedApostilaHTML, error) {
	claims, err := s.TokenModel.ParseJWT(token)
	if err != nil {
		log.Println("Error parsing JWT: ", err)
		return nil, err
	}

	u, err := uuid.Parse(id)
	if err != nil {
		fmt.Printf("Error parsing UUID: %v\n", err)
		return nil, err
	}

	htmlContent, err := s.ApostilaModel.GetEditedHTMLByID(ctx, u, claims.UserID)
	if err != nil {
		log.Println("Error getting edited HTML: ", err)
		return nil, err
	}

	log.Println("Retrieved edited HTML for apostila ID:", u)

	return htmlContent, nil
}

func (s *ApostilaService) EditApostila(ctx context.Context, input EditedApostilaInput, token string) error {
	claims, err := s.TokenModel.ParseJWT(token)
	if err != nil {
		log.Println("Error parsing JWT: ", err)
		return err
	}

	u, err := uuid.Parse(input.Data.Id)
	if err != nil {
		fmt.Printf("Error parsing UUID: %v\n", err)
		fmt.Println("Input ID was: ", input.Data.Id)
		return err
	}

	return s.ApostilaModel.UpdateEditedHTMLByID(ctx, u, input.Data.Html, claims.UserID)
}
