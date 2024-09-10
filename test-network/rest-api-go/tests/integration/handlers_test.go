package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/constants"
)

func TestAssetLifecycle(t *testing.T) {
	baseURL := fmt.Sprintf("http://localhost%s", constants.ServerAddr)

	// Step 1: Query the Ledger (Initial State)
	t.Run("Initial Query", func(t *testing.T) {
		url := baseURL + "/query?channelid=mychannel&chaincodeid=basic&function=GetAllAssets"
		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		defer resp.Body.Close()

		expectedResponse := map[string]interface{}{"docs": []interface{}{}}
		response, err := json.Marshal(expectedResponse)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.JSONEq(t, string(response), string(body))
	})

	// Step 2: Create an Asset
	t.Run("Create Asset", func(t *testing.T) {
		url := baseURL + "/invoke"
		assetData := map[string]interface{}{
			"channelid":   "mychannel",
			"chaincodeid": "basic",
			"function":    "createAsset",
			"args":        []string{"policy1", "Dono", "Smartphone ABC", "5000", "300", "12"},
		}

		jsonData, err := json.Marshal(assetData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		require.Contains(t, response, "TransactionID")
		require.Equal(t, true, response["Successful"])
	})

	// Step 3: Query All Assets
	t.Run("Query After Create", func(t *testing.T) {
		url := baseURL + "/query?channelid=mychannel&chaincodeid=basic&args=policy1&function=ReadAsset"
		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		defer resp.Body.Close()

		expectedResponse := map[string]interface{}{
			"ClaimStatus":    "Active",
			"CoverageAmount": 5000,
			"ID":             "policy1",
			"InsuredItem":    "Smartphone ABC",
			"Owner":          "Dono",
			"Premium":        300,
			"Term":           12,
		}
		response, err := json.Marshal(expectedResponse)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.JSONEq(t, string(response), string(body))
	})

	// Step 4: Update the Asset
	t.Run("Update Asset", func(t *testing.T) {
		url := baseURL + "/invoke"
		assetData := map[string]interface{}{
			"channelid":   "mychannel",
			"chaincodeid": "basic",
			"function":    "updateAsset",
			"args":        []string{"policy1", "Dono", "Smartphone XYZ", "5000", "300", "12", "Active"},
		}

		jsonData, err := json.Marshal(assetData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		require.Contains(t, response, "TransactionID")
		require.Equal(t, true, response["Successful"])
	})

	// Step 5: Query After Update
	t.Run("Query After Update", func(t *testing.T) {
		url := baseURL + "/query?channelid=mychannel&chaincodeid=basic&args=policy1&function=ReadAsset"
		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		defer resp.Body.Close()

		expectedResponse := map[string]interface{}{
			"ClaimStatus":    "Active",
			"CoverageAmount": 5000,
			"ID":             "policy1",
			"InsuredItem":    "Smartphone XYZ",
			"Owner":          "Dono",
			"Premium":        300,
			"Term":           12,
		}
		response, err := json.Marshal(expectedResponse)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.JSONEq(t, string(response), string(body))
	})

	// Step 6: Delete the Asset
	t.Run("Delete Asset", func(t *testing.T) {
		url := baseURL + "/invoke"
		assetData := map[string]interface{}{
			"channelid":   "mychannel",
			"chaincodeid": "basic",
			"function":    "DeleteAsset",
			"args":        []string{"policy1"},
		}

		jsonData, err := json.Marshal(assetData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		require.Contains(t, response, "TransactionID")
		require.Equal(t, true, response["Successful"])
	})

	// Step 7: Query After Delete
	t.Run("Query After Delete", func(t *testing.T) {
		url := baseURL + "/query?channelid=mychannel&chaincodeid=basic&args=policy1&function=ReadAsset"
		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		require.Contains(t, string(body), "message")
	})
}
