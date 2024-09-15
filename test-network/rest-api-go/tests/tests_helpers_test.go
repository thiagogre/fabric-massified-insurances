package tests

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/constants"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain"
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

func TestSetupTestEventLog(t *testing.T) {
	CleanupTestEventLog()
	err := SetupTestEventLog()
	require.NoError(t, err)
	defer CleanupTestEventLog()

	_, err = os.Stat(constants.EventLogFilename)
	require.False(t, os.IsNotExist(err), "The event log file should have been created")
}

func TestSeedTestEventLogData(t *testing.T) {
	f, err := os.Create(constants.EventLogFilename)
	require.NoError(t, err)
	defer CleanupTestEventLog()

	err = SeedTestEventLogData(f)
	require.NoError(t, err)

	file, err := os.Open(constants.EventLogFilename)
	require.NoError(t, err)
	defer file.Close()

	var events []domain.Event
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var event domain.Event
		err := json.Unmarshal([]byte(scanner.Text()), &event)
		require.NoError(t, err)

		events = append(events, event)
	}

	assert.Len(t, events, 10, "The event log file should contain 10 events")

}

func TestCleanupTestEventLog(t *testing.T) {
	f, err := os.Create(constants.EventLogFilename)
	require.NoError(t, err)
	f.Close()

	_, err = os.Stat(constants.EventLogFilename)
	require.False(t, os.IsNotExist(err))

	CleanupTestEventLog()

	_, err = os.Stat(constants.EventLogFilename)
	assert.True(t, os.IsNotExist(err), "The event log file should have been removed")
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
