package adapters

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/constants"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain/mocks"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/tests"
)

func TestClaimHandler_UploadClaim(t *testing.T) {
	tests.SetupLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	target := "/claim/evidence/upload"
	mockClaimService := mocks.NewMockClaimServiceInterface(ctrl)
	claimHandler := NewClaimHandler(mockClaimService)

	t.Run("successful upload", func(t *testing.T) {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("files", "test.pdf")
		require.NoError(t, err)
		part.Write([]byte("fake pdf content"))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, target, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()

		mockClaimService.EXPECT().StoreClaim(gomock.Any()).Return(nil)

		claimHandler.UploadEvidences(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), "All files uploaded successfully")
	})

	t.Run("file too large", func(t *testing.T) {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("files", "large_file.pdf")
		require.NoError(t, err)
		part.Write(make([]byte, constants.MaxFileSize+1))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, target, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()

		claimHandler.UploadEvidences(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Contains(t, rec.Body.String(), "File too large")
	})

	t.Run("invalid file type", func(t *testing.T) {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("files", "test.txt")
		require.NoError(t, err)
		part.Write([]byte("fake txt content"))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, target, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()

		claimHandler.UploadEvidences(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Contains(t, rec.Body.String(), "Invalid file type")
	})

	t.Run("error opening file", func(t *testing.T) {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		_, err := writer.CreateFormFile("files", "test.pdf")
		require.NoError(t, err)
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, target, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()

		mockClaimService.EXPECT().StoreClaim(gomock.Any()).Return(errors.New("error opening file"))

		claimHandler.UploadEvidences(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Contains(t, rec.Body.String(), "error opening file")
	})

	t.Run("error saving file", func(t *testing.T) {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("files", "test.pdf")
		require.NoError(t, err)
		part.Write([]byte("fake pdf content"))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, target, body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		rec := httptest.NewRecorder()

		mockClaimService.EXPECT().StoreClaim(gomock.Any()).Return(errors.New("unable to save file"))

		claimHandler.UploadEvidences(rec, req)

		require.Equal(t, http.StatusBadRequest, rec.Code)
		require.Contains(t, rec.Body.String(), "Unable to save file")
	})
}
