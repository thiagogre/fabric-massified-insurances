package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"rest-api-go/pkg/org"
)

type QueryHandler struct {
	OrgSetup org.OrgSetup
}

func InitQueryHandler(orgSetup org.OrgSetup) *QueryHandler {
	return &QueryHandler{OrgSetup: orgSetup}
}

func (h *QueryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Query request")
	queryParams := r.URL.Query()
	chainCodeName := queryParams.Get("chaincodeid")
	channelID := queryParams.Get("channelid")
	function := queryParams.Get("function")
	args := r.URL.Query()["args"]
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)
	network := h.OrgSetup.Gateway.GetNetwork(channelID)
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
