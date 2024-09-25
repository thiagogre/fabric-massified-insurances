package adapters

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/constants"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain"
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

	asset, err := h.ClaimService.GetAsset(username, utils.GetFullHostURL(r))
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

	if err := h.ClaimService.UpdateAssetClaimStatus(asset, "Pending", utils.GetFullHostURL(r)); err != nil {
		logger.Error("Error updating asset: " + err.Error())
		utils.ErrorResponse(w, http.StatusBadRequest, "Error updating asset: "+err.Error())
		return
	}

	if errorMessage != "" {
		logger.Error(errorMessage)
		utils.ErrorResponse(w, http.StatusBadRequest, errorMessage)
		return
	}

	response := domain.SuccessResponse[string]{Success: true, Data: "Claim in analysis"}
	logger.Success(response)
	utils.SuccessResponse(w, http.StatusOK, response)
}

func (h *ClaimHandler) GetPDFs(w http.ResponseWriter, r *http.Request) {
	logger.Info("Received a request")

	vars := mux.Vars(r)
	username := vars["username"]

	pdfURLs, err := h.ClaimService.ListPDFs(username, utils.GetFullHostURL(r))
	if err != nil {
		logger.Error("Error listing PDFs: " + err.Error())
		utils.ErrorResponse(w, http.StatusBadRequest, "Error listing PDFs: "+err.Error())
		return
	}

	response := domain.SuccessResponse[[]string]{Success: true, Data: pdfURLs}
	logger.Success(response)
	utils.SuccessResponse(w, http.StatusOK, response)
}

func (h *ClaimHandler) ServePDF(w http.ResponseWriter, r *http.Request) {
	logger.Info("Received a request")

	vars := mux.Vars(r)
	username := vars["username"]
	filename := vars["filename"]

	filePath := fmt.Sprintf("%s/%s/%s", constants.DefaultUploadDir, username, filename)

	if !h.ClaimService.IsExist(filePath) {
		logger.Error("File not found")
		utils.ErrorResponse(w, http.StatusNotFound, "File not found")
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.WriteHeader(http.StatusOK)
	http.ServeFile(w, r, filePath)
}

func (h *ClaimHandler) Validate(w http.ResponseWriter, r *http.Request) {
	logger.Info("Received a request")

	var body domain.ClaimValidateRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		logger.Error("Failed to parse request body" + err.Error())
		utils.ErrorResponse(w, http.StatusBadRequest, "Failed to parse request body")
		return
	}
	logger.Info(body)

	asset, err := h.ClaimService.GetAsset(body.Username, utils.GetFullHostURL(r))
	if err != nil {
		logger.Error("Error fetching asset: " + err.Error())
		utils.ErrorResponse(w, http.StatusBadRequest, "Error fetching asset: "+err.Error())
		return
	}

	var newClaimStatus string
	if body.IsApproved {
		newClaimStatus = "EvidencesApproved"
	} else {
		newClaimStatus = "EvidencesRejected"
	}

	if err := h.ClaimService.UpdateAssetClaimStatus(asset, newClaimStatus, utils.GetFullHostURL(r)); err != nil {
		logger.Error("Error updating asset: " + err.Error())
		utils.ErrorResponse(w, http.StatusBadRequest, "Error updating asset: "+err.Error())
		return
	}

	response := domain.SuccessResponse[domain.ClaimValidateRequest]{Success: true, Data: body}
	logger.Success(body)
	utils.SuccessResponse(w, http.StatusOK, response)
}

func (h *ClaimHandler) Finish(w http.ResponseWriter, r *http.Request) {
	logger.Info("Received a request")

	var body domain.ClaimValidateRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		logger.Error("Failed to parse request body" + err.Error())
		utils.ErrorResponse(w, http.StatusBadRequest, "Failed to parse request body")
		return
	}
	logger.Info(body)

	asset, err := h.ClaimService.GetAsset(body.Username, utils.GetFullHostURL(r))
	if err != nil {
		logger.Error("Error fetching asset: " + err.Error())
		utils.ErrorResponse(w, http.StatusBadRequest, "Error fetching asset: "+err.Error())
		return
	}

	var newClaimStatus string
	if body.IsApproved {
		newClaimStatus = "Approved"
	} else {
		newClaimStatus = "Rejected"
	}

	if err := h.ClaimService.UpdateAssetClaimStatus(asset, newClaimStatus, utils.GetFullHostURL(r)); err != nil {
		logger.Error("Error updating asset: " + err.Error())
		utils.ErrorResponse(w, http.StatusBadRequest, "Error updating asset: "+err.Error())
		return
	}

	response := domain.SuccessResponse[domain.ClaimValidateRequest]{Success: true, Data: body}
	logger.Success(body)
	utils.SuccessResponse(w, http.StatusOK, response)
}
