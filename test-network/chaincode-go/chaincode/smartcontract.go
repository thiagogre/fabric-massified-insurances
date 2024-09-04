package chaincode

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple insurance policy asset
// Insert struct field in alphabetic order => to achieve determinism across languages
// golang keeps the order when marshal to json but doesn't order automatically
type Asset struct {
	ClaimStatus    string `json:"ClaimStatus"`    // Status do sinistro (Ex: "Pending", "Approved", "Rejected")
	CoverageAmount int    `json:"CoverageAmount"` // Valor coberto pelo seguro
	ID             string `json:"ID"`             // ID único da apólice
	InsuredItem    string `json:"InsuredItem"`    // Descrição do item segurado (Ex: "Smartphone XYZ")
	Owner          string `json:"Owner"`          // Nome do proprietário da apólice
	Premium        int    `json:"Premium"`        // Valor do prêmio (custo do seguro)
	Term           int    `json:"Term"`           // Prazo do seguro em meses
}

// HistoryQueryResult structure used for returning result of history query
type HistoryQueryResult struct {
	Record    *Asset    `json:"record"`
	TxId      string    `json:"txId"`
	Timestamp time.Time `json:"timestamp"`
	IsDelete  bool      `json:"isDelete"`
}

// DocsResponse structure used for returning docs object instead of array
type DocsResponse struct {
	Docs []interface{} `json:"docs"`
}

// InitLedger adds a base set of insurance policies to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	policies := []Asset{
		{ID: "policy1", Owner: "Alice", InsuredItem: "Smartphone ABC", CoverageAmount: 5000, Premium: 300, Term: 12, ClaimStatus: "Active"},
		{ID: "policy2", Owner: "Bob", InsuredItem: "Smartphone DEF", CoverageAmount: 4000, Premium: 250, Term: 12, ClaimStatus: "Pending"},
		{ID: "policy3", Owner: "Charlie", InsuredItem: "Smartphone GHI", CoverageAmount: 6000, Premium: 350, Term: 12, ClaimStatus: "Approved"},
		{ID: "policy4", Owner: "Diana", InsuredItem: "Smartphone JKL", CoverageAmount: 4500, Premium: 275, Term: 12, ClaimStatus: "Reject"},
	}

	for _, policy := range policies {
		policyBytes, err := json.Marshal(policy)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(policy.ID, policyBytes)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

func (s *SmartContract) readState(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
	assetBytes, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %w", err)
	}
	if len(assetBytes) == 0 {
		return nil, nil
	}

	var asset Asset
	err = json.Unmarshal(assetBytes, &asset)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal asset JSON: %w", err)
	}

	return &asset, nil
}

// CreateAsset issues a new insurance policy asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, id string, owner string, insuredItem string, coverageAmount int, premium int, term int) error {
	exist, err := s.readState(ctx, id)
	if err != nil {
		return err
	}
	if exist != nil {
		return fmt.Errorf("the asset %s already exist", id)
	}

	asset := Asset{
		ID:             id,
		Owner:          owner,
		InsuredItem:    insuredItem,
		CoverageAmount: coverageAmount,
		Premium:        premium,
		Term:           term,
		ClaimStatus:    "Active",
	}

	assetBytes, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	ctx.GetStub().SetEvent("CreateAsset", assetBytes)
	return ctx.GetStub().PutState(id, assetBytes)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
	exist, err := s.readState(ctx, id)
	if err != nil {
		return nil, err
	}
	if exist == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	return exist, nil
}

// UpdateAsset updates an existing insurance policy in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, owner string, insuredItem string, coverageAmount int, premium int, term int, claimStatus string) error {
	exist, err := s.readState(ctx, id)
	if err != nil {
		return err
	}
	if exist == nil {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	asset := Asset{
		ID:             id,
		Owner:          owner,
		InsuredItem:    insuredItem,
		CoverageAmount: coverageAmount,
		Premium:        premium,
		Term:           term,
		ClaimStatus:    claimStatus,
	}

	assetBytes, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	ctx.GetStub().SetEvent("UpdateAsset", assetBytes)
	return ctx.GetStub().PutState(id, assetBytes)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	exist, err := s.readState(ctx, id)
	if err != nil {
		return err
	}
	if exist == nil {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	assetBytes, err := json.Marshal(exist)
	if err != nil {
		return err
	}

	ctx.GetStub().SetEvent("DeleteAsset", assetBytes)
	return ctx.GetStub().DelState(id)
}

// TransferAsset updates the owner field of asset with given id in world state, and returns the old owner.
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) (string, error) {
	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		return "", err
	}

	oldOwner := asset.Owner
	asset.Owner = newOwner

	assetBytes, err := json.Marshal(asset)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(id, assetBytes)
	if err != nil {
		return "", err
	}

	return oldOwner, nil
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) (*DocsResponse, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	records := []interface{}{}
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		records = append(records, &asset)
	}

	docsResponse := DocsResponse{
		Docs: records,
	}

	return &docsResponse, nil
}

// GetAssetRecords returns the chain of custody for an asset since issuance.
func (s *SmartContract) GetAssetRecords(ctx contractapi.TransactionContextInterface, id string) (*DocsResponse, error) {
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(id)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	records := []interface{}{}
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &asset)
			if err != nil {
				return nil, err
			}
		} else {
			asset = Asset{
				ID: id,
			}
		}

		timestampProto := response.Timestamp
		err = timestampProto.CheckValid()
		if err != nil {
			return nil, err
		}
		timestamp := timestampProto.AsTime()

		record := HistoryQueryResult{
			TxId:      response.TxId,
			Timestamp: timestamp,
			Record:    &asset,
			IsDelete:  response.IsDelete,
		}
		records = append(records, record)
	}

	docsResponse := DocsResponse{
		Docs: records,
	}

	return &docsResponse, nil
}
