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

func TestGenerateRandomCredentials(t *testing.T) {
	usernameLength := 20
	passwordLength := 24
	username, password, err := GenerateRandomCredentials(usernameLength, passwordLength)
	require.NoError(t, err)
	require.Len(t, username, usernameLength)
	require.Len(t, password, passwordLength)
}
