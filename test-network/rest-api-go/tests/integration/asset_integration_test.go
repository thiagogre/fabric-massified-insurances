package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/constants"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/dto"
)

func TestAssetLifecycle(t *testing.T) {
	baseURL := fmt.Sprintf("http://localhost%s/smartcontract", constants.ServerAddr)

	// Step 1: Query the Ledger (Initial State)
	t.Run("Initial Query", func(t *testing.T) {
		url := baseURL + "/query?channelid=mychannel&chaincodeid=basic&function=GetAllAssets"
		httpClient := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		require.NoError(t, err)

		resp, err := httpClient.Do(req)
		require.NoError(t, err)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		defer resp.Body.Close()

		expectedResponse := dto.QuerySuccessResponse{Success: true, Data: map[string]interface{}{"docs": []interface{}{}}}
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
			"args":        []string{"policy1", "Dono", "12", "5000", "0", "Varejista", "300"},
		}

		jsonData, err := json.Marshal(assetData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		httpClient := &http.Client{}
		resp, err := httpClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var response dto.InvokeSuccessResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		require.True(t, response.Success)
		require.NotEmpty(t, response.Data, "TransactionID")
		require.NotEmpty(t, response.Data, "Code")
		require.NotEmpty(t, response.Data, "BlockNumber")
	})

	// Step 3: Query All Assets
	t.Run("Query After Create", func(t *testing.T) {
		url := baseURL + "/query?channelid=mychannel&chaincodeid=basic&args=policy1&function=ReadAsset"
		httpClient := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		require.NoError(t, err)

		resp, err := httpClient.Do(req)
		require.NoError(t, err)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		defer resp.Body.Close()

		expectedResponse := dto.QuerySuccessResponse{Success: true, Data: map[string]interface{}{
			"ClaimStatus":      "Active",
			"CoverageDuration": 12,
			"CoverageType":     0,
			"CoverageValue":    5000,
			"ID":               "policy1",
			"Insured":          "Dono",
			"Partner":          "Varejista",
			"Premium":          300,
		}}
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
			"args":        []string{"policy1", "Dono", "12", "5000", "0", "Varejista", "300", "Pending"},
		}

		jsonData, err := json.Marshal(assetData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		httpClient := &http.Client{}
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var response dto.InvokeSuccessResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		require.True(t, response.Success)
		require.NotEmpty(t, response.Data, "TransactionID")
		require.NotEmpty(t, response.Data, "Code")
		require.NotEmpty(t, response.Data, "BlockNumber")
	})

	// Step 5: Query After Update
	t.Run("Query After Update", func(t *testing.T) {
		url := baseURL + "/query?channelid=mychannel&chaincodeid=basic&args=policy1&function=ReadAsset"
		httpClient := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		require.NoError(t, err)

		resp, err := httpClient.Do(req)
		require.NoError(t, err)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		defer resp.Body.Close()

		expectedResponse := dto.QuerySuccessResponse{Success: true, Data: map[string]interface{}{
			"ClaimStatus":      "Pending",
			"CoverageDuration": 12,
			"CoverageType":     0,
			"CoverageValue":    5000,
			"ID":               "policy1",
			"Insured":          "Dono",
			"Partner":          "Varejista",
			"Premium":          300,
		}}
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

		httpClient := &http.Client{}
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var response dto.InvokeSuccessResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		require.True(t, response.Success)
		require.NotEmpty(t, response.Data, "TransactionID")
		require.NotEmpty(t, response.Data, "Code")
		require.NotEmpty(t, response.Data, "BlockNumber")
	})

	// Step 7: Query After Delete
	t.Run("Query After Delete", func(t *testing.T) {
		url := baseURL + "/query?channelid=mychannel&chaincodeid=basic&args=policy1&function=ReadAsset"
		httpClient := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		require.NoError(t, err)

		resp, err := httpClient.Do(req)
		require.NoError(t, err)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		require.Contains(t, string(body), "message")
	})
}

func TestAssetRichQuery(t *testing.T) {
	baseURL := fmt.Sprintf("http://localhost%s/smartcontract", constants.ServerAddr)
	httpClient := &http.Client{} // Reuse the HTTP client

	var assets []map[string]interface{}

	for i := 1; i <= 5; i++ {
		assetID := "policy_" + strconv.Itoa(i)
		asset := map[string]interface{}{
			"channelid":   "mychannel",
			"chaincodeid": "basic",
			"function":    "createAsset",
			"args":        []string{assetID, "Dono", "12", "5000", "0", "Varejista", "300"},
		}
		assets = append(assets, asset)
	}

	for i, assetData := range assets {
		t.Run(fmt.Sprintf("Create Asset [%d]", i), func(t *testing.T) {
			url := baseURL + "/invoke"
			jsonData, err := json.Marshal(assetData)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp, err := httpClient.Do(req)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)

			var response dto.InvokeSuccessResponse
			err = json.NewDecoder(resp.Body).Decode(&response)
			require.NoError(t, err)
			require.True(t, response.Success)
			require.NotEmpty(t, response.Data, "TransactionID")
			require.NotEmpty(t, response.Data, "Code")
			require.NotEmpty(t, response.Data, "BlockNumber")
		})
	}

	t.Run("Update Asset", func(t *testing.T) {
		url := baseURL + "/invoke"
		assetData := map[string]interface{}{
			"channelid":   "mychannel",
			"chaincodeid": "basic",
			"function":    "updateAsset",
			"args":        []string{"policy_1", "Dono", "12", "5000", "0", "Varejista", "300", "Pending"},
		}

		jsonData, err := json.Marshal(assetData)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		httpClient := &http.Client{}
		resp, err := httpClient.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		var response dto.InvokeSuccessResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		require.True(t, response.Success)
		require.NotEmpty(t, response.Data, "TransactionID")
		require.NotEmpty(t, response.Data, "Code")
		require.NotEmpty(t, response.Data, "BlockNumber")
	})

	t.Run("Get assets by claim status equal to Pending", func(t *testing.T) {
		richQuery := `{"selector":{"ClaimStatus":"Pending"}}`
		url := baseURL + "/query?channelid=mychannel&chaincodeid=basic&function=GetAssetsByRichQuery&args=" + richQuery
		httpClient := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		require.NoError(t, err)

		resp, err := httpClient.Do(req)
		require.NoError(t, err)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		defer resp.Body.Close()

		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Contains(t, string(body), `"ClaimStatus":"Pending"`)
		require.NotContains(t, string(body), `"ClaimStatus":"Active"`)
	})
}
