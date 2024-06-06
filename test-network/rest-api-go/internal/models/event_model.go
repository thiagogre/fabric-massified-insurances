package models

type Event struct {
	BlockNumber   uint64
	TransactionID string
	ChaincodeName string
	EventName     string
	Payload       string
}
