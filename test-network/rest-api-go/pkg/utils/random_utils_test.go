package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateRandomBytes(t *testing.T) {
	length := 16
	bytes, err := GenerateRandomBytes(length)
	require.NoError(t, err)
	require.Len(t, bytes, length)
}

func TestGenerateRandomString(t *testing.T) {
	length := 12
	str, err := GenerateRandomString(length)
	require.NoError(t, err)
	require.Len(t, str, length)
}
