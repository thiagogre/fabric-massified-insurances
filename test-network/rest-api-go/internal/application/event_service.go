package application

import (
	"context"
	"encoding/json"

	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/logger"
)

type EventService struct {
	BlockchainGateway domain.BlockchainGatewayInterface
	EventRepository   domain.EventRepositoryInterface
}

func NewEventService(blockchainGateway domain.BlockchainGatewayInterface, eventRepository domain.EventRepositoryInterface) *EventService {
	return &EventService{EventRepository: eventRepository, BlockchainGateway: blockchainGateway}
}

func (s *EventService) ReplayEvents(ctx context.Context, channelID, chaincodeID string, txnBlockNumber int, txnID string) (<-chan *domain.Events, error) {
	event, err := s.BlockchainGateway.GetNetwork(channelID).GetEvents(ctx, chaincodeID, txnBlockNumber)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (s *EventService) HandleEvent(events <-chan *domain.Events, txnID, filename string) error {
	for event := range events {
		if event.TransactionID != txnID {
			continue
		}

		if err := s.AppendEvent(event, filename); err != nil {
			return err
		}

		logger.Info(string(event.Payload))
		break
	}

	return nil
}

func (s *EventService) AppendEvent(event *domain.Events, filename string) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return s.EventRepository.AppendData(eventBytes, filename)
}

func (s *EventService) GetEventsFromStorage() ([]domain.Event, error) {
	return s.EventRepository.GetAllEvents()
}
