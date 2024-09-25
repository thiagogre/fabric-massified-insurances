package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetFullHostURL_HTTP(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)

	str := GetFullHostURL(req)
	require.Equal(t, "example.com", req.Host)
	require.Equal(t, "http://example.com", str)
}

func TestGetFullHostURL_HTTPS(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "https://example.com", nil)

	str := GetFullHostURL(req)
	require.Equal(t, "example.com", req.Host)
	require.Equal(t, "https://example.com", str)
}
