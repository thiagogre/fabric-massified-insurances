package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	PASSWORD        = "TestUser"
	HASHED_PASSWORD = "$2a$10$Zlq4NHbbVu60EvboKwNY6eUUyyFy2fNldqfQCn7Cs5bnvnN10CjK6"
)

// TestCheckPasswordHash tests for different scenarios using shared setup
func TestCheckPasswordHash(t *testing.T) {
	t.Run("CorrectPassword", func(t *testing.T) {
		require.True(t, CheckPasswordHash(PASSWORD, HASHED_PASSWORD))
	})

	t.Run("WrongPassword", func(t *testing.T) {
		wrongPassword := "wrongpassword"
		require.False(t, CheckPasswordHash(wrongPassword, HASHED_PASSWORD))
	})

	t.Run("EmptyPassword", func(t *testing.T) {
		require.False(t, CheckPasswordHash("", HASHED_PASSWORD))
	})
}
