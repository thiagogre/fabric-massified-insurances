package main

import (
	"rest-api-go/constants"
	"rest-api-go/internal/routes"
	"rest-api-go/pkg/logger"
	"rest-api-go/pkg/org"
)

func main() {
	logger.Init()

	//Initialize setup for Org1
	cryptoPath := constants.TestNetworkPath + "organizations/peerOrganizations/org1.example.com"
	orgConfig := org.OrgSetup{
		OrgName:      "Org1",
		MSPID:        "Org1MSP",
		CertPath:     cryptoPath + "/users/BackendClient@org1.example.com/msp/signcerts/cert.pem",
		KeyPath:      cryptoPath + "/users/BackendClient@org1.example.com/msp/keystore/",
		TLSCertPath:  cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt",
		PeerEndpoint: "dns:///localhost:7051",
		GatewayPeer:  "peer0.org1.example.com",
	}

	orgSetup, err := org.Initialize(orgConfig)
	if err != nil {
		logger.Error("Error initializing setup for Org1: " + err.Error())
	}

	routes.Serve(org.OrgSetup(*orgSetup))
}
