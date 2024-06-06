package routes

import (
	"fmt"
	"net/http"

	"github.com/rs/cors"

	"rest-api-go/pkg/db"
	"rest-api-go/pkg/org"
)

const (
	dbType     = db.SQLite
	dbPath     = "../organizations/fabric-ca/org1/fabric-ca-server.db"
	serverAddr = ":3001"
)

// Serve starts http web server.
func Serve(orgSetup org.OrgSetup) {
	handler := cors.Default().Handler(http.DefaultServeMux)

	database, err := db.NewDatabase(dbType, dbPath)
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer database.Close()

	RegisterRoutes(database, orgSetup)

	fmt.Printf("Listening (http://localhost:%s/)...\n", serverAddr)
	if err := http.ListenAndServe(serverAddr, handler); err != nil {
		fmt.Println(err)
	}
}
