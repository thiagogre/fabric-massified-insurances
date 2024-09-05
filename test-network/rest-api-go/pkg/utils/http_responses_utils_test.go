package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func newTestResponseRecorder() *httptest.ResponseRecorder {
	return httptest.NewRecorder()
}

func TestErrorResponse(t *testing.T) {
	recorder := newTestResponseRecorder()
	statusCode := http.StatusBadRequest
	message := "An error occurred"

	ErrorResponse(recorder, statusCode, message)

	response := recorder.Result()
	require.Equal(t, statusCode, response.StatusCode)

	expectedBody := `{"success":false,"message":"An error occurred"}`
	body := recorder.Body.String()
	require.JSONEq(t, expectedBody, body)
}

func TestSuccessResponse(t *testing.T) {
	recorder := newTestResponseRecorder()
	statusCode := http.StatusOK
	responsePayload := map[string]string{"key": "value"}

	SuccessResponse(recorder, statusCode, responsePayload)

	response := recorder.Result()
	require.Equal(t, statusCode, response.StatusCode)

	expectedBody := `{"key":"value"}`
	body := recorder.Body.String()
	require.JSONEq(t, expectedBody, body)
}

func TestResponse(t *testing.T) {
	recorder := newTestResponseRecorder()
	statusCode := http.StatusOK
	payload := map[string]string{"key": "value"}

	Response(recorder, statusCode, payload)

	response := recorder.Result()
	require.Equal(t, statusCode, response.StatusCode)

	expectedBody := `{"key":"value"}`
	body := recorder.Body.String()
	require.JSONEq(t, expectedBody, body)
}

func TestResponse_ErrorEncoding(t *testing.T) {
	recorder := newTestResponseRecorder()
	statusCode := http.StatusOK
	payload := make(chan int) // Unmarshalable payload to simulate an encoding error

	Response(recorder, statusCode, payload)

	response := recorder.Result()
	require.Equal(t, http.StatusInternalServerError, response.StatusCode)

	expectedBody := "Failed to encode response\n"
	body := recorder.Body.String()
	require.Equal(t, expectedBody, body)
}
