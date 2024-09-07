package routes

import (
	"net/http"

	"github.com/rs/cors"

	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/constants"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/cmd"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/db"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/logger"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/org"
)

// Serve starts http web server.
func Serve(orgSetup org.OrgSetup) {
	handler := cors.Default().Handler(http.DefaultServeMux)

	database, err := db.NewDatabase(constants.DBType, constants.DBPath)
	if err != nil {
		logger.Error("Error opening database: " + err.Error())
		return
	}
	defer database.Close()

	commandExecutor := &cmd.CommandExecutor{}

	RegisterRoutes(database, orgSetup, commandExecutor)

	if err := http.ListenAndServe(constants.ServerAddr, handler); err != nil {
		logger.Error(err.Error())
	}

	logger.Info("Listening (http://localhost:3001)")
}
