package handlers

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
	"rest-api-go/internal/dto"
	"rest-api-go/internal/models"
)

const (
	filename = "events.log"
)

type EventHandler struct {
}

func InitEventHandler() *EventHandler {
	return &EventHandler{}
}

func (h *EventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	var events []models.Event

	// Read each line from the file
	for scanner.Scan() {
		var event models.Event

		// Parse the JSON data from the line
		if err := json.Unmarshal([]byte(scanner.Text()), &event); err != nil {
			http.Error(w, "Failed to parse file data", http.StatusInternalServerError)
			return
		}

		events = append(events, event)
	}

	// Marshal the events array into JSON
	docs := dto.DocsResponse[models.Event]{Docs: events}
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
