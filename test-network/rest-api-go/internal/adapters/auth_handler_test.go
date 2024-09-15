package adapters

import (
	"bytes"
	"encoding/json"
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

func TestServeHTTP_AuthenticateUser_Success(t *testing.T) {
	tests.SetupLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthServiceInterface(ctrl)
	handler := NewAuthHandler(mockAuthService)

	requestBody := dto.AuthRequest{
		Username: "testuser",
		Password: "testpassword",
	}
	body, _ := json.Marshal(requestBody)

	mockAuthService.EXPECT().
		AuthenticateUser(requestBody.Username, requestBody.Password).
		Return(&domain.User{Id: requestBody.Username}, nil)

	req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	expected := dto.SuccessResponse[dto.AuthRequest]{Success: true, Data: requestBody}
	utils.AssertJSONResponse(t, rec.Body.String(), expected)
}

func TestServeHTTP_AuthenticateUser_Fail_ParseBody(t *testing.T) {
	tests.SetupLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthServiceInterface(ctrl)
	handler := NewAuthHandler(mockAuthService)

	req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer([]byte("invalid json")))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "Failed to parse request body")
}

func TestServeHTTP_AuthenticateUser_Fail_UserNotFound(t *testing.T) {
	tests.SetupLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthServiceInterface(ctrl)
	handler := NewAuthHandler(mockAuthService)

	requestBody := dto.AuthRequest{
		Username: "testuser",
		Password: "testpassword",
	}
	body, _ := json.Marshal(requestBody)

	mockAuthService.EXPECT().
		AuthenticateUser(requestBody.Username, requestBody.Password).
		Return(nil, errors.New("err"))

	req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
	require.Contains(t, rec.Body.String(), "User not found")
}

func TestServeHTTP_AuthenticateUser_Fail_IncorrectPassword(t *testing.T) {
	tests.SetupLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockAuthServiceInterface(ctrl)
	handler := NewAuthHandler(mockAuthService)

	requestBody := dto.AuthRequest{
		Username: "testuser",
		Password: "testpassword",
	}
	body, _ := json.Marshal(requestBody)

	mockAuthService.EXPECT().
		AuthenticateUser(requestBody.Username, requestBody.Password).
		Return(nil, nil)

	req := httptest.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
	require.Contains(t, rec.Body.String(), "Incorrect password")
}
