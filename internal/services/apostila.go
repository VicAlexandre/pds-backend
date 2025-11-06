package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/VicAlexandre/pds-backend/internal/models"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
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

type RenderPDFInput struct {
	Data struct {
		Html string `json:"html"`
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

const cleanupScript = `
(function() {
    document.querySelectorAll('h2[role="button"]').forEach(h2 => {
        h2.setAttribute("aria-expanded", "true");
        let next = h2.nextElementSibling;
        if (next && next.classList.contains('ouvir')) {
            next = next.nextElementSibling;
        }
        if (next && next.classList.contains('content')) {
            next.removeAttribute('hidden');
        }
    });

    const controls = document.querySelector('.controls');
    if (controls) {
        controls.remove();
    }
    document.querySelectorAll('.ouvir').forEach(button => {
        button.remove();
    });
    
    document.querySelectorAll('.toggle-icon').forEach(icon => {
        icon.textContent = " "; 
    });

    document.querySelectorAll('details.spoiler').forEach(details => {
        const content = details.innerHTML;
        const newDiv = document.createElement('div');
        newDiv.innerHTML = content;
        details.parentNode.replaceChild(newDiv, details);
    });

    const scriptTag = document.querySelector('script');
    if (scriptTag) {
        scriptTag.remove();
    }

    document.body.classList.remove('dark');
})();
`

func (s *ApostilaService) RenderApostilaPDF(ctx context.Context, input RenderPDFInput) ([]byte, error) {
	cctx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	htmlB64 := base64.StdEncoding.EncodeToString([]byte(input.Data.Html))
	dataURL := fmt.Sprintf("data:text/html;base64,%s", htmlB64)

	var pdfBuf []byte
	var bodyContent string

	err := chromedp.Run(cctx,
		chromedp.Navigate(dataURL),

		chromedp.WaitReady("body", chromedp.ByQuery),

		chromedp.ActionFunc(func(ctx context.Context) error {
			ctxErr := chromedp.Evaluate(cleanupScript, nil).Do(ctx)
			if ctxErr != nil {
				return fmt.Errorf("erro ao executar script de limpeza: %w", ctxErr)
			}

			return nil
		}),

		chromedp.Sleep(2*time.Second),

		chromedp.Text("body", &bodyContent, chromedp.ByQuery, chromedp.NodeVisible),

		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfBuf, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithPaperWidth(8.27).   // A4
				WithPaperHeight(11.69). // A4
				Do(ctx)
			return err
		}),
	)

	if err != nil {
		return nil, fmt.Errorf("erro ao renderizar PDF: %w", err)
	}

	if len(bodyContent) == 0 {
		return nil, fmt.Errorf("PDF vazio: o HTML não foi renderizado")
	}

	return pdfBuf, nil
}

func (s *ApostilaService) GetAllApostilas(ctx context.Context, token string) ([]*models.Apostila, error) {
	claims, err := s.TokenModel.ParseJWT(token)
	if err != nil {
		log.Println("Error parsing JWT: ", err)
		return nil, err
	}

	apostilas, err := s.ApostilaModel.GetAllByUserID(ctx, claims.UserID)
	if err != nil {
		log.Println("Error getting apostilas: ", err)
		return nil, err
	}

	log.Println("Retrieved", len(apostilas), "apostilas for user", claims.UserID)

	return apostilas, nil
}

func (s *ApostilaService) GetApostilaByID(ctx context.Context, id string, token string) (*models.Apostila, error) {
	u, err := uuid.Parse(id)
	if err != nil {
		fmt.Printf("Error parsing UUID: %v\n", err)
		return nil, err
	}

	// Se tiver token, valida e verifica se é do usuário
	if token != "" {
		claims, err := s.TokenModel.ParseJWT(token)
		if err != nil {
			log.Println("Error parsing JWT: ", err)
			return nil, err
		}

		apostila, err := s.ApostilaModel.GetByID(ctx, u)
		if err != nil {
			log.Println("Error getting apostila: ", err)
			return nil, err
		}

		// Verifica se a apostila pertence ao usuário
		if apostila.UserID != claims.UserID {
			return nil, fmt.Errorf("unauthorized access to apostila")
		}

		return apostila, nil
	}

	// Se não tiver token, busca a apostila (para compartilhamento público)
	apostila, err := s.ApostilaModel.GetByID(ctx, u)
	if err != nil {
		log.Println("Error getting apostila: ", err)
		return nil, err
	}

	return apostila, nil
}

func (s *ApostilaService) DeleteApostila(ctx context.Context, id string, token string) error {
	claims, err := s.TokenModel.ParseJWT(token)
	if err != nil {
		log.Println("Error parsing JWT: ", err)
		return err
	}

	u, err := uuid.Parse(id)
	if err != nil {
		fmt.Printf("Error parsing UUID: %v\n", err)
		return err
	}

	err = s.ApostilaModel.Delete(ctx, u, claims.UserID)
	if err != nil {
		log.Println("Error deleting apostila: ", err)
		return err
	}

	log.Println("Deleted apostila ID:", u, "for user", claims.UserID)

	return nil
}
