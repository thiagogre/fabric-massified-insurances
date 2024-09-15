package adapters

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain/mocks"
)

func TestEventRepository_AppendData_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := []byte(`{"TransactionID":"txn123","Payload":"test payload"}`)
	filename := "testfile.log"

	mockRepo := mocks.NewMockEventRepositoryInterface(ctrl)
	mockRepo.EXPECT().AppendData(data, filename).Return(nil)

	err := mockRepo.AppendData(data, filename)
	require.NoError(t, err)
}

func TestEventRepository_AppendData_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	data := []byte(`{"TransactionID":"txn123","Payload":"test payload"}`)
	filename := "/invalid/path/to/file"

	mockRepo := mocks.NewMockEventRepositoryInterface(ctrl)
	mockRepo.EXPECT().AppendData(data, filename).Return(errors.New("err"))

	err := mockRepo.AppendData(data, filename)
	require.Error(t, err)
}

func TestEventRepository_GetAllEvents_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEventRepositoryInterface(ctrl)

	events := []domain.Event{
		{TransactionID: "txn123", Payload: "test payload"},
		{TransactionID: "txn456", Payload: "another test payload"},
	}

	mockRepo.EXPECT().GetAllEvents().Return(events, nil)

	returnedEvents, err := mockRepo.GetAllEvents()
	require.NoError(t, err)
	require.Len(t, returnedEvents, len(events))

	for i, event := range events {
		require.Equal(t, event.TransactionID, returnedEvents[i].TransactionID)
		require.Equal(t, event.Payload, returnedEvents[i].Payload)
	}
}

func TestEventRepository_GetAllEvents_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockEventRepositoryInterface(ctrl)
	mockRepo.EXPECT().GetAllEvents().Return(nil, errors.New("err"))

	_, err := mockRepo.GetAllEvents()
	require.Error(t, err)
}
