package main

import (
	"flag"

	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/constants"
	config "github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/config/routes"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/logger"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/org"
)

const cryptoPath = constants.TestNetworkPath + "organizations/peerOrganizations/org1.example.com"

func main() {
	logger.Init()

	orgName := flag.String("orgName", "Org1", "Name of the organization")
	mspID := flag.String("mspID", "Org1MSP", "MSP ID for the organization")
	certPath := flag.String("certPath", cryptoPath+"/users/BackendClient@org1.example.com/msp/signcerts/cert.pem", "Path to the certificate")
	keyPath := flag.String("keyPath", cryptoPath+"/users/BackendClient@org1.example.com/msp/keystore/", "Path to the private key")
	tlsCertPath := flag.String("tlsCertPath", cryptoPath+"/peers/peer0.org1.example.com/tls/ca.crt", "Path to the TLS certificate")
	peerEndpoint := flag.String("peerEndpoint", "dns:///localhost:7051", "Peer endpoint")
	gatewayPeer := flag.String("gatewayPeer", "peer0.org1.example.com", "Gateway peer")
	port := flag.Int("port", 3001, "Server port number")

	flag.Parse()

	//Initialize setup for Org1
	orgConfig := org.OrgSetup{
		OrgName:      *orgName,
		MSPID:        *mspID,
		CertPath:     *certPath,
		KeyPath:      *keyPath,
		TLSCertPath:  *tlsCertPath,
		PeerEndpoint: *peerEndpoint,
		GatewayPeer:  *gatewayPeer,
	}

	logger.Info(orgConfig)

	orgSetup, err := org.Initialize(orgConfig)
	if err != nil {
		logger.Error("Error initializing setup for Org1: " + err.Error())
	}

	config.Serve(org.OrgSetup(*orgSetup), *port)
}
