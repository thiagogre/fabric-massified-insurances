package domain

import (
	"mime/multipart"
)

type Asset struct {
	ClaimStatus      string `json:"ClaimStatus"`
	CoverageDuration int    `json:"CoverageDuration"`
	CoverageType     int    `json:"CoverageType"`
	CoverageValue    int    `json:"CoverageValue"`
	Evidences        string `json:"Evidences"`
	ID               string `json:"ID"`
	Insured          string `json:"Insured"`
	Partner          string `json:"Partner"`
	Premium          int    `json:"Premium"`
}

type ClaimServiceInterface interface {
	StoreClaim(file *multipart.FileHeader, uploadDir string) error
	GetAsset(username string) (*Asset, error)
	UpdateAsset(asset *Asset, uploadDir string) error
}

type ClaimRepositoryInterface interface {
	SaveFile(file *multipart.FileHeader, uploadDir, filename string) error
	GetAsset(username string) (*Asset, error)
	UpdateAsset(asset *Asset, uploadDir string) error
}
