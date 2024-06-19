package utils

import (
	"encoding/json"
	"net/http"

	"rest-api-go/internal/dto"
)

// ErrorResponse sends an error response with a given status code and message.
func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	Response(w, statusCode, dto.ErrorResponse{Success: false, Message: message})
}

// SuccessResponse sends a success response with a given status code and response payload.
func SuccessResponse(w http.ResponseWriter, statusCode int, response interface{}) {
	Response(w, statusCode, response)
}

// Response sends a JSON response with a given status code and payload.
func Response(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}
