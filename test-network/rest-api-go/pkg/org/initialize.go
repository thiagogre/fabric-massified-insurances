package org

import (
	"context"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

// Initialize the orgSetup for the organization.
func Initialize(orgSetup OrgSetup) (*OrgSetup, error) {
	log.Printf("Initializing connection for %s...\n", orgSetup.OrgName)
	clientConnection := orgSetup.newGrpcConnection()
	id := orgSetup.newIdentity()
	sign := orgSetup.newSign()

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
		panic(err)
	}
	orgSetup.Gateway = *gateway

	// Context used for event listening
	ctx, cancel := context.WithCancel(context.Background())
	orgSetup.Context = ctx
	orgSetup.CancelContext = cancel

	log.Println("Initialization complete")

	return &orgSetup, nil
}

// newGrpcConnection creates a gRPC connection to the Gateway server.
func (orgSetup OrgSetup) newGrpcConnection() *grpc.ClientConn {
	certificate, err := loadCertificate(orgSetup.TLSCertPath)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, orgSetup.GatewayPeer)

	connection, err := grpc.NewClient(orgSetup.PeerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}

// newIdentity creates a client identity for this Gateway connection using an X.509 certificate.
func (orgSetup OrgSetup) newIdentity() *identity.X509Identity {
	certificate, err := loadCertificate(orgSetup.CertPath)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(orgSetup.MSPID, certificate)
	if err != nil {
		panic(err)
	}

	return id
}

// newSign creates a function that generates a digital signature from a message digest using a private key.
func (orgSetup OrgSetup) newSign() identity.Sign {
	files, err := os.ReadDir(orgSetup.KeyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key directory: %w", err))
	}
	privateKeyPEM, err := os.ReadFile(path.Join(orgSetup.KeyPath, files[0].Name()))

	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}

func loadCertificate(filename string) (*x509.Certificate, error) {
	certificatePEM, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}
	return identity.CertificateFromPEM(certificatePEM)
}
