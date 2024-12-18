package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/constants"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain"
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

func (r *ClaimRepository) IsFileOrDirExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func (r *ClaimRepository) ListPDFFiles(username string) ([]string, error) {
	folderPath := fmt.Sprintf("%s/%s", constants.DefaultUploadDir, username)

	if !r.IsFileOrDirExist(folderPath) {
		return nil, fmt.Errorf("folder not found")
	}

	files, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	var pdfFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".pdf") {
			pdfFiles = append(pdfFiles, file.Name())
		}
	}

	if len(pdfFiles) == 0 {
		return nil, fmt.Errorf("no PDF files found")
	}

	return pdfFiles, nil
}

func (r *ClaimRepository) GetAsset(username string, host string) (*domain.Asset, error) {
	URL := fmt.Sprintf(
		"%s/smartcontract/query?channelid=%s&chaincodeid=%s&function=GetAssetsByRichQuery&args=%s",
		host,
		constants.ChannelID,
		constants.ChaincodeID,
		fmt.Sprintf(`{"selector":{"Insured":"%s"}}`, username),
	)

	resp, err := http.Get(URL)
	if err != nil {
		logger.Error("Failed to call API: " + err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error(fmt.Sprintf("API request failed with status: %d", resp.StatusCode))
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var response domain.SuccessResponse[domain.DocsResponse[domain.Asset]]
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		logger.Error("Failed to decode response: " + err.Error())
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	if len(response.Data.Docs) == 0 {
		logger.Error("No assets found in response")
		return nil, fmt.Errorf("no assets found for user: %s", username)
	}

	return &response.Data.Docs[0], nil
}

func (r *ClaimRepository) UpdateAsset(body *domain.InvokeRequest, host string) error {
	url := fmt.Sprintf("%s/smartcontract/invoke", host)

	jsonBody, err := json.Marshal(body)
	if err != nil {
		logger.Error("Failed to marshal JSON: " + err.Error())
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		logger.Error("Failed to call API: " + err.Error())
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error(fmt.Sprintf("API request failed with status: %d", resp.StatusCode))
		return fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	return nil
}
