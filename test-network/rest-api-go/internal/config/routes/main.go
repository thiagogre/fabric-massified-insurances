package config

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/constants"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/adapters"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/application"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/db"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/logger"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/org"
)

func Serve(orgSetup org.OrgSetup) {
	router := mux.NewRouter()

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
	claimRepository := adapters.NewClaimRepository()

	authService := application.NewAuthService(userRepository)
	eventService := application.NewEventService(blockchainGateway, eventRepository)
	identityService := application.NewIdentityService(commandExecutor)
	invokeService := application.NewInvokeService(blockchainGateway)
	queryService := application.NewQueryService(blockchainGateway)
	claimService := application.NewClaimService(claimRepository)

	authHandler := adapters.NewAuthHandler(authService)
	eventHandler := adapters.NewEventHandler(eventService)
	identityHandler := adapters.NewIdentityHandler(identityService)
	invokeHandler := adapters.NewInvokeHandler(invokeService, eventService)
	smartContractHandler := adapters.NewSmartContractHandler()
	queryHandler := adapters.NewQueryHandler(queryService)
	claimHandler := adapters.NewClaimHandler(claimService)

	router.HandleFunc("/auth", authHandler.Execute).Methods("POST")
	router.HandleFunc("/event", eventHandler.GetAll).Methods("GET")
	router.HandleFunc("/identity", identityHandler.Execute).Methods("POST")

	smartContractRoutes := router.PathPrefix("/smartcontract").Subrouter()
	smartContractRoutes.HandleFunc("", smartContractHandler.Info).Methods("GET")
	smartContractRoutes.HandleFunc("/invoke", invokeHandler.Execute).Methods("POST")
	smartContractRoutes.HandleFunc("/query", queryHandler.Execute).Methods("GET")

	claimRoutes := router.PathPrefix("/claim").Subrouter()
	claimRoutes.HandleFunc("", claimHandler.Execute).Methods("POST")
	claimRoutes.HandleFunc("/evidence/{username}", claimHandler.GetPDFs).Methods("GET")
	claimRoutes.HandleFunc("/evidence/{username}/{filename}", claimHandler.ServePDF).Methods("GET")
	claimRoutes.HandleFunc("/evidence/validate", claimHandler.Validate).Methods("POST")
	claimRoutes.HandleFunc("/finish", claimHandler.Finish).Methods("POST")

	handler := cors.Default().Handler(router)
	if err := http.ListenAndServe(constants.ServerAddr, handler); err != nil {
		logger.Error(err.Error())
	}

	logger.Info(fmt.Sprintf("Listening: (http://localhost%s)", constants.ServerAddr))
}
