package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"

	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/constants"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/dto"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/tests"
)

func TestAuthHandler(t *testing.T) {
	testDB, err := tests.Setup()
	require.NoError(t, err)
	defer testDB.Close()

	authHandler := InitAuthHandler(testDB)

	t.Run("ValidCredentials", func(t *testing.T) {
		requestBody := dto.AuthRequest{Username: constants.TestUsername, Password: constants.TestPassword}
		body, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		w := httptest.NewRecorder()

		authHandler.ServeHTTP(w, req)

		resp := w.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var response dto.SuccessResponse[dto.AuthRequest]
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		require.True(t, response.Success)
		dataJsonAsByte, err := json.Marshal(response.Data)
		require.NoError(t, err)
		require.JSONEq(t, string(body), string(dataJsonAsByte))
	})

	t.Run("InvalidCredentials", func(t *testing.T) {
		requestBody := dto.AuthRequest{Username: constants.TestUsername, Password: "wrongpassword"}
		body, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		w := httptest.NewRecorder()

		authHandler.ServeHTTP(w, req)

		resp := w.Result()
		require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		requestBody := dto.AuthRequest{Username: "nonexistentuser", Password: "password"}
		body, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
		w := httptest.NewRecorder()

		authHandler.ServeHTTP(w, req)

		resp := w.Result()
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
