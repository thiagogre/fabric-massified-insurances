package adapters

import (
	"net/http"

	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/dto"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/logger"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/utils"
)

type QueryHandler struct {
	QueryService domain.QueryInterface
}

func NewQueryHandler(queryService domain.QueryInterface) *QueryHandler {
	return &QueryHandler{QueryService: queryService}
}

func (h *QueryHandler) Execute(w http.ResponseWriter, r *http.Request) {
	logger.Info("Received a request")

	queryParams := r.URL.Query()
	chainCodeName := queryParams.Get("chaincodeid")
	channelID := queryParams.Get("channelid")
	function := queryParams.Get("function")
	args := r.URL.Query()["args"]

	logger.Info(queryParams)

	data, err := h.QueryService.ExecuteQuery(channelID, chainCodeName, function, args)
	if err != nil {
		logger.Error("Error executing query: " + err.Error())
		utils.ErrorResponse(w, http.StatusInternalServerError, "Error executing query: "+err.Error())
		return
	}

	response := dto.QuerySuccessResponse{Success: true, Data: data}
	logger.Success(response)
	utils.SuccessResponse(w, http.StatusOK, response)
}
