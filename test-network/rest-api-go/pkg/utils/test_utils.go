package utils

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func AssertJSONResponse(t *testing.T, actualBody string, expected interface{}) {
	var actual interface{}
	err := json.Unmarshal([]byte(actualBody), &actual)
	require.NoError(t, err, "Failed to unmarshal actual response body")

	expectedJSON, err := json.Marshal(expected)
	require.NoError(t, err, "Failed to marshal expected response")

	require.JSONEq(t, string(expectedJSON), actualBody, "Response body does not match expected JSON")
}
