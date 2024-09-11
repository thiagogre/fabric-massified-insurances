package services

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCredentials(t *testing.T) {
	credentials := NewCredentials()

	t.Run("GenerateRandomBytes", func(t *testing.T) {
		length := 16
		bytes, err := credentials.generateRandomBytes(length)
		require.NoError(t, err)
		require.Len(t, bytes, length)
	})

	t.Run("GenerateRandomString", func(t *testing.T) {
		length := 12
		str, err := credentials.generateRandomString(length)
		require.NoError(t, err)
		require.Len(t, str, length)
	})

	t.Run("Generate", func(t *testing.T) {
		usernameLength := 20
		passwordLength := 24
		err := credentials.Generate(usernameLength, passwordLength)
		require.NoError(t, err)
		require.Len(t, credentials.Username, usernameLength)
		require.Len(t, credentials.Password, passwordLength)
	})
}
