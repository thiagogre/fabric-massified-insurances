package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/dto"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/repositories"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/services"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/db"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/logger"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/utils"
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
	logger.Info("Received a request")

	var body dto.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		logger.Error("Failed to parse request body" + err.Error())
		utils.ErrorResponse(w, http.StatusBadRequest, "Failed to parse request body")
		return
	}

	logger.Info(body)

	if user, err := h.AuthService.AuthenticateUser(body.Username, body.Password); err != nil {
		logger.Error("User not found" + err.Error())
		utils.ErrorResponse(w, http.StatusNotFound, "User not found")
		return
	} else if user == nil {
		logger.Error("Incorrect password")
		utils.ErrorResponse(w, http.StatusUnauthorized, "Incorrect password")
		return
	}

	response := dto.SuccessResponse[dto.AuthRequest]{Success: true, Data: body}
	logger.Success(response)
	utils.SuccessResponse(w, http.StatusOK, response)
}
