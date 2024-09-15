package domain

import (
	"context"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

type BlockchainGatewayInterface interface {
	GetNetwork(channelID string) NetworkInterface
}

type NetworkInterface interface {
	GetContract(chaincodeID string) ContractInterface
	GetEvents(ctx context.Context, chaincodeID string, txnBlockNumber int) (<-chan *Events, error)
}

type ContractInterface interface {
	EvaluateTransaction(function string, args ...string) ([]byte, error)
	NewTransactionProposal(transactionName string, args ...string) (*TransactionProposalStatus, error)
}

type TransactionProposalStatus = client.Status

type Events = client.ChaincodeEvent
