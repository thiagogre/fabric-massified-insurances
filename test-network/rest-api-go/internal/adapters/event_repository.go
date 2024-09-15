package adapters

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/constants"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/logger"
)

type EventRepository struct{}

func NewEventRepository() *EventRepository {
	return &EventRepository{}
}

func (repo *EventRepository) AppendData(data []byte, filename string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Error("Failed to open file: " + err.Error())
		return err
	}
	defer file.Close()

	data = append(data, '\n')

	if _, err := file.Write(data); err != nil {
		logger.Error("Failed to write data to file: " + err.Error())
		return err
	}

	return nil
}

func (r *EventRepository) GetAllEvents() ([]domain.Event, error) {
	file, err := os.Open(constants.EventLogFilename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var events []domain.Event
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		var event domain.Event

		if err := json.Unmarshal([]byte(scanner.Text()), &event); err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
