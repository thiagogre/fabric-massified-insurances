package web

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// InvokeRequest represents the structure of the JSON object containing invocation request data.
type InvokeRequest struct {
	ChaincodeID string   `json:"chaincodeid"`
	ChannelID   string   `json:"channelid"`
	Function    string   `json:"function"`
	Args        []string `json:"args"`
}

// InvokeResponse represents the structure of the JSON object to be returned in the response.
type InvokeResponse struct {
	TransactionID string `json:"transaction_id"`
	Result        string `json:"result"`
}

// Invoke handles chaincode invoke requests.
func (setup *OrgSetup) Invoke(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Invoke request")

	// Decode the JSON object
	var invokeReq InvokeRequest
	if err := json.NewDecoder(r.Body).Decode(&invokeReq); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %v\n", invokeReq.ChannelID, invokeReq.ChaincodeID, invokeReq.Function, invokeReq.Args)

	network := setup.Gateway.GetNetwork(invokeReq.ChannelID)
	contract := network.GetContract(invokeReq.ChaincodeID)

	txnProposal, err := contract.NewProposal(invokeReq.Function, client.WithArguments(invokeReq.Args...))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating txn proposal: %s", err), http.StatusInternalServerError)
		return
	}
	txnEndorsed, err := txnProposal.Endorse()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error endorsing txn: %s", err), http.StatusInternalServerError)
		return
	}
	txnCommitted, err := txnEndorsed.Submit()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error submitting transaction: %s", err), http.StatusInternalServerError)
		return
	}

	status, err := txnCommitted.Status()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get transaction commit status: %s", err), http.StatusInternalServerError)
		return
	}

	if !status.Successful {
		http.Error(w, fmt.Sprintf("Failed to commit transaction with status code %v", status.Code), http.StatusInternalServerError)
		return
	}

	// Response struct
	response := InvokeResponse{
		TransactionID: txnCommitted.TransactionID(),
		Result:        string(txnEndorsed.Result()), // Converting byte slice to string
	}

	// Convert response to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

	// start a new goroutine
	go replayChaincodeEvents(setup.Context, network, invokeReq.ChaincodeID, status.BlockNumber)
}

// Replay events from the block containing the first transaction
func replayChaincodeEvents(ctx context.Context, network *client.Network, chaincodeID string, startBlock uint64) {
	fmt.Println("\n*** Start chaincode event replay")

	events, err := network.ChaincodeEvents(ctx, chaincodeID, client.WithStartBlock(0))
	if err != nil {
		fmt.Printf("Failed to start chaincode event listening: %v\n", err)
		return
	}

	timeout := time.After(30 * time.Second) // Set a timeout for event replay
	var eventsBuffer bytes.Buffer

	for {
		select {
		case <-timeout:
			fmt.Println("Event replay timeout reached")
			// Write buffer content to the file
			if err := writeBufferToFile(&eventsBuffer, "events.log"); err != nil {
				fmt.Println("Error writing buffer to file:", err)
				return
			}

		case event, ok := <-events:
			if !ok {
				fmt.Println("Event channel closed")
				return
			}

			// writeEventToFile(event, "events.json")
			writeEventToBuffer(event, &eventsBuffer)

			fmt.Printf("\n<-- Replayed event payload %s\n", event.Payload)

			// // Example condition to stop listening; modify as needed.
			// if event.EventName == "DeleteAsset" {
			// 	fmt.Println("Stopping event replay on 'DeleteAsset' event")
			// 	return
			// }

		case <-ctx.Done():
			fmt.Println("Context canceled, stopping event replay")
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
