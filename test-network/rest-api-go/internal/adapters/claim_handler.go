package adapters

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/constants"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/dto"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/logger"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/utils"
)

type ClaimHandler struct {
	ClaimService domain.ClaimServiceInterface
}

func NewClaimHandler(claimService domain.ClaimServiceInterface) *ClaimHandler {
	return &ClaimHandler{ClaimService: claimService}
}

func (h *ClaimHandler) Execute(w http.ResponseWriter, r *http.Request) {
	logger.Info("Received a request")

	err := r.ParseMultipartForm(constants.MaxFileSize)
	if err != nil {
		logger.Error("Error parsing form: " + err.Error())
		utils.ErrorResponse(w, http.StatusBadRequest, "Error parsing form")
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		logger.Error("No files uploaded")
		utils.ErrorResponse(w, http.StatusBadRequest, "No files uploaded")
		return
	}

	username := ""
	if r.MultipartForm.Value["username"] != nil && r.MultipartForm.Value["username"][0] != "" {
		username = r.MultipartForm.Value["username"][0]
	} else {
		logger.Error("username is required")
		utils.ErrorResponse(w, http.StatusBadRequest, "username is required")
		return
	}

	uploadDir := constants.DefaultUploadDir + "/" + username

	asset, err := h.ClaimService.GetAsset(username)
	if err != nil {
		logger.Error("Error fetching asset: " + err.Error())
		utils.ErrorResponse(w, http.StatusBadRequest, "Error fetching asset: "+err.Error())
		return
	}

	var errorMessage string
	for _, fileHeader := range files {
		if fileHeader.Size > constants.MaxFileSize {
			errorMessage = fmt.Sprintf("File too large (%d > %d bytes): %s", fileHeader.Size, constants.MaxFileSize, fileHeader.Filename)
			logger.Error(errorMessage)
			break
		}

		file, err := fileHeader.Open()
		if err != nil {
			errorMessage = fmt.Sprintf("Error opening file: %s, %v", fileHeader.Filename, err)
			logger.Error(errorMessage)
			break
		}
		defer file.Close()

		if filepath.Ext(fileHeader.Filename) != ".pdf" {
			errorMessage = fmt.Sprintf("Invalid file type for file: %s", fileHeader.Filename)
			logger.Error(errorMessage)
			break
		}

		err = h.ClaimService.StoreClaim(fileHeader, uploadDir)
		if err != nil {
			errorMessage = fmt.Sprintf("Unable to save file: %s, %v", fileHeader.Filename, err)
			logger.Error(errorMessage)
			break
		}
	}

	if err := h.ClaimService.UpdateAsset(asset, uploadDir); err != nil {
		logger.Error("Error updating asset: " + err.Error())
		utils.ErrorResponse(w, http.StatusBadRequest, "Error updating asset: "+err.Error())
		return
	}

	if errorMessage != "" {
		logger.Error(errorMessage)
		utils.ErrorResponse(w, http.StatusBadRequest, errorMessage)
		return
	}

	response := dto.SuccessResponse[string]{Success: true, Data: "Claim in analysis"}
	logger.Success(response)
	utils.SuccessResponse(w, http.StatusOK, response)
}
