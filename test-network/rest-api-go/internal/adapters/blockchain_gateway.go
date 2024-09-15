package adapters

import (
	"context"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain"
)

type blockchainGatewayImpl struct {
	gateway *client.Gateway
}

func NewBlockchainGateway(gateway *client.Gateway) *blockchainGatewayImpl {
	return &blockchainGatewayImpl{gateway}
}

type networkImpl struct {
	network *client.Network
}

func NewNetwork(network *client.Network) *networkImpl {
	return &networkImpl{network}
}

type contractImpl struct {
	contract *client.Contract
}

func NewContract(contract *client.Contract) *contractImpl {
	return &contractImpl{contract}
}

func (g *blockchainGatewayImpl) GetNetwork(channelID string) domain.NetworkInterface {
	network := g.gateway.GetNetwork(channelID)
	return &networkImpl{network}
}

func (n *networkImpl) GetContract(chaincodeID string) domain.ContractInterface {
	contract := n.network.GetContract(chaincodeID)
	return &contractImpl{contract}
}

func (n *networkImpl) GetEvents(ctx context.Context, chaincodeID string, txnBlockNumber int) (<-chan *domain.Events, error) {
	events, err := n.network.ChaincodeEvents(ctx, chaincodeID, client.WithStartBlock(uint64(txnBlockNumber)))
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (c *contractImpl) EvaluateTransaction(function string, args ...string) ([]byte, error) {
	return c.contract.EvaluateTransaction(function, args...)
}

func (c *contractImpl) NewTransactionProposal(transactionName string, args ...string) (*domain.TransactionProposalStatus, error) {
	proposal, err := c.contract.NewProposal(transactionName, client.WithArguments(args...))
	if err != nil {
		return nil, err
	}

	transaction, err := proposal.Endorse()
	if err != nil {
		return nil, err
	}

	commit, err := transaction.Submit()
	if err != nil {
		return nil, err
	}

	status, err := commit.Status()
	if err != nil {
		return nil, err
	}

	return status, nil
}
