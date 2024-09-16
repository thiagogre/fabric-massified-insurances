package tests

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"os"

	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/constants"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/db"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/logger"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/utils"
)

func CreateTestFileHeader(filename string) *multipart.FileHeader {
	buf := new(bytes.Buffer)
	return &multipart.FileHeader{
		Filename: filename,
		Size:     int64(buf.Len()),
	}
}

func CleanupTestEventLog() {
	os.Remove(constants.EventLogFilename)
}

func SetupTestEventLog() error {
	f, err := os.Create(constants.EventLogFilename)
	if err != nil {
		return err
	}

	if err := SeedTestEventLogData(f); err != nil {
		return err
	}

	return nil
}

func SeedTestEventLogData(f *os.File) error {
	events := []domain.Event{
		{BlockNumber: 10, TransactionID: "fea0ff7b4dd6df3a90ea7ed45d5c352eb92c7041aa6d762f8563a8edc0f937f9", ChaincodeName: "basic", EventName: "DeleteAsset", Payload: "eyJDbGFpbVN0YXR1cyI6IkFjdGl2ZSIsIkNvdmVyYWdlQW1vdW50Ijo1MDAwLCJJRCI6InBvbGljeTEiLCJJbnN1cmVkSXRlbSI6IlNtYXJ0cGhvbmUgQURCIiwiT3duZXIiOiJEb25vIiwiUHJlbWl1bSI6MzAwLCJUZXJtIjoxMn0="},
		{BlockNumber: 11, TransactionID: "aced617a5f09b372aabc3e4abf57db25d7fd4fa2956d0688bdb112927693f463", ChaincodeName: "basic", EventName: "CreateAsset", Payload: "eyJDbGFpbVN0YXR1cyI6IkFjdGl2ZSIsIkNvdmVyYWdlQW1vdW50Ijo1MDAwLCJJRCI6ImIxMmZhY2M1LTY2NzQtNDBjZC04ZDNmLWM5NmU2OGZjZGI0OSIsIkluc3VyZWRJdGVtIjoiU21hcnRwaG9uZSBBQkMiLCJPd25lciI6IkRvbm8iLCJQcmVtaXVtIjozMDAsIlRlcm0iOjEyfQ=="},
		{BlockNumber: 11, TransactionID: "1964747d81c8c1346e71c9ebdef6856a39946122844ff6f7de77925c3cced5c2", ChaincodeName: "basic", EventName: "CreateAsset", Payload: "eyJDbGFpbVN0YXR1cyI6IkFjdGl2ZSIsIkNvdmVyYWdlQW1vdW50Ijo1MDAwLCJJRCI6ImRhM2FlZDU4LTg0MjgtNGQyNS05OTQ2LWMzOWZmNjhjMWFkNyIsIkluc3VyZWRJdGVtIjoiU21hcnRwaG9uZSBBQkMiLCJPd25lciI6IkRvbm8iLCJQcmVtaXVtIjozMDAsIlRlcm0iOjEyfQ=="},
		{BlockNumber: 11, TransactionID: "ea9db3a87be405f0051c1f07266b68db19750eea4f917396404cd1eafa250954", ChaincodeName: "basic", EventName: "CreateAsset", Payload: "eyJDbGFpbVN0YXR1cyI6IkFjdGl2ZSIsIkNvdmVyYWdlQW1vdW50Ijo1MDAwLCJJRCI6IjA4MmVkYjJkLWQ0Y2YtNDI5MC05ZTA1LTdkMzk2MzAwODRmYyIsIkluc3VyZWRJdGVtIjoiU21hcnRwaG9uZSBBQkMiLCJPd25lciI6IkRvbm8iLCJQcmVtaXVtIjozMDAsIlRlcm0iOjEyfQ=="},
		{BlockNumber: 11, TransactionID: "62039ffe61cedd018fa0db65d31af7f54b4db26440fb4c47290c444ddfef617f", ChaincodeName: "basic", EventName: "CreateAsset", Payload: "eyJDbGFpbVN0YXR1cyI6IkFjdGl2ZSIsIkNvdmVyYWdlQW1vdW50Ijo1MDAwLCJJRCI6IjI2YmYxY2Y2LWU2NzctNGM3ZC05MzczLTU4NTRkMzI5YzVlYSIsIkluc3VyZWRJdGVtIjoiU21hcnRwaG9uZSBBQkMiLCJPd25lciI6IkRvbm8iLCJQcmVtaXVtIjozMDAsIlRlcm0iOjEyfQ=="},
		{BlockNumber: 11, TransactionID: "1e3fbeacdc920f5ed34e5c24f3889b97bd9892dc893d33afed45187899188797", ChaincodeName: "basic", EventName: "CreateAsset", Payload: "eyJDbGFpbVN0YXR1cyI6IkFjdGl2ZSIsIkNvdmVyYWdlQW1vdW50Ijo1MDAwLCJJRCI6ImFkODQ3NWMwLTViYmQtNGU1NS1hYzkyLTNkMmJmNGIyMGJhNyIsIkluc3VyZWRJdGVtIjoiU21hcnRwaG9uZSBBQkMiLCJPd25lciI6IkRvbm8iLCJQcmVtaXVtIjozMDAsIlRlcm0iOjEyfQ=="},
		{BlockNumber: 11, TransactionID: "3a80fd54aeda5f2b81989e019cf423f4ac85827f7bb61a042abba66f6b943202", ChaincodeName: "basic", EventName: "CreateAsset", Payload: "eyJDbGFpbVN0YXR1cyI6IkFjdGl2ZSIsIkNvdmVyYWdlQW1vdW50Ijo1MDAwLCJJRCI6ImYzOGRlYmE1LTAxYjAtNGI1NS1iYjJkLTYwZGUyMDMwNTg3ZSIsIkluc3VyZWRJdGVtIjoiU21hcnRwaG9uZSBBQkMiLCJPd25lciI6IkRvbm8iLCJQcmVtaXVtIjozMDAsIlRlcm0iOjEyfQ=="},
		{BlockNumber: 12, TransactionID: "9495e73605f28b3bcb814705e57df8279096286733ae9364c3eda307804937eb", ChaincodeName: "basic", EventName: "CreateAsset", Payload: "eyJDbGFpbVN0YXR1cyI6IkFjdGl2ZSIsIkNvdmVyYWdlQW1vdW50Ijo1MDAwLCJJRCI6ImU4NjFjOTFjLTIyOGUtNDgyOS04OGY2LTdjNjk4NmJkNGU3YSIsIkluc3VyZWRJdGVtIjoiU21hcnRwaG9uZSBBQkMiLCJPd25lciI6IkRvbm8iLCJQcmVtaXVtIjozMDAsIlRlcm0iOjEyfQ=="},
		{BlockNumber: 12, TransactionID: "af9a8a6074cf8dc723cafeed060b0ab8a2d88a413a0e853adb77dcc299a356a7", ChaincodeName: "basic", EventName: "CreateAsset", Payload: "eyJDbGFpbVN0YXR1cyI6IkFjdGl2ZSIsIkNvdmVyYWdlQW1vdW50Ijo1MDAwLCJJRCI6ImU1MTQ0YWUzLTQ5M2YtNDBjYS1iZGI4LWIwNDU1NWI3YzZlNSIsIkluc3VyZWRJdGVtIjoiU21hcnRwaG9uZSBBQkMiLCJPd25lciI6IkRvbm8iLCJQcmVtaXVtIjozMDAsIlRlcm0iOjEyfQ=="},
		{BlockNumber: 12, TransactionID: "a00c06beae77291c535a60ad982a77bc7d2d02cde150e456d7f6796ab5436517", ChaincodeName: "basic", EventName: "CreateAsset", Payload: "eyJDbGFpbVN0YXR1cyI6IkFjdGl2ZSIsIkNvdmVyYWdlQW1vdW50Ijo1MDAwLCJJRCI6ImEwMTM3MzBkLTJiODEtNDNiOC04NDVkLTY3M2EzYzg1ZjkzNiIsIkluc3VyZWRJdGVtIjoiU21hcnRwaG9uZSBBQkMiLCJPd25lciI6IkRvbm8iLCJQcmVtaXVtIjozMDAsIlRlcm0iOjEyfQ=="},
	}

	for _, event := range events {
		eventData, err := json.Marshal(event)
		if err != nil {
			return err
		}

		if _, err := f.WriteString(string(eventData) + "\n"); err != nil {
			return err
		}

	}

	return nil
}

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
