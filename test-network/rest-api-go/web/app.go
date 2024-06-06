package web

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rs/cors"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// OrgSetup contains organization's config to interact with the network.
type OrgSetup struct {
	OrgName       string
	MSPID         string
	CryptoPath    string
	CertPath      string
	KeyPath       string
	TLSCertPath   string
	PeerEndpoint  string
	GatewayPeer   string
	Gateway       client.Gateway
	Context       context.Context
	CancelContext context.CancelFunc
}

// Serve starts http web server.
func Serve(setup OrgSetup) {
	// Enable CORS middleware
	handler := cors.Default().Handler(http.DefaultServeMux)

	http.HandleFunc("/query", setup.Query)
	http.HandleFunc("/invoke", setup.Invoke)
	http.HandleFunc("/events", Events)
	http.HandleFunc("/identity", Identity)
	http.HandleFunc("/auth", Auth)

	fmt.Println("Listening (http://localhost:3001/)...")
	if err := http.ListenAndServe(":3001", handler); err != nil {
		fmt.Println(err)
	}
}
