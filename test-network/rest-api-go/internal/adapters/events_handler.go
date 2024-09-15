package adapters

import (
	"net/http"

	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/dto"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/logger"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/utils"
)

type EventHandler struct {
	EventService domain.EventInterface
}

func NewEventHandler(eventService domain.EventInterface) *EventHandler {
	return &EventHandler{EventService: eventService}
}

func (h *EventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger.Info("Received a request")

	events, err := h.EventService.GetEventsFromStorage()
	if err != nil {
		logger.Error("Failed to retrieve events: " + err.Error())
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve events")
		return
	}

	response := dto.DocsResponse[domain.Event]{Docs: events}
	logger.Success(response)
	utils.SuccessResponse(w, http.StatusOK, response)
}
