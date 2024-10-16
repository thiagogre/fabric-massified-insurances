package domain

import (
	"mime/multipart"
)

type ClaimValidateRequest struct {
	Username   string `json:"username"`
	IsApproved bool   `json:"isApproved"`
}

type Asset struct {
	ClaimStatus      string `json:"ClaimStatus"`
	CoverageDuration int    `json:"CoverageDuration"`
	CoverageType     int    `json:"CoverageType"`
	CoverageValue    int    `json:"CoverageValue"`
	ID               string `json:"ID"`
	Insured          string `json:"Insured"`
	Partner          string `json:"Partner"`
	Premium          int    `json:"Premium"`
}

type ClaimServiceInterface interface {
	StoreClaim(file *multipart.FileHeader, uploadDir string) error
	ListPDFs(username, host string) ([]string, error)
	IsExist(filePath string) bool
	GetAsset(username string, host string) (*Asset, error)
	UpdateAssetClaimStatus(asset *Asset, newClaimStatus string, host string) error
}

type ClaimRepositoryInterface interface {
	SaveFile(file *multipart.FileHeader, uploadDir, filename string) error
	ListPDFFiles(username string) ([]string, error)
	IsFileOrDirExist(path string) bool
	GetAsset(username string, host string) (*Asset, error)
	UpdateAsset(body *InvokeRequest, host string) error
}
