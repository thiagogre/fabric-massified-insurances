package adapters

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain/mocks"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/dto"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/utils"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/tests"
)

func TestGetAll_GetEvents_Success(t *testing.T) {
	tests.SetupLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEventService := mocks.NewMockEventInterface(ctrl)
	handler := NewEventHandler(mockEventService)

	events := []domain.Event{
		{BlockNumber: 16, TransactionID: "568bfb40797282fceaf8e5e1dc1dd073ca471f3e97aba1a0cf0706e88fd55bff", ChaincodeName: "basic", EventName: "DeleteAsset", Payload: "eyJDbGFpbVN0YXR1cyI6IkFjdGl2ZSIsIkNvdmVyYWdlQW1vdW50Ijo1MDAwLCJJRCI6InBvbGljeTEiLCJJbnN1cmVkSXRlbSI6IlNtYXJ0cGhvbmUgQURCIiwiT3duZXIiOiJEb25vIiwiUHJlbWl1bSI6MzAwLCJUZXJtIjoxMn0="},
		{BlockNumber: 19, TransactionID: "62c0afd66310a59bb98520d977f7d06fd9e834f095202fdf9b85a2a56c775338", ChaincodeName: "basic", EventName: "CreateAsset", Payload: "eyJDbGFpbVN0YXR1cyI6IkFjdGl2ZSIsIkNvdmVyYWdlQW1vdW50Ijo1MDAwLCJJRCI6InBvbGljeTEiLCJJbnN1cmVkSXRlbSI6IlNtYXJ0cGhvbmUgQUJDIiwiT3duZXIiOiJEb25vIiwiUHJlbWl1bSI6MzAwLCJUZXJtIjoxMn0="},
	}

	mockEventService.EXPECT().
		GetEventsFromStorage().
		Return(events, nil)

	req := httptest.NewRequest(http.MethodGet, "/event", nil)
	rec := httptest.NewRecorder()

	handler.GetAll(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	expected := dto.DocsResponse[domain.Event]{Docs: events}
	utils.AssertJSONResponse(t, rec.Body.String(), expected)
}

func TestGetAll_GetEvents_Fail(t *testing.T) {
	tests.SetupLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEventService := mocks.NewMockEventInterface(ctrl)
	handler := NewEventHandler(mockEventService)

	mockEventService.EXPECT().
		GetEventsFromStorage().
		Return(nil, errors.New("error retrieving events"))

	req := httptest.NewRequest(http.MethodGet, "/event", nil)
	rec := httptest.NewRecorder()

	handler.GetAll(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	require.Contains(t, rec.Body.String(), "Failed to retrieve events")
}
