package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/constants"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/dto"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/handlers"
)

func TestAuthHandler(t *testing.T) {
	testDB, err := Setup()
	require.NoError(t, err)
	defer testDB.Close()

	authHandler := handlers.InitAuthHandler(testDB)

	t.Run("ValidCredentials", func(t *testing.T) {
		requestBody := dto.AuthRequest{Username: constants.TestUsername, Password: constants.TestPassword}
		body, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		w := httptest.NewRecorder()

		authHandler.ServeHTTP(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response dto.SuccessResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
	})

	t.Run("InvalidCredentials", func(t *testing.T) {
		requestBody := dto.AuthRequest{Username: constants.TestUsername, Password: "wrongpassword"}
		body, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		w := httptest.NewRecorder()

		authHandler.ServeHTTP(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		requestBody := dto.AuthRequest{Username: "nonexistentuser", Password: "password"}
		body, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		w := httptest.NewRecorder()

		authHandler.ServeHTTP(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
