package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"

	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/constants"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/dto"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/handlers"
)

func TestIdentityHandlerSuccess(t *testing.T) {
	SetupLogger()

	mockCommandExecutor := &MockCommandExecutor{
		Output: []byte("mocked output"),
		Err:    nil,
	}

	identityHandler := handlers.InitIdentityHandler(mockCommandExecutor)

	requestBody := dto.IdentityRequest{Username: constants.TestUsername}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/identity", bytes.NewReader(body))
	w := httptest.NewRecorder()

	identityHandler.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response dto.SuccessResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
}

func TestIdentityHandlerErrorExecutingScript(t *testing.T) {
	SetupLogger()

	mockCommandExecutor := &MockCommandExecutor{
		Output: []byte("mocked output"),
		Err:    errors.New("Error executing script"),
	}

	identityHandler := handlers.InitIdentityHandler(mockCommandExecutor)

	requestBody := dto.IdentityRequest{Username: constants.TestUsername}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/identity", bytes.NewReader(body))
	w := httptest.NewRecorder()

	identityHandler.ServeHTTP(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var response dto.SuccessResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
}
