package utils

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetQueryString(t *testing.T) {
	t.Run("Given valid query parameters with strings", func(t *testing.T) {
		queryParams := map[string]interface{}{
			"name": "John",
			"city": "New York",
		}

		queryString, err := GetQueryString(queryParams)

		require.NoError(t, err, "should not return an error")
		expected := url.Values{
			"name": []string{"John"},
			"city": []string{"New York"},
		}.Encode()
		require.Equal(t, expected, queryString, "should generate a correct query string")
	})

	t.Run("Given valid query parameters with string slices", func(t *testing.T) {
		queryParams := map[string]interface{}{
			"name": []string{"John", "Doe"},
			"city": "New York",
		}

		queryString, err := GetQueryString(queryParams)

		require.NoError(t, err, "should not return an error")
		expected := url.Values{
			"name": []string{"John", "Doe"},
			"city": []string{"New York"},
		}.Encode()
		require.Equal(t, expected, queryString, "should generate a correct query string with slices")
	})

	t.Run("Given invalid query parameter type", func(t *testing.T) {
		queryParams := map[string]interface{}{
			"name": 123, // Invalid type (int)
		}

		queryString, err := GetQueryString(queryParams)

		require.Error(t, err, "should return an error for invalid type")
		require.Equal(t, err.Error(), "invalid query parameter type for key name: int", "error message should indicate the invalid type")
		require.Equal(t, "", queryString, "query string should be empty when there is an error")
	})

	t.Run("Given empty query parameters", func(t *testing.T) {
		queryParams := map[string]interface{}{}

		queryString, err := GetQueryString(queryParams)

		require.NoError(t, err, "should not return an error for empty parameters")
		require.Equal(t, "", queryString, "should return an empty query string")
	})

	t.Run("Given mixed query parameters with strings and slices", func(t *testing.T) {
		queryParams := map[string]interface{}{
			"first_name": []string{"Alice", "Bob"},
			"last_name":  "Smith",
		}

		queryString, err := GetQueryString(queryParams)

		require.NoError(t, err, "should not return an error")
		expected := url.Values{
			"first_name": []string{"Alice", "Bob"},
			"last_name":  []string{"Smith"},
		}.Encode()
		require.Equal(t, expected, queryString, "should generate a correct query string with mixed types")
	})
}
