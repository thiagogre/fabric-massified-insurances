package handlers

import (
	"net/http"

	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/constants"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/dto"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/cmd"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/logger"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/utils"
)

type IdentityHandler struct {
	commandExecutor cmd.CommandExecutorInterface
}

func InitIdentityHandler(commandExecutor cmd.CommandExecutorInterface) *IdentityHandler {
	return &IdentityHandler{commandExecutor: commandExecutor}
}

// Register and enroll an identity.
func (h *IdentityHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger.Info("Received a request")

	randomUsername, randomPassword, err := utils.GenerateRandomCredentials(constants.DefaultUsernameLength, constants.DefaultPasswordLength)
	if err != nil {
		logger.Error("Error generating random credentials: " + err.Error())
		utils.ErrorResponse(w, http.StatusInternalServerError, "Error generating random credentials")
		return
	}

	output, err := h.commandExecutor.ExecuteCommand("/bin/bash", "./registerEnrollIdentity.sh", randomUsername, randomPassword)
	if err != nil {
		logger.Error("Error executing script: " + err.Error())
		logger.Error("Script output: " + string(output))
		utils.ErrorResponse(w, http.StatusInternalServerError, "Error executing script")
		return
	}

	logger.Info("Script output: " + string(output))

	response := dto.SuccessResponse[dto.IdentityResponse]{Success: true, Data: dto.IdentityResponse{Username: randomUsername, Password: randomPassword}}
	logger.Success(response)
	utils.SuccessResponse(w, http.StatusOK, response)
}
