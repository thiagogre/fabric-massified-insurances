package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rest-api-go/internal/dto"
	"rest-api-go/internal/repositories"
	"rest-api-go/internal/services"
	"rest-api-go/pkg/db"
)

type AuthHandler struct {
	AuthService *services.AuthService
}

func InitAuthHandler(db db.Database) *AuthHandler {
	userRepository := &repositories.SQLUserRepository{DB: db}
	authService := &services.AuthService{UserRepository: userRepository}
	return &AuthHandler{AuthService: authService}
}

func (h *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[AUTH] Received a request")

	var body dto.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	_, authenticated, err := h.AuthService.AuthenticateUser(body.Username, body.Password)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if !authenticated {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	successResponse := dto.SuccessResponse{Success: true}
	jsonResponse, err := json.Marshal(successResponse)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
