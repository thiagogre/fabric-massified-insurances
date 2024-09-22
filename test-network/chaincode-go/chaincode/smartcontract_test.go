package chaincode_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	"github.com/stretchr/testify/require"

	"github.com/thiagogre/fabric-massified-insurances/test-network/chaincode-go/chaincode"
	"github.com/thiagogre/fabric-massified-insurances/test-network/chaincode-go/chaincode/mocks"
)

//go:generate counterfeiter -o mocks/transaction.go -fake-name TransactionContext . transactionContext
type transactionContext interface {
	contractapi.TransactionContextInterface
}

//go:generate counterfeiter -o mocks/chaincodestub.go -fake-name ChaincodeStub . chaincodeStub
type chaincodeStub interface {
	shim.ChaincodeStubInterface
}

//go:generate counterfeiter -o mocks/statequeryiterator.go -fake-name StateQueryIterator . stateQueryIterator
type stateQueryIterator interface {
	shim.StateQueryIteratorInterface
}

func TestInitLedger(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	asset := chaincode.SmartContract{}
	err := asset.InitLedger(transactionContext)
	require.NoError(t, err)

	chaincodeStub.PutStateReturns(fmt.Errorf("failed inserting key"))
	err = asset.InitLedger(transactionContext)
	require.EqualError(t, err, "failed to put to world state. failed inserting key")
}

func TestCreateAsset(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	asset := chaincode.SmartContract{}
	err := asset.CreateAsset(transactionContext, "", "", 0, 0, 0, "", 0)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns([]byte{}, nil)
	err = asset.CreateAsset(transactionContext, "policy1", "", 0, 0, 0, "", 0)
	require.NoError(t, err)

	expectedAsset := &chaincode.Asset{ID: "policy1"}
	bytes, err := json.Marshal(expectedAsset)
	require.NoError(t, err)
	chaincodeStub.GetStateReturns(bytes, nil)
	err = asset.CreateAsset(transactionContext, "policy1", "", 0, 0, 0, "", 0)
	require.EqualError(t, err, "the asset policy1 already exist")

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	err = asset.CreateAsset(transactionContext, "policy1", "", 0, 0, 0, "", 0)
	require.EqualError(t, err, "failed to read from world state: unable to retrieve asset")
}

func TestReadAsset(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	expectedAsset := &chaincode.Asset{ID: "policy1"}
	bytes, err := json.Marshal(expectedAsset)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(bytes, nil)
	asset := chaincode.SmartContract{}
	a, err := asset.ReadAsset(transactionContext, "")
	require.NoError(t, err)
	require.Equal(t, expectedAsset, a)

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	_, err = asset.ReadAsset(transactionContext, "")
	require.EqualError(t, err, "failed to read from world state: unable to retrieve asset")

	chaincodeStub.GetStateReturns(nil, nil)
	a, err = asset.ReadAsset(transactionContext, "policy1")
	require.EqualError(t, err, "the asset policy1 does not exist")
	require.Nil(t, a)
}

func TestUpdateAsset(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	expectedAsset := &chaincode.Asset{ID: "policy1"}
	bytes, err := json.Marshal(expectedAsset)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(bytes, nil)
	asset := chaincode.SmartContract{}
	err = asset.UpdateAsset(transactionContext, "", "", 0, 0, 0, "", 0, "")
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(nil, nil)
	err = asset.UpdateAsset(transactionContext, "policy1", "", 0, 0, 0, "", 0, "")
	require.EqualError(t, err, "the asset policy1 does not exist")

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	err = asset.UpdateAsset(transactionContext, "policy1", "", 0, 0, 0, "", 0, "")
	require.EqualError(t, err, "failed to read from world state: unable to retrieve asset")
}

func TestDeleteAsset(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	a := &chaincode.Asset{ID: "policy1"}
	bytes, err := json.Marshal(a)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(bytes, nil)
	chaincodeStub.DelStateReturns(nil)
	asset := chaincode.SmartContract{}
	err = asset.DeleteAsset(transactionContext, "")
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(nil, nil)
	err = asset.DeleteAsset(transactionContext, "policy1")
	require.EqualError(t, err, "the asset policy1 does not exist")

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	err = asset.DeleteAsset(transactionContext, "")
	require.EqualError(t, err, "failed to read from world state: unable to retrieve asset")
}

func TestTransferAsset(t *testing.T) {
	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	a := &chaincode.Asset{ID: "policy1"}
	bytes, err := json.Marshal(a)
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(bytes, nil)
	asset := chaincode.SmartContract{}
	_, err = asset.TransferAsset(transactionContext, "", "")
	require.NoError(t, err)

	chaincodeStub.GetStateReturns(nil, fmt.Errorf("unable to retrieve asset"))
	_, err = asset.TransferAsset(transactionContext, "", "")
	require.EqualError(t, err, "failed to read from world state: unable to retrieve asset")
}

func TestGetAllAssets(t *testing.T) {
	a := &chaincode.Asset{ID: "policy1"}
	bytes, err := json.Marshal(a)
	require.NoError(t, err)

	iterator := &mocks.StateQueryIterator{}
	iterator.HasNextReturnsOnCall(0, true)
	iterator.HasNextReturnsOnCall(1, false)
	iterator.NextReturns(&queryresult.KV{Value: bytes}, nil)

	chaincodeStub := &mocks.ChaincodeStub{}
	transactionContext := &mocks.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	chaincodeStub.GetStateByRangeReturns(iterator, nil)
	asset := &chaincode.SmartContract{}
	assets, err := asset.GetAllAssets(transactionContext)
	require.NoError(t, err)
	require.Equal(t, []interface{}{a}, assets.Docs)

	iterator.HasNextReturns(true)
	iterator.NextReturns(nil, fmt.Errorf("failed retrieving next item"))
	assets, err = asset.GetAllAssets(transactionContext)
	require.EqualError(t, err, "failed retrieving next item")
	require.Nil(t, assets)

	chaincodeStub.GetStateByRangeReturns(nil, fmt.Errorf("failed retrieving all assets"))
	assets, err = asset.GetAllAssets(transactionContext)
	require.EqualError(t, err, "failed retrieving all assets")
	require.Nil(t, assets)
}
