package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/dto"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/models"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/tests"
)

func TestEventHandlerSuccess(t *testing.T) {
	tests.SetupLogger()
	err := tests.SetupTestEventLog()
	require.NoError(t, err)

	eventHandler := InitEventHandler()

	t.Run("Valid request with events", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/events", nil)
		w := httptest.NewRecorder()
		eventHandler.ServeHTTP(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response dto.DocsResponse[models.Event]
		err := json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Len(t, response.Docs, 10)
		assert.Equal(t, 10, int(response.Docs[0].BlockNumber))
	})

	t.Run("Failed to open file", func(t *testing.T) {
		tests.CleanupTestEventLog()

		req := httptest.NewRequest(http.MethodGet, "/events", nil)
		w := httptest.NewRecorder()
		eventHandler.ServeHTTP(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
