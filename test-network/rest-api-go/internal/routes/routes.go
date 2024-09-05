package routes

import (
	"net/http"

	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/handlers"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/db"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/org"
)

func RegisterRoutes(database db.Database, orgSetup org.OrgSetup) {
	http.Handle("/auth", handlers.InitAuthHandler(database))
	http.Handle("/events", handlers.InitEventHandler())
	http.Handle("/identity", handlers.InitIdentityHandler())
	http.Handle("/invoke", handlers.InitInvokeHandler(orgSetup))
	http.Handle("/query", handlers.InitQueryHandler(orgSetup))
}
