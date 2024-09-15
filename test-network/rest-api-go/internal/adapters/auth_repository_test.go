package adapters

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/constants"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/tests"
)

func TestAuthenticateUser(t *testing.T) {
	testDB, err := tests.Setup()
	require.NoError(t, err)
	defer testDB.Close()

	repo := AuthRepository{DB: testDB}

	t.Run("ValidUser", func(t *testing.T) {
		user, err := repo.GetUserById(constants.TestUsername)
		require.NoError(t, err)
		require.Equal(t, constants.TestUsername, user.Id)
	})

	t.Run("NonexistentUser", func(t *testing.T) {
		user, err := repo.GetUserById("nonexistent")
		require.Error(t, err)
		require.Nil(t, user)
	})
}
