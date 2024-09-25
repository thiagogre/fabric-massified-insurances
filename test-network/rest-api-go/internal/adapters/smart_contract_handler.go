package adapters

import (
	"net/http"

	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/logger"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/utils"
)

type SmartContractHandler struct {
}

func NewSmartContractHandler() *SmartContractHandler {
	return &SmartContractHandler{}
}

func (h *SmartContractHandler) Info(w http.ResponseWriter, r *http.Request) {
	logger.Info("Received a request")

	routes := []string{
		"/smartcontract/query",
		"/smartcontract/invoke",
	}

	response := domain.SuccessResponse[[]string]{Success: true, Data: routes}
	logger.Success(response)
	utils.SuccessResponse(w, http.StatusOK, response)
}
