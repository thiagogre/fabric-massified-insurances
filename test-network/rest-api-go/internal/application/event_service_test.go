package application

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain/mocks"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/tests"
)

func TestEventService_AppendEvent_Success(t *testing.T) {
	tests.SetupLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEventRepo := mocks.NewMockEventRepositoryInterface(ctrl)
	eventService := NewEventService(nil, mockEventRepo)

	event := &domain.Events{TransactionID: "txn123", Payload: []byte(`{"message":"event occurred"}`)}

	mockEventRepo.EXPECT().AppendData(gomock.Any(), "testfile").Return(nil).Times(1)

	err := eventService.AppendEvent(event, "testfile")
	require.NoError(t, err)
}

func TestEventService_AppendEvent_Failure(t *testing.T) {
	tests.SetupLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEventRepo := mocks.NewMockEventRepositoryInterface(ctrl)
	eventService := NewEventService(nil, mockEventRepo)

	event := &domain.Events{TransactionID: "txn123", Payload: []byte(`{"message":"event occurred"}`)}

	mockEventRepo.EXPECT().AppendData(gomock.Any(), "testfile").Return(errors.New("failed to append data")).Times(1)

	err := eventService.AppendEvent(event, "testfile")
	require.Error(t, err)
}

func TestEventService_GetEventsFromStorage_Success(t *testing.T) {
	tests.SetupLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEventRepo := mocks.NewMockEventRepositoryInterface(ctrl)
	eventService := NewEventService(nil, mockEventRepo)

	mockEventRepo.EXPECT().GetAllEvents().Return([]domain.Event{{TransactionID: "txn123"}}, nil).Times(1)

	events, err := eventService.GetEventsFromStorage()
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.Equal(t, "txn123", events[0].TransactionID)
}

func TestEventService_GetEventsFromStorage_Failure(t *testing.T) {
	tests.SetupLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEventRepo := mocks.NewMockEventRepositoryInterface(ctrl)
	eventService := NewEventService(nil, mockEventRepo)

	mockEventRepo.EXPECT().GetAllEvents().Return(nil, errors.New("storage error")).Times(1)

	events, err := eventService.GetEventsFromStorage()
	require.Error(t, err)
	require.Nil(t, events)
}
