package tests

import (
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/constants"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/db"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/logger"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/utils"
)

func SetupLogger() {
	// FIXME: If we don't initialize logger our tests don't run.
	logger.Init()
}

func SetupTestDatabase() (db.Database, error) {
	testDB, err := db.NewDatabase(db.SQLite, ":memory:")
	if err != nil {
		return nil, err
	}

	if err := CreateTestSchema(testDB); err != nil {
		return nil, err
	}

	if err := SeedTestData(testDB); err != nil {
		return nil, err
	}

	return testDB, nil
}

func CreateTestSchema(testDB db.Database) error {
	_, err := testDB.Exec(`
        CREATE TABLE users (
            id TEXT PRIMARY KEY,
            token TEXT
        );
    `)
	return err
}

func SeedTestData(testDB db.Database) error {
	hashedPassword, err := utils.GenerateHash(constants.TestPassword)
	if err != nil {
		return err
	}

	_, err = testDB.Exec(`INSERT INTO users (id, token) VALUES (?, ?)`, constants.TestUsername, string(hashedPassword))
	if err != nil {
		return err
	}

	return nil
}

func Setup() (db.Database, error) {
	SetupLogger()

	testDB, err := SetupTestDatabase()
	if err != nil {
		return nil, err
	}

	return testDB, nil
}
