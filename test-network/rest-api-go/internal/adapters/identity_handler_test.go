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
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/utils"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/tests"
)

func TestExecute_Create_Success(t *testing.T) {
	tests.SetupLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIdentityService := mocks.NewMockIdentityInterface(ctrl)
	handler := NewIdentityHandler(mockIdentityService)

	credentials := &domain.Credentials{Username: "testuser", Password: "testpassword"}

	mockIdentityService.EXPECT().
		Create().
		Return(credentials, nil)

	req := httptest.NewRequest(http.MethodPost, "/identity", nil)
	rec := httptest.NewRecorder()

	handler.Execute(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	expected := domain.SuccessResponse[domain.IdentityResponse]{Success: true, Data: domain.IdentityResponse{Username: credentials.Username, Password: credentials.Password}}
	utils.AssertJSONResponse(t, rec.Body.String(), expected)
}

func TestExecute_Create_Fail(t *testing.T) {
	tests.SetupLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockIdentityService := mocks.NewMockIdentityInterface(ctrl)
	handler := NewIdentityHandler(mockIdentityService)

	mockIdentityService.EXPECT().
		Create().
		Return(nil, errors.New("error creating credentials"))

	req := httptest.NewRequest(http.MethodPost, "/identity", nil)
	rec := httptest.NewRecorder()

	handler.Execute(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	require.Contains(t, rec.Body.String(), "Error creating new random credentials")
}
