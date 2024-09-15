package application

import (
	"encoding/json"

	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain"
)

type QueryServiceImpl struct {
	BlockchainGateway domain.BlockchainGatewayInterface
}

func NewQueryService(blockchainGateway domain.BlockchainGatewayInterface) *QueryServiceImpl {
	return &QueryServiceImpl{
		BlockchainGateway: blockchainGateway,
	}
}

func (s *QueryServiceImpl) ExecuteQuery(channelID, chaincodeID, function string, args []string) (interface{}, error) {
	contract := s.BlockchainGateway.GetNetwork(channelID).GetContract(chaincodeID)
	txn, err := contract.EvaluateTransaction(function, args...)
	if err != nil {
		return nil, err
	}

	return s.ParseTransaction(txn)
}

// Decode base64-encoded protobuf binary data into JSON
func (s *QueryServiceImpl) ParseTransaction(b []byte) (interface{}, error) {
	var response interface{}
	err := json.Unmarshal(b, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
