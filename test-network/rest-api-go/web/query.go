package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Query handles chaincode query requests.
func (setup OrgSetup) Query(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Query request")
	queryParams := r.URL.Query()
	chainCodeName := queryParams.Get("chaincodeid")
	channelID := queryParams.Get("channelid")
	function := queryParams.Get("function")
	args := r.URL.Query()["args"]
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	evaluateResponse, err := contract.EvaluateTransaction(function, args...)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}

	// Parse the JSON response into a GenericResponse interface
	var genericResponse interface{}
	err = json.Unmarshal(evaluateResponse, &genericResponse)
	if err != nil {
		http.Error(w, "Failed to parse response JSON", http.StatusInternalServerError)
		return
	}

	// Marshal the GenericResponse struct back into a JSON string
	jsonEncodedResponse, err := json.Marshal(genericResponse)
	if err != nil {
		http.Error(w, "Failed to encode response JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonEncodedResponse)
}
