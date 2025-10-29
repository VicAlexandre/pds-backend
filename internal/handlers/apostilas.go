package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/VicAlexandre/pds-backend/internal/services"
)

type ApostilasHandler struct {
	ApostilaService *services.ApostilaService
}

func extractToken(r *http.Request) (string, error) {
	ctx := r.Context()

	req := r.WithContext(ctx)

	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		log.Println("Authorization header missing: ", authHeader)
		return "", http.ErrNoCookie
	}

	parts := strings.SplitN(authHeader, " ", 2)

	token := parts[1]

	return token, nil
}

/* receives the id of a new apostila, authenticated user via jwt token and prints the id and jwt data */
func (h *ApostilasHandler) AddApostila(w http.ResponseWriter, r *http.Request) {
	var input services.AddApostilaInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	token, err := extractToken(r)
	if err != nil {
		http.Error(w, "authorization header missing", http.StatusUnauthorized)
	}

	apostila, err := h.ApostilaService.AddApostila(r.Context(), input, token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(apostila)
}

func (h *ApostilasHandler) GetEditedApostilaHTML(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id query parameter is required", http.StatusBadRequest)
		return
	}

	token, err := extractToken(r)
	if err != nil {
		http.Error(w, "authorization header missing", http.StatusUnauthorized)
		return
	}

	htmlContent, err := h.ApostilaService.GetEditedApostilaHTML(r.Context(), id, token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(htmlContent)
}

func (h *ApostilasHandler) EditApostila(w http.ResponseWriter, r *http.Request) {
	var input services.EditedApostilaInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	token, err := extractToken(r)
	if err != nil {
		http.Error(w, "authorization header missing", http.StatusUnauthorized)
		return
	}
	h.ApostilaService.EditApostila(r.Context(), input, token)
}
