package web

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
)

type Event struct {
	BlockNumber   uint64 `json:"BlockNumber"`
	TransactionID string `json:"TransactionID"`
	ChaincodeName string `json:"ChaincodeName"`
	EventName     string `json:"EventName"`
	Payload       string `json:"Payload"`
}

type Docs struct {
	Docs []Event `json:"docs"`
}

func Events(w http.ResponseWriter, r *http.Request) {
	filename := "events.log"

	// Open the file
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	var events []Event

	// Read each line from the file
	for scanner.Scan() {
		var event Event

		// Parse the JSON data from the line
		if err := json.Unmarshal([]byte(scanner.Text()), &event); err != nil {
			http.Error(w, "Failed to parse file data", http.StatusInternalServerError)
			return
		}

		events = append(events, event)
	}

	// Marshal the events array into JSON
	docs := Docs{Docs: events}
	responseJSON, err := json.Marshal(docs)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	// Set the response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Write the response JSON to the client
	w.Write(responseJSON)
}
