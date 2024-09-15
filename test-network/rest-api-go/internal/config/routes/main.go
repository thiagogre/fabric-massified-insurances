package config

import (
	"fmt"
	"net/http"

	"github.com/rs/cors"

	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/constants"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/adapters"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/application"
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

	blockchainGateway := adapters.NewBlockchainGateway(&orgSetup.Gateway)
	eventRepository := adapters.NewEventRepository()
	userRepository := adapters.NewAuthRepository(database)
	commandExecutor := adapters.NewCommandExecutor()

	authHandler := adapters.NewAuthHandler(application.NewAuthService(userRepository))
	eventHandler := adapters.NewEventHandler(application.NewEventService(blockchainGateway, eventRepository))
	identityHandler := adapters.NewIdentityHandler(application.NewIdentityService(commandExecutor))
	invokeHandler := adapters.NewInvokeHandler(application.NewInvokeService(blockchainGateway), application.NewEventService(blockchainGateway, eventRepository))
	queryHandler := adapters.NewQueryHandler(application.NewQueryService(blockchainGateway))

	routes := map[string]http.Handler{
		"/auth":     authHandler,
		"/events":   eventHandler,
		"/identity": identityHandler,
		"/invoke":   invokeHandler,
		"/query":    queryHandler,
	}
	for path, handler := range routes {
		http.Handle(path, handler)
	}

	if err := http.ListenAndServe(constants.ServerAddr, handler); err != nil {
		logger.Error(err.Error())
	}

	logger.Info(fmt.Sprintf("Listening: (http://localhost%s)", constants.ServerAddr))
}
