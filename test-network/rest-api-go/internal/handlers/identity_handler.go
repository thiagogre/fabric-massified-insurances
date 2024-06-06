package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"rest-api-go/constants"
	"rest-api-go/internal/dto"
)

type IdentityHandler struct {
}

func InitIdentityHandler() *IdentityHandler {
	return &IdentityHandler{}
}

// Register and enroll an identity.
func (h *IdentityHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Identity request")

	// Decode the JSON object
	var identityRequest dto.IdentityRequest
	if err := json.NewDecoder(r.Body).Decode(&identityRequest); err != nil {
		http.Error(w, "Failed to parse request body", http.StatusBadRequest)
		return
	}

	cmd := exec.Command("/bin/bash", "./registerEnrollIdentity.sh", identityRequest.Username)
	cmd.Dir = constants.TestNetworkPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing script:", err)
		fmt.Println("Script output:", string(output))
		return
	}
	fmt.Println("Script output:", string(output))

	// Response struct
	response := dto.IdentityRequest{
		Username: identityRequest.Username,
	}

	// Convert response to JSON
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
