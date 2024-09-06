package handlers

import (
	"bufio"
	"encoding/json"
	"net/http"
	"os"

	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/constants"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/dto"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/models"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/logger"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/utils"
)

type EventHandler struct {
}

func InitEventHandler() *EventHandler {
	return &EventHandler{}
}

func (h *EventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger.Info("Received a request")

	file, err := os.Open(constants.EventLogFilename)
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
