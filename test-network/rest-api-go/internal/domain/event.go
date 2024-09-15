package domain

import "context"

type Event struct {
	BlockNumber   uint64
	TransactionID string
	ChaincodeName string
	EventName     string
	Payload       string
}

type EventInterface interface {
	ReplayEvents(ctx context.Context, channelID, chaincodeID string, txnBlockNumber int, txnID string) (<-chan *Events, error)
	HandleEvent(events <-chan *Events, txnID, filename string) error
	AppendEvent(event *Events, filename string) error
	GetEventsFromStorage() ([]Event, error)
}

type EventRepositoryInterface interface {
	AppendData(data []byte, filename string) error
	GetAllEvents() ([]Event, error)
}
