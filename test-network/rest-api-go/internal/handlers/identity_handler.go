package handlers

import (
	"encoding/json"
	"net/http"
	"os/exec"

	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/constants"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/dto"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/logger"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/utils"
)

type IdentityHandler struct {
}

func InitIdentityHandler() *IdentityHandler {
	return &IdentityHandler{}
}

// Register and enroll an identity.
func (h *IdentityHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger.Info("Received a request")

	var body dto.IdentityRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Failed to parse request body")
		return
	}

	logger.Info(body)

	cmd := exec.Command("/bin/bash", "./registerEnrollIdentity.sh", body.Username)
	cmd.Dir = constants.TestNetworkPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Error executing script: " + err.Error())
		logger.Error("Script output: " + string(output))
		utils.ErrorResponse(w, http.StatusInternalServerError, "Error executing script")
		return
	}

	logger.Info("Script output: " + string(output))

	response := dto.SuccessResponse{Success: true}
	logger.Success(response)
	utils.SuccessResponse(w, http.StatusOK, response)
}
