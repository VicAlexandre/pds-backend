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
    log.Printf("Starting PDF generation for HTML length: %d", len(input.Data.Html))

    chromePaths := []string{
        "/usr/bin/google-chrome",
        "/usr/bin/chromium-browser", 
        "/usr/bin/chromium",
    }
    
    var chromePath string
    for _, path := range chromePaths {
        if _, err := os.Stat(path); err == nil {
            chromePath = path
            log.Printf("Found Chrome at: %s", path)
            break
        }
    }
    
    if chromePath == "" {
        return nil, fmt.Errorf("no Chrome executable found")
    }

    opts := []chromedp.ExecAllocatorOption{
        chromedp.ExecPath(chromePath),
        chromedp.Headless,
        chromedp.NoSandbox,
        chromedp.DisableGPU,
    }

    allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
    defer cancel()

    browserCtx, cancel := chromedp.NewContext(allocCtx)
    defer cancel()

    timeoutCtx, cancel := context.WithTimeout(browserCtx, 45*time.Second)
    defer cancel()

    var pdfBuf []byte

    htmlB64 := base64.StdEncoding.EncodeToString([]byte(input.Data.Html))
    dataURL := fmt.Sprintf("data:text/html;base64,%s", htmlB64)

    log.Printf("Navigating to data URL and generating PDF...")

    err := chromedp.Run(timeoutCtx,
        chromedp.Navigate(dataURL),
        chromedp.WaitReady("body", chromedp.ByQuery),
        chromedp.Sleep(1*time.Second),
        chromedp.ActionFunc(func(ctx context.Context) error {
            var err error
            pdfBuf, _, err = page.PrintToPDF().
                WithPrintBackground(false),
                WithPaperWidth(8.27).
                WithPaperHeight(11.69).
                WithMarginTop(0.3).
                WithMarginBottom(0.3).
                WithMarginLeft(0.3).
                WithMarginRight(0.3).
                Do(ctx)
            if err != nil {
                return fmt.Errorf("page.PrintToPDF failed: %w", err)
            }
            return nil
        }),
    )

    if err != nil {
        log.Printf("ChromeDP execution failed: %v", err)
        return nil, fmt.Errorf("failed to generate PDF: %w", err)
    }

    if len(pdfBuf) == 0 {
        return nil, fmt.Errorf("generated PDF is empty")
    }

    log.Printf("PDF generated successfully: %d bytes", len(pdfBuf))
    return pdfBuf, nil
}