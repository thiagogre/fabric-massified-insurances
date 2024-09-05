package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/logger"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/org"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/utils"
)

type QueryHandler struct {
	OrgSetup org.OrgSetup
}

func InitQueryHandler(orgSetup org.OrgSetup) *QueryHandler {
	return &QueryHandler{OrgSetup: orgSetup}
}

func (h *QueryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger.Info("Received a request")

	queryParams := r.URL.Query()
	chainCodeName := queryParams.Get("chaincodeid")
	channelID := queryParams.Get("channelid")
	function := queryParams.Get("function")
	args := r.URL.Query()["args"]

	logger.Info(queryParams)

	network := h.OrgSetup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	txn, err := contract.EvaluateTransaction(function, args...)
	if err != nil {
		logger.Error("Error evaluating transaction " + err.Error())
		utils.ErrorResponse(w, http.StatusInternalServerError, "Error evaluating transaction "+err.Error())
		return
	}

	response, err := parseEvaluatedTransaction(txn)
	if err != nil {
		logger.Error("Error parsing transaction " + err.Error())
		utils.ErrorResponse(w, http.StatusInternalServerError, "Error parsing transaction "+err.Error())
		return
	}

	logger.Success(response)
	utils.SuccessResponse(w, http.StatusOK, response)
}

// Decode base64-encoded protobuf binary data into JSON
func parseEvaluatedTransaction(b []byte) (interface{}, error) {
	var response interface{}
	err := json.Unmarshal(b, &response)
	if err != nil {
		return make([]byte, 0), err
	}

	return response, nil
}
