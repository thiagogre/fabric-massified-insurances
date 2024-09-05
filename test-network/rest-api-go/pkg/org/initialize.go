package org

import (
	"crypto/x509"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/logger"
)

// OrgSetup contains the organization's configuration for interacting with the network.
type OrgSetup struct {
	OrgName      string
	MSPID        string
	CryptoPath   string
	CertPath     string
	KeyPath      string
	TLSCertPath  string
	PeerEndpoint string
	GatewayPeer  string
	Gateway      client.Gateway
}

// Initialize sets up the organization configuration for interacting with the network.
func Initialize(orgSetup OrgSetup) (*OrgSetup, error) {
	logger.Info("Initializing connection for " + orgSetup.OrgName)

	clientConnection, err := orgSetup.newGrpcConnection()
	if err != nil {
		return nil, handleError(err, "Failed to create gRPC connection")
	}

	id, err := orgSetup.newIdentity()
	if err != nil {
		return nil, handleError(err, "Failed to create identity")
	}

	sign, err := orgSetup.newSign()
	if err != nil {
		return nil, handleError(err, "Failed to create sign function")
	}

	gateway, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		return nil, handleError(err, "Failed to connect to gateway")
	}

	orgSetup.Gateway = *gateway
	logger.Info("Initialization complete")

	return &orgSetup, nil
}

// newGrpcConnection creates a gRPC connection to the Gateway server.
func (orgSetup OrgSetup) newGrpcConnection() (*grpc.ClientConn, error) {
	certificate, err := loadCertificate(orgSetup.TLSCertPath)
	if err != nil {
		return nil, handleError(err, "Failed to load TLS certificate")
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, orgSetup.GatewayPeer)

	connection, err := grpc.NewClient(orgSetup.PeerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		return nil, handleError(err, "Failed to create gRPC connection")
	}

	return connection, nil
}

// newIdentity creates a client identity for this Gateway connection using an X.509 certificate.
func (orgSetup OrgSetup) newIdentity() (*identity.X509Identity, error) {
	certificate, err := loadCertificate(orgSetup.CertPath)
	if err != nil {
		return nil, handleError(err, "Failed to load certificate")
	}

	id, err := identity.NewX509Identity(orgSetup.MSPID, certificate)
	if err != nil {
		return nil, handleError(err, "Failed to create X509 identity")
	}

	return id, nil
}

// newSign creates a function that generates a digital signature from a message digest using a private key.
func (orgSetup OrgSetup) newSign() (identity.Sign, error) {
	files, err := os.ReadDir(orgSetup.KeyPath)
	if err != nil {
		return nil, handleError(err, "Failed to read private key directory")
	}

	privateKeyPEM, err := os.ReadFile(path.Join(orgSetup.KeyPath, files[0].Name()))
	if err != nil {
		return nil, handleError(err, "Failed to read private key file")
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		return nil, handleError(err, "Failed to parse private key")
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		return nil, handleError(err, "Failed to create sign function")
	}

	return sign, nil
}

// loadCertificate loads an X.509 certificate from a PEM file.
func loadCertificate(filename string) (*x509.Certificate, error) {
	certificatePEM, err := os.ReadFile(filename)
	if err != nil {
		return nil, handleError(err, "Failed to read certificate file")
	}
	return identity.CertificateFromPEM(certificatePEM)
}

// handleError logs the error and returns a formatted error with context.
func handleError(err error, context string) error {
	errMessage := context
	logger.Error(errMessage + ": " + err.Error())
	return fmt.Errorf("%s: %w", errMessage, err)
}
