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

	testCases := []struct {
		name          string
		fileName      string
		fileContent   []byte
		expectedCode  int
		expectedBody  string
		mockReturn    error
		mockStoreCall bool
		hasUsername   bool
	}{
		{
			name:          "successful upload",
			fileName:      "test.pdf",
			fileContent:   []byte("fake pdf content"),
			expectedCode:  http.StatusOK,
			expectedBody:  "All files uploaded successfully",
			mockReturn:    nil,
			mockStoreCall: true,
			hasUsername:   true,
		},
		{
			name:          "no files uploaded",
			fileName:      "",
			fileContent:   make([]byte, 0),
			expectedCode:  http.StatusBadRequest,
			expectedBody:  "No files uploaded",
			mockStoreCall: false,
			hasUsername:   true,
		},
		{
			name:          "username is required",
			fileName:      "test.pdf",
			fileContent:   []byte("fake pdf content"),
			expectedCode:  http.StatusBadRequest,
			expectedBody:  "username is required",
			mockReturn:    nil,
			mockStoreCall: false,
			hasUsername:   false,
		},
		{
			name:          "file too large",
			fileName:      "large_file.pdf",
			fileContent:   make([]byte, constants.MaxFileSize+1),
			expectedCode:  http.StatusBadRequest,
			expectedBody:  "File too large",
			mockStoreCall: false,
			hasUsername:   true,
		},
		{
			name:          "invalid file type",
			fileName:      "test.txt",
			fileContent:   []byte("fake txt content"),
			expectedCode:  http.StatusBadRequest,
			expectedBody:  "Invalid file type",
			mockStoreCall: false,
			hasUsername:   true,
		},
		{
			name:          "error saving file",
			fileName:      "test.pdf",
			fileContent:   []byte("fake pdf content"),
			expectedCode:  http.StatusBadRequest,
			expectedBody:  "Unable to save file",
			mockReturn:    errors.New("unable to save file"),
			mockStoreCall: true,
			hasUsername:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)

			part, err := writer.CreateFormFile("files", tc.fileName)
			require.NoError(t, err)

			part.Write(tc.fileContent)

			if tc.hasUsername {
				err = writer.WriteField("username", "testuser")
				require.NoError(t, err)
			}

			writer.Close()

			req := httptest.NewRequest(http.MethodPost, target, body)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			rec := httptest.NewRecorder()

			if tc.mockStoreCall {
				mockClaimService.EXPECT().StoreClaim(gomock.Any(), "./uploads/testuser").Return(tc.mockReturn)
			}

			claimHandler.UploadEvidences(rec, req)

			require.Equal(t, tc.expectedCode, rec.Code)
			require.Contains(t, rec.Body.String(), tc.expectedBody)
		})
	}
}
