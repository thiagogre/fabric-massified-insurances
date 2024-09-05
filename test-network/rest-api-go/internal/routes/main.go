package routes

import (
	"net/http"

	"github.com/rs/cors"

	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/constants"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/db"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/logger"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/org"
)

const (
	dbType     = db.SQLite
	dbPath     = constants.TestNetworkPath + "organizations/fabric-ca/org1/fabric-ca-server.db"
	serverAddr = ":3001"
)

// Serve starts http web server.
func Serve(orgSetup org.OrgSetup) {
	handler := cors.Default().Handler(http.DefaultServeMux)

	database, err := db.NewDatabase(dbType, dbPath)
	if err != nil {
		logger.Error("Error opening database: " + err.Error())
		return
	}
	defer database.Close()

	RegisterRoutes(database, orgSetup)

	if err := http.ListenAndServe(serverAddr, handler); err != nil {
		logger.Error(err.Error())
	}

	logger.Info("Listening (http://localhost:3001)")
}
