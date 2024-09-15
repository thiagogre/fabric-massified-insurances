package adapters

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain/mocks"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/dto"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/utils"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/tests"
)

func TestServeHTTP_ExecuteQuery_Success(t *testing.T) {
	tests.SetupLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueryService := mocks.NewMockQueryInterface(ctrl)
	handler := NewQueryHandler(mockQueryService)

	queryParams := map[string]string{
		"chaincodeid": "testChainCode",
		"channelid":   "testChannel",
		"function":    "testFunction",
	}
	queryParamsArgs := []string{"arg1", "arg2"}

	mockQueryService.EXPECT().
		ExecuteQuery(queryParams["channelid"], queryParams["chaincodeid"], queryParams["function"], queryParamsArgs).
		Return([]byte(`{"result":"success"}`), nil)

	req := httptest.NewRequest(http.MethodGet, "/query?chaincodeid=testChainCode&channelid=testChannel&function=testFunction&args=arg1&args=arg2", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	expected := dto.QuerySuccessResponse{Success: true, Data: []byte(`{"result":"success"}`)}
	utils.AssertJSONResponse(t, rec.Body.String(), expected)
}

func TestServeHTTP_ExecuteQuery_Fail(t *testing.T) {
	tests.SetupLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQueryService := mocks.NewMockQueryInterface(ctrl)
	handler := NewQueryHandler(mockQueryService)

	queryParams := map[string]string{
		"chaincodeid": "testChainCode",
		"channelid":   "testChannel",
		"function":    "testFunction",
	}
	queryParamsArgs := []string{"arg1", "arg2"}

	mockQueryService.EXPECT().
		ExecuteQuery(queryParams["channelid"], queryParams["chaincodeid"], queryParams["function"], queryParamsArgs).
		Return(nil, errors.New("error executing query"))

	req := httptest.NewRequest(http.MethodGet, "/query?chaincodeid=testChainCode&channelid=testChannel&function=testFunction&args=arg1&args=arg2", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusInternalServerError, rec.Code)
	require.Contains(t, rec.Body.String(), "Error executing query")
}
