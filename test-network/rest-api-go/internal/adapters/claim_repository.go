package adapters

import (
	"mime/multipart"
	"os"
	"path/filepath"

	"io"

	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/logger"
)

type ClaimRepository struct{}

func NewClaimRepository() *ClaimRepository {
	return &ClaimRepository{}
}

func (r *ClaimRepository) SaveFile(file *multipart.FileHeader, uploadDir, filename string) error {
	err := os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		logger.Error("Error creating upload directory: " + err.Error())
		return err
	}

	out, err := os.Create(filepath.Join(uploadDir, filename))
	if err != nil {
		logger.Error("Error creating file: " + err.Error())
		return err
	}
	defer out.Close()

	src, err := file.Open()
	if err != nil {
		logger.Error("Error opening file: " + err.Error())
		return err
	}
	defer src.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		logger.Error("Error saving file: " + err.Error())
		return err
	}

	return nil
}
