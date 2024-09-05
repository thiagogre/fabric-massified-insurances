package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// hashPassword helper function to hash a password for testing purposes
func hashPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword)
}

// TestCheckPasswordHash tests for different scenarios using shared setup
func TestCheckPasswordHash(t *testing.T) {
	// Setup common test data
	password := "valid_password"
	hashedPassword := hashPassword(password)

	t.Run("CorrectPassword", func(t *testing.T) {
		require.True(t, CheckPasswordHash(password, hashedPassword))
	})

	t.Run("WrongPassword", func(t *testing.T) {
		wrongPassword := "wrongpassword"
		require.False(t, CheckPasswordHash(wrongPassword, hashedPassword))
	})

	t.Run("EmptyPassword", func(t *testing.T) {
		require.False(t, CheckPasswordHash("", hashedPassword))
	})
}
