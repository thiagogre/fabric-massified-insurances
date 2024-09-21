package adapters

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/constants"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain/mocks"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/tests"
)

func setupTest(t *testing.T) (*mocks.MockClaimServiceInterface, *ClaimHandler, *gomock.Controller) {
	tests.SetupLogger()
	ctrl := gomock.NewController(t)
	mockClaimService := mocks.NewMockClaimServiceInterface(ctrl)
	claimHandler := NewClaimHandler(mockClaimService)
	return mockClaimService, claimHandler, ctrl
}

// THIS SETUP IS REQUIRED 'CAUSE THE TEST CASES THAT USED IT
// GET PATH PARAMETERS WITH MUX PKG
func setupClaimRouterForTest(claimHandler *ClaimHandler) *mux.Router {
	router := mux.NewRouter()
	claimRoutes := router.PathPrefix("/claim").Subrouter()
	claimRoutes.HandleFunc("/evidence/{username}", claimHandler.GetPDFs).Methods(http.MethodGet)
	claimRoutes.HandleFunc("/evidence/{username}/{filename}", claimHandler.ServePDF).Methods(http.MethodGet)
	return router
}

func TestClaimHandler_Execute_SuccessfulUpload(t *testing.T) {
	mockClaimService, claimHandler, ctrl := setupTest(t)
	defer ctrl.Finish()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("files", "test.pdf")
	require.NoError(t, err)
	part.Write([]byte("fake pdf content"))
	err = writer.WriteField("username", "testuser")
	require.NoError(t, err)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/claim", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	mockAsset := &domain.Asset{ID: "123", Insured: "testuser", Evidences: "evidence123"}
	mockClaimService.EXPECT().GetAsset("testuser").Return(mockAsset, nil).AnyTimes()
	mockClaimService.EXPECT().StoreClaim(gomock.Any(), "./uploads/testuser").Return(nil)
	mockClaimService.EXPECT().UpdateAsset(gomock.Any(), "./uploads/testuser").Return(nil)

	claimHandler.Execute(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "Claim in analysis")
}

func TestClaimHandler_Execute_NoFilesUploaded(t *testing.T) {
	_, claimHandler, ctrl := setupTest(t)
	defer ctrl.Finish()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/claim", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	claimHandler.Execute(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "No files uploaded")
}

func TestClaimHandler_Execute_UsernameRequired(t *testing.T) {
	_, claimHandler, ctrl := setupTest(t)
	defer ctrl.Finish()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("files", "test.pdf")
	require.NoError(t, err)
	part.Write([]byte("fake pdf content"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/claim", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	claimHandler.Execute(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "username is required")
}

func TestClaimHandler_Execute_FileTooLarge(t *testing.T) {
	mockClaimService, claimHandler, ctrl := setupTest(t)
	defer ctrl.Finish()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("files", "large_file.pdf")
	require.NoError(t, err)
	part.Write(make([]byte, constants.MaxFileSize+1))
	err = writer.WriteField("username", "testuser")
	require.NoError(t, err)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/claim", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	mockAsset := &domain.Asset{ID: "123", Insured: "testuser", Evidences: "evidence123"}
	mockClaimService.EXPECT().GetAsset("testuser").Return(mockAsset, nil).AnyTimes()
	mockClaimService.EXPECT().UpdateAsset(gomock.Any(), "./uploads/testuser").Return(nil)

	claimHandler.Execute(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "File too large")
}

func TestClaimHandler_Execute_InvalidFileType(t *testing.T) {
	mockClaimService, claimHandler, ctrl := setupTest(t)
	defer ctrl.Finish()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("files", "test.txt")
	require.NoError(t, err)
	part.Write([]byte("fake txt content"))
	err = writer.WriteField("username", "testuser")
	require.NoError(t, err)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/claim", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	mockAsset := &domain.Asset{ID: "123", Insured: "testuser", Evidences: "evidence123"}
	mockClaimService.EXPECT().GetAsset("testuser").Return(mockAsset, nil).AnyTimes()
	mockClaimService.EXPECT().UpdateAsset(gomock.Any(), "./uploads/testuser").Return(nil)

	claimHandler.Execute(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "Invalid file type")
}

func TestClaimHandler_Execute_ErrorFetchingAsset(t *testing.T) {
	mockClaimService, claimHandler, ctrl := setupTest(t)
	defer ctrl.Finish()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("files", "test.pdf")
	require.NoError(t, err)
	part.Write([]byte("fake pdf content"))
	err = writer.WriteField("username", "testuser")
	require.NoError(t, err)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/claim", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	mockClaimService.EXPECT().GetAsset("testuser").Return(nil, errors.New("failed to fetch asset"))

	claimHandler.Execute(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "Error fetching asset")
}

func TestClaimHandler_Execute_ErrorSavingFile(t *testing.T) {
	mockClaimService, claimHandler, ctrl := setupTest(t)
	defer ctrl.Finish()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("files", "test.pdf")
	require.NoError(t, err)
	part.Write([]byte("fake pdf content"))
	err = writer.WriteField("username", "testuser")
	require.NoError(t, err)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/claim", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	mockAsset := &domain.Asset{ID: "123", Insured: "testuser", Evidences: "evidence123"}
	mockClaimService.EXPECT().GetAsset("testuser").Return(mockAsset, nil)
	mockClaimService.EXPECT().StoreClaim(gomock.Any(), "./uploads/testuser").Return(errors.New("unable to save file"))
	mockClaimService.EXPECT().UpdateAsset(gomock.Any(), "./uploads/testuser").Return(nil)

	claimHandler.Execute(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "Unable to save file")
}

func TestClaimHandler_Execute_ErrorUpdatingAsset(t *testing.T) {
	mockClaimService, claimHandler, ctrl := setupTest(t)
	defer ctrl.Finish()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("files", "test.pdf")
	require.NoError(t, err)
	part.Write([]byte("fake pdf content"))
	err = writer.WriteField("username", "testuser")
	require.NoError(t, err)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/claim", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	mockAsset := &domain.Asset{ID: "123", Insured: "testuser", Evidences: "evidence123"}
	mockClaimService.EXPECT().GetAsset("testuser").Return(mockAsset, nil)
	mockClaimService.EXPECT().StoreClaim(gomock.Any(), "./uploads/testuser").Return(nil)
	mockClaimService.EXPECT().UpdateAsset(gomock.Any(), "./uploads/testuser").Return(errors.New("failed to update asset"))

	claimHandler.Execute(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "Error updating asset")
}

func TestClaimHandler_GetPDFs_Success(t *testing.T) {
	mockClaimService, claimHandler, ctrl := setupTest(t)
	defer ctrl.Finish()

	router := setupClaimRouterForTest(claimHandler)

	req := httptest.NewRequest(http.MethodGet, "/claim/evidence/testuser", nil)
	rec := httptest.NewRecorder()

	pdfURLs := []string{"http://localhost/uploads/testuser/file1.pdf", "http://localhost/uploads/testuser/file2.pdf"}
	mockClaimService.EXPECT().ListPDFs("testuser", req.Host).Return(pdfURLs, nil)

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "file1.pdf")
	require.Contains(t, rec.Body.String(), "file2.pdf")
}

func TestClaimHandler_GetPDFs_ErrorListingPDFs(t *testing.T) {
	mockClaimService, claimHandler, ctrl := setupTest(t)
	defer ctrl.Finish()

	router := setupClaimRouterForTest(claimHandler)

	req := httptest.NewRequest(http.MethodGet, "/claim/evidence/testuser", nil)
	rec := httptest.NewRecorder()

	mockClaimService.EXPECT().ListPDFs("testuser", req.Host).Return(nil, errors.New("failed to list PDFs"))

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "Error listing PDFs")
}

func TestClaimHandler_GetPDFs_NoPDFsFound(t *testing.T) {
	mockClaimService, claimHandler, ctrl := setupTest(t)
	defer ctrl.Finish()

	router := setupClaimRouterForTest(claimHandler)

	req := httptest.NewRequest(http.MethodGet, "/claim/evidence/testuser", nil)
	rec := httptest.NewRecorder()

	mockClaimService.EXPECT().ListPDFs("testuser", req.Host).Return([]string{}, nil)

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "[]")
}

func TestClaimHandler_ServePDF_Success(t *testing.T) {
	mockClaimService, claimHandler, ctrl := setupTest(t)
	defer ctrl.Finish()

	router := setupClaimRouterForTest(claimHandler)

	req := httptest.NewRequest(http.MethodGet, "/claim/evidence/testuser/test.pdf", nil)
	rec := httptest.NewRecorder()

	mockClaimService.EXPECT().IsExist("./uploads/testuser/test.pdf").Return(true)

	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestClaimHandler_ServePDF_FileNotFound(t *testing.T) {
	mockClaimService, claimHandler, ctrl := setupTest(t)
	defer ctrl.Finish()

	router := setupClaimRouterForTest(claimHandler)

	req := httptest.NewRequest(http.MethodGet, "/claim/evidence/testuser/test.pdf", nil)
	rec := httptest.NewRecorder()

	mockClaimService.EXPECT().IsExist("./uploads/testuser/test.pdf").Return(false)

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
	require.Contains(t, rec.Body.String(), "File not found")
}
