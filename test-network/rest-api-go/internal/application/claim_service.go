package application

import (
	"fmt"
	"mime/multipart"
	"path/filepath"

	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/logger"
)

type ClaimService struct {
	ClaimRepository domain.ClaimRepositoryInterface
}

func NewClaimService(claimRepository domain.ClaimRepositoryInterface) *ClaimService {
	return &ClaimService{ClaimRepository: claimRepository}
}

func (s *ClaimService) StoreClaim(file *multipart.FileHeader, uploadDir string) error {
	filename := filepath.Base(file.Filename)
	err := s.ClaimRepository.SaveFile(file, uploadDir, filename)
	if err != nil {
		logger.Error("Error storing file: " + err.Error())
		return err
	}
	return nil
}

func (s *ClaimService) ListPDFs(username, host string) ([]string, error) {
	pdfFiles, err := s.ClaimRepository.ListPDFFiles(username)
	if err != nil {
		return nil, err
	}

	var pdfURLs []string
	for _, pdfFile := range pdfFiles {
		pdfURL := fmt.Sprintf("%s/claim/evidence/%s/%s", host, username, pdfFile)
		pdfURLs = append(pdfURLs, pdfURL)
	}

	return pdfURLs, nil
}

func (s *ClaimService) IsExist(filePath string) bool {
	return s.ClaimRepository.IsFileOrDirExist(filePath)
}

func (s *ClaimService) GetAsset(username string) (*domain.Asset, error) {
	return s.ClaimRepository.GetAsset(username)
}

func (s *ClaimService) UpdateAsset(asset *domain.Asset, uploadDir string) error {
	return s.ClaimRepository.UpdateAsset(asset, uploadDir)
}
