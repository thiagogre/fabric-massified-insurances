package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"rest-api-go/constants"
	"rest-api-go/internal/routes"
	"rest-api-go/pkg/org"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

func main() {
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
		fmt.Println("Error initializing setup for Org1: ", err)
	}
	defer orgSetup.CancelContext()

	network := orgSetup.Gateway.GetNetwork("mychannel")

	startChaincodeEventListening(orgSetup.Context, network, "mychannel")

	routes.Serve(org.OrgSetup(*orgSetup))
}

func startChaincodeEventListening(ctx context.Context, network *client.Network, chaincodeID string) {
	fmt.Println("\n*** Start chaincode event listening")

	events, err := network.ChaincodeEvents(ctx, chaincodeID)
	if err != nil {
		panic(fmt.Errorf("failed to start chaincode event listening: %w", err))
	}

	go func() {
		for event := range events {
			asset := formatJSON(event.Payload)
			fmt.Printf("\n<-- Chaincode event received: %s - %s\n", event.EventName, asset)
		}
	}()
}

func formatJSON(data []byte) string {
	var result bytes.Buffer
	if err := json.Indent(&result, data, "", "  "); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return result.String()
}
