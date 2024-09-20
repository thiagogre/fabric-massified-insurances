package chaincode

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ABAC (claim_analyst, partner, evidence_analyst)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple insurance policy asset
// Insert struct field in alphabetic order => to achieve determinism across languages
// golang keeps the order when marshal to json but doesn't order automatically
type Asset struct {
	ClaimStatus      string `json:"ClaimStatus"`      // Status do sinistro (Ex: "Active", "Pending", "Approved", "Rejected")
	CoverageDuration int    `json:"CoverageDuration"` // Prazo do seguro em meses
	CoverageType     int    `json:"CoverageType"`     // 0 = Cobertura contra roubo/furto de celular
	CoverageValue    int    `json:"CoverageValue"`    // Valor coberto pelo seguro
	Evidences        string `json:"Evidences"`        // Evidências
	ID               string `json:"ID"`               // ID único da apólice
	Insured          string `json:"Insured"`          // Proprietário da apólice
	Partner          string `json:"Partner"`          // Parceiro de distribuição
	Premium          int    `json:"Premium"`          // Valor do prêmio (custo do seguro)
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
		{ID: "policy1", Insured: "Alice", Evidences: "", CoverageDuration: 12, CoverageType: 0, CoverageValue: 5000, Premium: 300, Partner: "Parceiro varejista", ClaimStatus: "Active"},
		{ID: "policy2", Insured: "Bob", Evidences: "", CoverageDuration: 12, CoverageType: 0, CoverageValue: 4000, Premium: 250, Partner: "Parceiro varejista", ClaimStatus: "Pending"},
		{ID: "policy3", Insured: "Charlie", Evidences: "", CoverageDuration: 12, CoverageType: 0, CoverageValue: 6000, Premium: 350, Partner: "Parceiro varejista", ClaimStatus: "Approved"},
		{ID: "policy4", Insured: "Diana", Evidences: "", CoverageDuration: 12, CoverageType: 0, CoverageValue: 4500, Premium: 275, Partner: "Parceiro varejista", ClaimStatus: "Reject"},
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
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, id string, insured string, coverageDuration int, converageValue int, coverageType int, partner string, premium int) error {
	// if err := validateABAC(ctx, []string{"abac.partner"}); err != nil {
	// 	return err
	// }

	asset, err := s.readState(ctx, id)
	if err != nil {
		return err
	}
	if asset != nil {
		return fmt.Errorf("the asset %s already exist", id)
	}

	// clientID, err := s.getSubmittingClientIdentity(ctx)
	// if err != nil {
	// 	return err
	// }

	newAsset := Asset{
		ID:               id,
		Insured:          insured,
		Partner:          partner,
		CoverageDuration: coverageDuration,
		CoverageType:     coverageType,
		CoverageValue:    converageValue,
		Evidences:        "",
		Premium:          premium,
		ClaimStatus:      "Active",
	}

	newAssetBytes, err := json.Marshal(newAsset)
	if err != nil {
		return err
	}

	ctx.GetStub().SetEvent("CreateAsset", newAssetBytes)
	return ctx.GetStub().PutState(id, newAssetBytes)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
	asset, err := s.readState(ctx, id)
	if err != nil {
		return nil, err
	}
	if asset == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	return asset, nil
}

// UpdateAsset updates an existing insurance policy in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, insured string, coverageDuration int, converageValue int, coverageType int, partner string, premium int, claimStatus string, evidences string) error {
	// if err := validateABAC(ctx, []string{"abac.evidence_analyst"}); err != nil {
	// 	return err
	// }

	asset, err := s.readState(ctx, id)
	if err != nil {
		return err
	}
	if asset == nil {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	newAsset := Asset{
		ID:               id,
		Insured:          insured,
		Partner:          partner,
		CoverageDuration: coverageDuration,
		CoverageType:     coverageType,
		CoverageValue:    converageValue,
		Evidences:        evidences,
		Premium:          premium,
		ClaimStatus:      claimStatus,
	}

	newAssetBytes, err := json.Marshal(newAsset)
	if err != nil {
		return err
	}

	ctx.GetStub().SetEvent("UpdateAsset", newAssetBytes)
	return ctx.GetStub().PutState(id, newAssetBytes)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	asset, err := s.readState(ctx, id)
	if err != nil {
		return err
	}
	if asset == nil {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	assetBytes, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	ctx.GetStub().SetEvent("DeleteAsset", assetBytes)
	return ctx.GetStub().DelState(id)
}

// TransferAsset updates the insured field of asset with given id in world state, and returns the old insured.
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newInsured string) (string, error) {
	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		return "", err
	}

	oldInsured := asset.Insured
	asset.Insured = newInsured

	assetBytes, err := json.Marshal(asset)
	if err != nil {
		return "", err
	}

	err = ctx.GetStub().PutState(id, assetBytes)
	if err != nil {
		return "", err
	}

	return oldInsured, nil
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

// GetAssetsByRichQuery return all assets given a rich query
func (s *SmartContract) GetAssetsByRichQuery(ctx contractapi.TransactionContextInterface, queryString string) (*DocsResponse, error) {
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	records := []interface{}{}
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		var asset Asset
		err = json.Unmarshal(queryResult.Value, &asset)
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

// GetSubmittingClientIdentity returns the name and issuer of the identity that
// invokes the smart contract. This function base64 decodes the identity string
// before returning the value to the client or smart contract.
func (s *SmartContract) getSubmittingClientIdentity(ctx contractapi.TransactionContextInterface) (string, error) {
	b64ID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("failed to read clientID: %v", err)
	}

	decodeID, err := base64.StdEncoding.DecodeString(b64ID)
	if err != nil {
		return "", fmt.Errorf("failed to base64 decode clientID: %v", err)
	}

	return string(decodeID), nil
}

// validateABAC checks whether the client identity has at least one of the
// required roles to perform an action.
func validateABAC(ctx contractapi.TransactionContextInterface, roles []string) error {
	authorized := false
	for _, role := range roles {
		err := ctx.GetClientIdentity().AssertAttributeValue(role, "true")
		if err == nil {
			authorized = true
			break
		}
	}

	if !authorized {
		return fmt.Errorf("submitting client not authorized to perform this action, does not have the required role %v", roles)
	}

	return nil
}
