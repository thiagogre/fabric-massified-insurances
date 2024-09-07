package services

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/constants"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/repositories"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/tests"
)

func TestAuthenticateUser(t *testing.T) {
	testDB, err := tests.Setup()
	require.NoError(t, err)
	defer testDB.Close()

	userRepository := &repositories.SQLUserRepository{DB: testDB}
	authService := &AuthService{UserRepository: userRepository}

	t.Run("ValidUser", func(t *testing.T) {
		user, err := authService.AuthenticateUser(constants.TestUsername, constants.TestPassword)
		require.NoError(t, err)
		require.NotNil(t, user)
		require.Equal(t, constants.TestUsername, user.Id)
	})

	t.Run("InvalidUser", func(t *testing.T) {
		user, err := authService.AuthenticateUser("nonexistent_username", constants.TestPassword)
		require.Error(t, err)
		require.Nil(t, user)
	})

	t.Run("InvalidPassword", func(t *testing.T) {
		user, err := authService.AuthenticateUser(constants.TestUsername, "invalid_pass")
		require.NoError(t, err)
		require.Equal(t, nil, err)
		require.Nil(t, user)
	})
}
