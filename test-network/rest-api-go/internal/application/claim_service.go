package application

import (
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
