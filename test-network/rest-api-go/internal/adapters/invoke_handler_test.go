package adapters_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/constants"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/adapters"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain/mocks"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/dto"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/tests"
)

func TestInvokeHandler_ServeHTTP(t *testing.T) {
	tests.SetupLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInvokeService := mocks.NewMockInvokeInterface(ctrl)
	mockEventService := mocks.NewMockEventInterface(ctrl)

	handler := adapters.NewInvokeHandler(mockInvokeService, mockEventService)

	t.Run("should return 400 when failed to decode request body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/invoke", bytes.NewBufferString("invalid-body"))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Contains(t, w.Body.String(), "Failed to parse request body")
	})

	t.Run("should return 500 when ExecuteInvoke returns an error", func(t *testing.T) {
		body := dto.InvokeRequest{
			ChannelID:   "test-channel",
			ChaincodeID: "test-chaincode",
			Function:    "test-function",
			Args:        []string{"arg1", "arg2"},
		}
		bodyBytes, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/invoke", bytes.NewBuffer(bodyBytes))
		w := httptest.NewRecorder()

		mockInvokeService.EXPECT().ExecuteInvoke("test-channel", "test-chaincode", "test-function", []string{"arg1", "arg2"}).
			Return(&domain.TransactionProposalStatus{}, errors.New("invoke error"))

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		require.Contains(t, w.Body.String(), "Error executing invoke")
	})

	t.Run("should return 500 when ReplayEvents returns an error", func(t *testing.T) {
		body := dto.InvokeRequest{
			ChannelID:   "test-channel",
			ChaincodeID: "test-chaincode",
			Function:    "test-function",
			Args:        []string{"arg1", "arg2"},
		}
		bodyBytes, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/invoke", bytes.NewBuffer(bodyBytes))
		w := httptest.NewRecorder()

		txnStatus := &domain.TransactionProposalStatus{
			BlockNumber:   10,
			TransactionID: "test-txn-id",
		}
		mockInvokeService.EXPECT().ExecuteInvoke("test-channel", "test-chaincode", "test-function", []string{"arg1", "arg2"}).
			Return(txnStatus, nil)

		mockEventService.EXPECT().ReplayEvents(gomock.Any(), "test-channel", "test-chaincode", 10, "test-txn-id").
			Return(nil, errors.New("replay events error"))

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		require.Contains(t, w.Body.String(), "Error replaying event")
	})

	t.Run("should return 500 when HandleEvent returns an error", func(t *testing.T) {
		body := dto.InvokeRequest{
			ChannelID:   "test-channel",
			ChaincodeID: "test-chaincode",
			Function:    "test-function",
			Args:        []string{"arg1", "arg2"},
		}
		bodyBytes, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/invoke", bytes.NewBuffer(bodyBytes))
		w := httptest.NewRecorder()

		txnStatus := &domain.TransactionProposalStatus{
			BlockNumber:   10,
			TransactionID: "test-txn-id",
		}
		mockInvokeService.EXPECT().ExecuteInvoke("test-channel", "test-chaincode", "test-function", []string{"arg1", "arg2"}).
			Return(txnStatus, nil)

		mockEventService.EXPECT().ReplayEvents(gomock.Any(), "test-channel", "test-chaincode", 10, "test-txn-id").
			Return(make(<-chan *domain.Events), nil)

		mockEventService.EXPECT().HandleEvent(gomock.Any(), "test-txn-id", constants.EventLogFilename).
			Return(errors.New("handle event error"))

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		require.Contains(t, w.Body.String(), "Error handling event")
	})

	t.Run("should return 200 when request is successful", func(t *testing.T) {
		body := dto.InvokeRequest{
			ChannelID:   "test-channel",
			ChaincodeID: "test-chaincode",
			Function:    "test-function",
			Args:        []string{"arg1", "arg2"},
		}
		bodyBytes, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/invoke", bytes.NewBuffer(bodyBytes))
		w := httptest.NewRecorder()

		txnStatus := &domain.TransactionProposalStatus{
			BlockNumber:   10,
			TransactionID: "test-txn-id",
		}
		mockInvokeService.EXPECT().ExecuteInvoke("test-channel", "test-chaincode", "test-function", []string{"arg1", "arg2"}).
			Return(txnStatus, nil)

		mockEventService.EXPECT().ReplayEvents(gomock.Any(), "test-channel", "test-chaincode", 10, "test-txn-id").
			Return(make(<-chan *domain.Events), nil)

		mockEventService.EXPECT().HandleEvent(gomock.Any(), "test-txn-id", constants.EventLogFilename).
			Return(nil)

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		require.Contains(t, w.Body.String(), `"success":true`)
	})
}
