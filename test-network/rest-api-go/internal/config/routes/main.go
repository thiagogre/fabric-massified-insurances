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

	authHandler := adapters.NewAuthHandler(application.NewAuthService(userRepository))
	eventHandler := adapters.NewEventHandler(application.NewEventService(blockchainGateway, eventRepository))
	identityHandler := adapters.NewIdentityHandler(application.NewIdentityService(commandExecutor))
	invokeHandler := adapters.NewInvokeHandler(application.NewInvokeService(blockchainGateway), application.NewEventService(blockchainGateway, eventRepository))
	queryHandler := adapters.NewQueryHandler(application.NewQueryService(blockchainGateway))
	claimHandler := adapters.NewClaimHandler(application.NewClaimService(claimRepository))

	router.HandleFunc("/auth", authHandler.Execute).Methods("POST")
	router.HandleFunc("/event", eventHandler.GetAll).Methods("GET")
	router.HandleFunc("/identity", identityHandler.Execute).Methods("POST")
	router.HandleFunc("/invoke", invokeHandler.Execute).Methods("POST")
	router.HandleFunc("/query", queryHandler.Execute).Methods("GET")
	claimRoutes := router.PathPrefix("/claim").Subrouter()
	claimRoutes.HandleFunc("/evidence/upload", claimHandler.UploadEvidences).Methods("POST")
	// claimRoutes.HandleFunc("/evidence/{id}", claimHandler.GetEvidenceByID).Methods("GET")

	handler := cors.Default().Handler(router)
	if err := http.ListenAndServe(constants.ServerAddr, handler); err != nil {
		logger.Error(err.Error())
	}

	logger.Info(fmt.Sprintf("Listening: (http://localhost%s)", constants.ServerAddr))
}
