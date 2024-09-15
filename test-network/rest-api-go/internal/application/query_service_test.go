package application

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain/mocks"
)

func TestExecuteQuery_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGateway := mocks.NewMockBlockchainGatewayInterface(ctrl)
	mockNetwork := mocks.NewMockNetworkInterface(ctrl)
	mockContract := mocks.NewMockContractInterface(ctrl)

	mockGateway.EXPECT().GetNetwork("testChannel").Return(mockNetwork)
	mockNetwork.EXPECT().GetContract("testChaincode").Return(mockContract)
	mockContract.EXPECT().EvaluateTransaction("testFunction", "arg1", "arg2").Return([]byte(`{"key": "value"}`), nil)

	service := NewQueryService(mockGateway)
	result, err := service.ExecuteQuery("testChannel", "testChaincode", "testFunction", []string{"arg1", "arg2"})

	require.NoError(t, err)
	require.Equal(t, map[string]interface{}{"key": "value"}, result)
}

func TestExecuteQuery_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Using the mocks you previously generated
	mockGateway := mocks.NewMockBlockchainGatewayInterface(ctrl)
	mockNetwork := mocks.NewMockNetworkInterface(ctrl)
	mockContract := mocks.NewMockContractInterface(ctrl)

	mockGateway.EXPECT().GetNetwork("testChannel").Return(mockNetwork)
	mockNetwork.EXPECT().GetContract("testChaincode").Return(mockContract)
	mockContract.EXPECT().EvaluateTransaction("testFunction", "arg1", "arg2").Return(nil, errors.New("transaction failed"))

	service := NewQueryService(mockGateway)
	result, err := service.ExecuteQuery("testChannel", "testChaincode", "testFunction", []string{"arg1", "arg2"})

	require.Error(t, err)
	require.Nil(t, result)
}

func TestParseTransaction(t *testing.T) {
	service := NewQueryService(nil)
	data := []byte(`{"key": "value"}`)

	result, err := service.ParseTransaction(data)

	require.NoError(t, err)
	require.Equal(t, map[string]interface{}{"key": "value"}, result)
}

func TestParseTransaction_InvalidJSON(t *testing.T) {
	service := NewQueryService(nil)
	data := []byte(`{key: value}`) // invalid JSON

	result, err := service.ParseTransaction(data)

	require.Error(t, err)
	require.Nil(t, result)
}
