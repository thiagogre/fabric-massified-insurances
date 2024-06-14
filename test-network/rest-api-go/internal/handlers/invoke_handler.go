package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"rest-api-go/internal/dto"
	"rest-api-go/pkg/logger"
	"rest-api-go/pkg/org"
	"rest-api-go/pkg/utils"

	"github.com/hyperledger/fabric-gateway/pkg/client"
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

	logger.Success(status)
	utils.SuccessResponse(w, http.StatusOK, status)

	// start a new goroutine
	go replayChaincodeEvents(h.OrgSetup.Context, network, body.ChaincodeID, status.BlockNumber)
}

// Replay events from the block containing the first transaction
func replayChaincodeEvents(ctx context.Context, network *client.Network, chaincodeID string, startBlock uint64) {
	logger.Info("*** Start chaincode event replay ***")

	events, err := network.ChaincodeEvents(ctx, chaincodeID, client.WithStartBlock(0))
	if err != nil {
		logger.Error("Failed to start chaincode event listening " + err.Error())
		return
	}

	timeout := time.After(30 * time.Second) // Set a timeout for event replay
	var eventsBuffer bytes.Buffer

	for {
		select {
		case <-timeout:
			logger.Info("*** Event replay timeout reached ***")

			// Write buffer content to the file
			if err := writeBufferToFile(&eventsBuffer, "events.log"); err != nil {
				logger.Error("Error writing buffer to file " + err.Error())
				return
			}

		case event, ok := <-events:
			if !ok {
				logger.Info("*** Event channel closed ***")
				return
			}

			writeEventToBuffer(event, &eventsBuffer)

			logger.Info(string(event.Payload))

			// // Example condition to stop listening; modify as needed.
			// if event.EventName == "DeleteAsset" {
			// 	return
			// }

		case <-ctx.Done():
			logger.Info("*** Context canceled, stopping event replay ***")
			return
		}
	}
}

// Write the chaincode event to a buffer
func writeEventToBuffer(event *client.ChaincodeEvent, buffer *bytes.Buffer) error {
	// Marshal the event into JSON
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// Write the event bytes to the buffer
	_, err = buffer.Write(eventBytes)
	if err != nil {
		return err
	}

	// Add a newline separator between events
	_, err = buffer.WriteString("\n")
	if err != nil {
		return err
	}

	return nil
}

// Write the buffer content to a new file
func writeBufferToFile(buffer *bytes.Buffer, filename string) error {
	// Create a new file or truncate an existing file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the buffer content to the file
	_, err = buffer.WriteTo(file)
	if err != nil {
		return err
	}

	return nil
}
