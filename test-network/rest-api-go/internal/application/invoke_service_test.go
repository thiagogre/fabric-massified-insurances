package application

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain/mocks"
)

func TestExecuteInvoke_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGateway := mocks.NewMockBlockchainGatewayInterface(ctrl)
	mockNetwork := mocks.NewMockNetworkInterface(ctrl)
	mockContract := mocks.NewMockContractInterface(ctrl)

	expectedTxn := &domain.TransactionProposalStatus{}

	mockGateway.EXPECT().GetNetwork("testChannel").Return(mockNetwork)
	mockNetwork.EXPECT().GetContract("testChaincode").Return(mockContract)
	mockContract.EXPECT().NewTransactionProposal("testFunction", "arg1", "arg2").Return(expectedTxn, nil)

	service := NewInvokeService(mockGateway)
	result, err := service.ExecuteInvoke("testChannel", "testChaincode", "testFunction", []string{"arg1", "arg2"})

	require.NoError(t, err)
	require.Equal(t, expectedTxn, result)
}

func TestExecuteInvoke_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGateway := mocks.NewMockBlockchainGatewayInterface(ctrl)
	mockNetwork := mocks.NewMockNetworkInterface(ctrl)
	mockContract := mocks.NewMockContractInterface(ctrl)

	mockGateway.EXPECT().GetNetwork("testChannel").Return(mockNetwork)
	mockNetwork.EXPECT().GetContract("testChaincode").Return(mockContract)
	mockContract.EXPECT().NewTransactionProposal("testFunction", "arg1", "arg2").Return(nil, errors.New("transaction proposal failed"))

	service := NewInvokeService(mockGateway)
	result, err := service.ExecuteInvoke("testChannel", "testChaincode", "testFunction", []string{"arg1", "arg2"})

	require.Error(t, err)
	require.Nil(t, result)
}
