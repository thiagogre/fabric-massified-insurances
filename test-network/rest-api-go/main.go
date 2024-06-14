package main

import (
	"bytes"
	"context"
	"encoding/json"
	"rest-api-go/constants"
	"rest-api-go/internal/routes"
	"rest-api-go/pkg/logger"
	"rest-api-go/pkg/org"

	"github.com/hyperledger/fabric-gateway/pkg/client"
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
	defer orgSetup.CancelContext()

	network := orgSetup.Gateway.GetNetwork("mychannel")

	startChaincodeEventListening(orgSetup.Context, network, "mychannel")

	routes.Serve(org.OrgSetup(*orgSetup))
}

func startChaincodeEventListening(ctx context.Context, network *client.Network, chaincodeID string) {
	logger.Info("Start chaincode event listening\n")

	events, err := network.ChaincodeEvents(ctx, chaincodeID)
	if err != nil {
		logger.Error("Failed to start chaincode event listening: " + err.Error())
		panic(err)
	}

	go func() {
		for event := range events {
			asset := formatJSON(event.Payload)
			logger.Info("Chaincode event received: " + event.EventName + "\n" + asset)
		}
	}()
}

func formatJSON(data []byte) string {
	var result bytes.Buffer
	if err := json.Indent(&result, data, "", "  "); err != nil {
		logger.Error("failed to parse JSON: " + err.Error())
		panic(err)
	}

	return result.String()
}
