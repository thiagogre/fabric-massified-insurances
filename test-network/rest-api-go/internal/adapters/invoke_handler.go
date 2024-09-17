package adapters

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/constants"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/dto"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/logger"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/utils"
)

type InvokeHandler struct {
	InvokeService domain.InvokeInterface
	EventService  domain.EventInterface
}

func NewInvokeHandler(invokeService domain.InvokeInterface, eventService domain.EventInterface) *InvokeHandler {
	return &InvokeHandler{InvokeService: invokeService, EventService: eventService}
}

func (h *InvokeHandler) Execute(w http.ResponseWriter, r *http.Request) {
	logger.Info("Received a request")

	var body dto.InvokeRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		logger.Error("Failed to parse request body" + err.Error())
		utils.ErrorResponse(w, http.StatusBadRequest, "Failed to parse request body")
		return
	}
	logger.Info(body)

	data, err := h.InvokeService.ExecuteInvoke(body.ChannelID, body.ChaincodeID, body.Function, body.Args)
	if err != nil {
		logger.Error("Error executing invoke: " + err.Error())
		utils.ErrorResponse(w, http.StatusInternalServerError, "Error executing invoke: "+err.Error())
		return
	}

	logger.Info("*** Start chaincode event replay ***")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	events, err := h.EventService.ReplayEvents(ctx, body.ChannelID, body.ChaincodeID, int(data.BlockNumber), data.TransactionID)
	if err != nil {
		logger.Error("Error replaying event: " + err.Error())
		utils.ErrorResponse(w, http.StatusInternalServerError, "Error replaying event: "+err.Error())
		return
	}

	if err := h.EventService.HandleEvent(events, data.TransactionID, constants.EventLogFilename); err != nil {
		logger.Error("Error handling event: " + err.Error())
		utils.ErrorResponse(w, http.StatusInternalServerError, "Error handling event: "+err.Error())
		return
	}
	logger.Info("*** Finish chaincode event replay ***")

	response := dto.InvokeSuccessResponse{Success: true, Data: data}
	logger.Success(data)
	utils.SuccessResponse(w, http.StatusOK, response)
}
