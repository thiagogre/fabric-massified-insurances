package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/dto"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/tests"
)

func TestIdentityHandler_Successful(t *testing.T) {
	tests.SetupLogger()

	mockCommandExecutor := &tests.MockCommandExecutor{
		Output: []byte("mocked output"),
		Err:    nil,
	}

	identityHandler := InitIdentityHandler(mockCommandExecutor)

	body, _ := json.Marshal("{}")
	req := httptest.NewRequest(http.MethodPost, "/identity", bytes.NewReader(body))
	w := httptest.NewRecorder()

	identityHandler.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var response dto.SuccessResponse[dto.IdentityResponse]
	err := json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	require.True(t, response.Success)
	require.NotEmpty(t, response.Data, "Username")
	require.NotEmpty(t, response.Data, "Password")
}

func TestIdentityHandler_Error_Executing_Script(t *testing.T) {
	tests.SetupLogger()

	mockCommandExecutor := &tests.MockCommandExecutor{
		Output: []byte("mocked output"),
		Err:    errors.New("Error executing script"),
	}

	identityHandler := InitIdentityHandler(mockCommandExecutor)
	body, _ := json.Marshal("{}")
	req := httptest.NewRequest(http.MethodPost, "/identity", bytes.NewReader(body))
	w := httptest.NewRecorder()

	identityHandler.ServeHTTP(w, req)

	resp := w.Result()
	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var response dto.ErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)
	require.False(t, response.Success)
	require.NotEmpty(t, response, "Message")
	require.Equal(t, "Error executing script", response.Message)
}
