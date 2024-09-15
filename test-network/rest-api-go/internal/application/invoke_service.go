package application

import (
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain"
)

type InvokeServiceImpl struct {
	BlockchainGateway domain.BlockchainGatewayInterface
}

func NewInvokeService(blockchainGateway domain.BlockchainGatewayInterface) *InvokeServiceImpl {
	return &InvokeServiceImpl{
		BlockchainGateway: blockchainGateway,
	}
}

func (s *InvokeServiceImpl) ExecuteInvoke(channelID, chaincodeID, function string, args []string) (*domain.TransactionProposalStatus, error) {
	contract := s.BlockchainGateway.GetNetwork(channelID).GetContract(chaincodeID)
	txn, err := contract.NewTransactionProposal(function, args...)
	if err != nil {
		return nil, err
	}

	return txn, nil
}
