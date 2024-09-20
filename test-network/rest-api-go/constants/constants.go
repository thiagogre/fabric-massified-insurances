package constants

import "github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/db"

const (
	TestNetworkPath             = "../"
	DBType                      = db.SQLite
	DBPath                      = TestNetworkPath + "organizations/fabric-ca/org1/fabric-ca-server.db"
	ServerAddr                  = ":3001"
	DefaultUploadDir            = "./uploads"
	DefaultUsernameLength       = 16
	DefaultPasswordLength       = 24
	MaxFileSize           int64 = 10 << 20 // 10MB

	RedColorOuput      = "\033[31m"
	GreenColorOuput    = "\033[32m"
	YellowColorOutput  = "\033[33m"
	BlueColorOutput    = "\033[34m"
	DefaultColorOutput = "\033[0m"

	EventLogFilename = "events.log"

	TestUsername       = "valid_username"
	TestPassword       = "TestUser"
	TestHashedPassword = "$2a$10$Zlq4NHbbVu60EvboKwNY6eUUyyFy2fNldqfQCn7Cs5bnvnN10CjK6"
)
