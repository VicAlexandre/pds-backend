package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/VicAlexandre/pds-backend/internal/services"
)

type ApostilasHandler struct {
	ApostilaService *services.ApostilaService
}

func extractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Println("Authorization header missing")
		return "", http.ErrNoCookie
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 {
		log.Printf("Invalid authorization header format: %s", authHeader)
		return "", fmt.Errorf("invalid authorization header format")
	}

	if strings.ToLower(parts[0]) != "bearer" {
		log.Printf("Invalid authorization scheme: %s", parts[0])
		return "", fmt.Errorf("invalid authorization scheme")
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		log.Println("Empty token")
		return "", fmt.Errorf("empty token")
	}

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

	w.WriteHeader(http.StatusOK)
}

func (h *ApostilasHandler) RenderApostilaPDF(w http.ResponseWriter, r *http.Request) {
	var input services.RenderPDFInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	pdf, err := h.ApostilaService.RenderApostilaPDF(r.Context(), input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")

	// Optional but recommended: Specify Content-Length for efficiency
	w.Header().Set("Content-Length", strconv.Itoa(len(pdf)))

	// Optional: Suggest a filename to the browser (for direct link access)
	w.Header().Set("Content-Disposition", "attachment; filename=\"apostila.pdf\"")

	w.Write(pdf)
}

func (h *ApostilasHandler) GetAllApostilas(w http.ResponseWriter, r *http.Request) {
	token, err := extractToken(r)
	if err != nil {
		http.Error(w, "authorization header missing", http.StatusUnauthorized)
		return
	}

	apostilas, err := h.ApostilaService.GetAllApostilas(r.Context(), token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apostilas)
}

func (h *ApostilasHandler) GetApostilaByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		// Tenta pegar do path parameter (Chi router)
		id = r.PathValue("id")
	}
	if id == "" {
		http.Error(w, "id query parameter or path parameter is required", http.StatusBadRequest)
		return
	}

	// Token é opcional para permitir compartilhamento público
	token, _ := extractToken(r)

	apostila, err := h.ApostilaService.GetApostilaByID(r.Context(), id, token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apostila)
}

func (h *ApostilasHandler) DeleteApostila(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		// Tenta pegar do path parameter (Chi router)
		id = r.PathValue("id")
	}
	if id == "" {
		http.Error(w, "id query parameter or path parameter is required", http.StatusBadRequest)
		return
	}

	token, err := extractToken(r)
	if err != nil {
		http.Error(w, "authorization header missing", http.StatusUnauthorized)
		return
	}

	err = h.ApostilaService.DeleteApostila(r.Context(), id, token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "apostila deleted successfully"})
}
