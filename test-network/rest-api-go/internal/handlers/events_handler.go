package handlers

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"
	"rest-api-go/internal/dto"
	"rest-api-go/internal/models"
	"rest-api-go/pkg/logger"
	"rest-api-go/pkg/utils"
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
	logger.Info("Received a request")

	file, err := os.Open(filename)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to open file")
		return
	}
	defer file.Close()

	var events []models.Event
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var event models.Event

		// Parse the JSON data from the line
		if err := json.Unmarshal([]byte(scanner.Text()), &event); err != nil {
			utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to parse file data")
			return
		}

		events = append(events, event)
	}

	response := dto.DocsResponse[models.Event]{Docs: events}
	logger.Success(response)
	utils.SuccessResponse(w, http.StatusOK, response)
}
