package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/constants"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/db"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/logger"
)

func TestSetupLogger(t *testing.T) {
	SetupLogger()

	require.NotNil(t, logger.InfoLogger)
	require.NotNil(t, logger.WarnLogger)
	require.NotNil(t, logger.SuccessLogger)
	require.NotNil(t, logger.ErrorLogger)
}

func TestCreateTestSchema(t *testing.T) {
	testDB, err := db.NewDatabase(db.SQLite, ":memory:")
	require.NoError(t, err)
	defer testDB.Close()

	err = CreateTestSchema(testDB)
	require.NoError(t, err)

	_, err = testDB.Query("SELECT * FROM users LIMIT 1;")
	require.NoError(t, err)
}

func TestSeedTestData(t *testing.T) {
	testDB, err := db.NewDatabase(db.SQLite, ":memory:")
	require.NoError(t, err)
	defer testDB.Close()

	err = CreateTestSchema(testDB)
	require.NoError(t, err)

	err = SeedTestData(testDB)
	require.NoError(t, err)

	var id string
	var token string
	row := testDB.QueryRow("SELECT id, token FROM users WHERE id = ?", constants.TestUsername)
	err = row.Scan(&id, &token)
	require.NoError(t, err)

	assert.Equal(t, constants.TestUsername, id)
	assert.NotEmpty(t, token)
}

func TestSetupTestDatabase(t *testing.T) {
	testDB, err := SetupTestDatabase()
	require.NoError(t, err)
	defer testDB.Close()
}

func TestSetup(t *testing.T) {
	testDB, err := Setup()
	require.NoError(t, err)
	defer testDB.Close()
}
