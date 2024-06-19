package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"rest-api-go/internal/dto"
	"rest-api-go/pkg/logger"
	"rest-api-go/pkg/org"
	"rest-api-go/pkg/utils"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

const (
	EventLogFilename = "events.log"
)

type InvokeHandler struct {
	OrgSetup org.OrgSetup
}

func InitInvokeHandler(orgSetup org.OrgSetup) *InvokeHandler {
	return &InvokeHandler{OrgSetup: orgSetup}
}

func (h *InvokeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger.Info("Received a request")

	var body dto.InvokeRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		logger.Error("Failed to parse request body" + err.Error())
		utils.ErrorResponse(w, http.StatusBadRequest, "Failed to parse request body")
		return
	}

	logger.Info(body)

	network := h.OrgSetup.Gateway.GetNetwork(body.ChannelID)
	contract := network.GetContract(body.ChaincodeID)
	txnProposal, err := contract.NewProposal(body.Function, client.WithArguments(body.Args...))
	if err != nil {
		logger.Error("Error creating txn proposal " + err.Error())
		utils.ErrorResponse(w, http.StatusInternalServerError, "Error creating txn proposal "+err.Error())
		return
	}

	txnEndorsed, err := txnProposal.Endorse()
	if err != nil {
		logger.Error("Error endorsing txn " + err.Error())
		utils.ErrorResponse(w, http.StatusInternalServerError, "Error endorsing txn "+err.Error())
		return
	}

	txnCommitted, err := txnEndorsed.Submit()
	if err != nil {
		logger.Error("Error submitting transaction " + err.Error())
		utils.ErrorResponse(w, http.StatusInternalServerError, "Error submitting transaction "+err.Error())
		return
	}

	status, err := txnCommitted.Status()
	if err != nil {
		logger.Error("Failed to get transaction commit status " + err.Error())
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to get transaction commit status "+err.Error())
		return
	}

	if !status.Successful {
		logger.Error("Failed to commit transaction with status code " + string(status.Code))
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to commit transaction with status code "+string(status.Code))
		return
	}

	event := NewEvent(network, body.ChaincodeID, status.BlockNumber, status.TransactionID)
	go event.Replay()

	logger.Success(status)
	utils.SuccessResponse(w, http.StatusOK, status)
}

type Event struct {
	network        *client.Network
	chaincodeID    string
	txnBlockNumber uint64
	txnID          string
}

func NewEvent(network *client.Network, chaincodeID string, txnBlockNumber uint64, txnID string) *Event {
	return &Event{
		network:        network,
		chaincodeID:    chaincodeID,
		txnBlockNumber: txnBlockNumber,
		txnID:          txnID,
	}
}

func (e *Event) Replay() {
	logger.Info("*** Start chaincode event replay ***")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	events, err := e.network.ChaincodeEvents(ctx, e.chaincodeID, client.WithStartBlock(e.txnBlockNumber))
	if err != nil {
		logger.Error("Failed to start chaincode event listening " + err.Error())
		return
	}

	for event := range events {
		if event.TransactionID != e.txnID {
			continue
		}

		if err := e.Append(event, EventLogFilename); err != nil {
			logger.Error("Error appending event to file " + err.Error())
			return
		}

		logger.Info(string(event.Payload))
		break
	}

	logger.Info("*** Finish chaincode event replay ***")
}

func (e *Event) Append(event *client.ChaincodeEvent, filename string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Error("Failed to open file " + err.Error())
		return err
	}
	defer file.Close()

	eventBytes, err := json.Marshal(event)
	if err != nil {
		logger.Error("Failed to marshal event to JSON " + err.Error())
		return err
	}
	eventBytes = append(eventBytes, '\n')

	if _, err := file.Write(eventBytes); err != nil {
		logger.Error("Failed to write event to file " + err.Error())
		return err
	}

	return nil
}
