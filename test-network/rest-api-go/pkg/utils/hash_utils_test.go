package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/constants"
)

// TestCheckPasswordHash tests for different scenarios using shared setup
func TestCheckPasswordHash(t *testing.T) {
	t.Run("CorrectPassword", func(t *testing.T) {
		require.True(t, CheckPasswordHash(constants.TestPassword, constants.TestHashedPassword))
	})

	t.Run("WrongPassword", func(t *testing.T) {
		wrongPassword := "wrongpassword"
		require.False(t, CheckPasswordHash(wrongPassword, constants.TestHashedPassword))
	})

	t.Run("EmptyPassword", func(t *testing.T) {
		require.False(t, CheckPasswordHash("", constants.TestHashedPassword))
	})
}
