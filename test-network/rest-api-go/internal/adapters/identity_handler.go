package adapters

import (
	"net/http"

	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/logger"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/utils"
)

type IdentityHandler struct {
	IdentityService domain.IdentityInterface
}

func NewIdentityHandler(identityService domain.IdentityInterface) *IdentityHandler {
	return &IdentityHandler{IdentityService: identityService}
}

// Register and enroll an identity.
func (h *IdentityHandler) Execute(w http.ResponseWriter, r *http.Request) {
	logger.Info("Received a request")

	credentials, err := h.IdentityService.Create()
	if err != nil {
		logger.Error("Error creating new random credentials: " + err.Error())
		utils.ErrorResponse(w, http.StatusInternalServerError, "Error creating new random credentials: "+err.Error())
		return
	}

	response := domain.SuccessResponse[domain.IdentityResponse]{Success: true, Data: domain.IdentityResponse{Username: credentials.Username, Password: credentials.Password}}
	logger.Success(response)
	utils.SuccessResponse(w, http.StatusOK, response)
}
