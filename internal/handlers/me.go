package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/VicAlexandre/pds-backend/internal/services"
)

type MeHandler struct {
	UserService *services.UserService
}

func NewMeHandler(userService *services.UserService) *MeHandler {
	return &MeHandler{
		UserService: userService,
	}
}

func (h *MeHandler) FetchUserData(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "missing authorization header", http.StatusUnauthorized)
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
		return
	}

	token := parts[1]

	user, err := h.UserService.GetUserByID(r.Context(), token)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
