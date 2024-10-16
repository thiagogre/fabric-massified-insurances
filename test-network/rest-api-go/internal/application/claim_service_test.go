package application

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/constants"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain/mocks"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/tests"
)

func TestStoreClaim_Success(t *testing.T) {
	tests.SetupLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockClaimRepositoryInterface(ctrl)
	service := NewClaimService(mockRepo)

	fileHeader := tests.CreateTestFileHeader("test.pdf")

	mockRepo.EXPECT().
		SaveFile(fileHeader, "/test-uploads", "test.pdf").
		Return(nil).
		Times(1)

	err := service.StoreClaim(fileHeader, "/test-uploads")
	require.NoError(t, err)
}

func TestStoreClaim_ErrorSavingFile(t *testing.T) {
	tests.SetupLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockClaimRepositoryInterface(ctrl)
	service := NewClaimService(mockRepo)

	fileHeader := tests.CreateTestFileHeader("test.pdf")

	mockRepo.EXPECT().
		SaveFile(fileHeader, "/test-uploads", "test.pdf").
		Return(fmt.Errorf("mock save file error")).
		Times(1)

	err := service.StoreClaim(fileHeader, "/test-uploads")
	require.Error(t, err)
	require.EqualError(t, err, "mock save file error")
}

func TestListPDFs_Success(t *testing.T) {
	tests.SetupLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockClaimRepositoryInterface(ctrl)
	service := NewClaimService(mockRepo)

	username := "testuser"
	host := "http://localhost"

	mockRepo.EXPECT().
		ListPDFFiles(username).
		Return([]string{"file1.pdf", "file2.pdf"}, nil).
		Times(1)

	pdfURLs, err := service.ListPDFs(username, host)
	require.NoError(t, err)
	require.Equal(t, []string{
		"http://localhost/claim/evidence/testuser/file1.pdf",
		"http://localhost/claim/evidence/testuser/file2.pdf",
	}, pdfURLs)
}

func TestListPDFs_ErrorReadingFiles(t *testing.T) {
	tests.SetupLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockClaimRepositoryInterface(ctrl)
	service := NewClaimService(mockRepo)

	username := "testuser"
	host := "http://localhost"

	mockRepo.EXPECT().
		ListPDFFiles(username).
		Return(nil, fmt.Errorf("mock error listing files")).
		Times(1)

	pdfURLs, err := service.ListPDFs(username, host)
	require.Error(t, err)
	require.Nil(t, pdfURLs)
	require.EqualError(t, err, "mock error listing files")
}

func TestIsExist_FileExists(t *testing.T) {
	tests.SetupLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockClaimRepositoryInterface(ctrl)
	service := NewClaimService(mockRepo)

	filePath := "/uploads/testuser/test.pdf"

	mockRepo.EXPECT().
		IsFileOrDirExist(filePath).
		Return(true).
		Times(1)

	exists := service.IsExist(filePath)
	require.True(t, exists)
}

func TestIsExist_FileDoesNotExist(t *testing.T) {
	tests.SetupLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockClaimRepositoryInterface(ctrl)
	service := NewClaimService(mockRepo)

	filePath := "/uploads/testuser/nonexistent.pdf"

	mockRepo.EXPECT().
		IsFileOrDirExist(filePath).
		Return(false).
		Times(1)

	exists := service.IsExist(filePath)
	require.False(t, exists)
}

func TestGetAsset_Success(t *testing.T) {
	tests.SetupLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockClaimRepositoryInterface(ctrl)
	service := NewClaimService(mockRepo)

	username := "testuser"
	expectedAsset := &domain.Asset{ID: "assetID", Insured: username}

	mockRepo.EXPECT().
		GetAsset(username, "").
		Return(expectedAsset, nil).
		Times(1)

	asset, err := service.GetAsset(username, "")
	require.NoError(t, err)
	require.Equal(t, expectedAsset, asset)
}

func TestGetAsset_Error(t *testing.T) {
	tests.SetupLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockClaimRepositoryInterface(ctrl)
	service := NewClaimService(mockRepo)

	username := "testuser"

	mockRepo.EXPECT().
		GetAsset(username, "").
		Return(nil, fmt.Errorf("mock error fetching asset")).
		Times(1)

	asset, err := service.GetAsset(username, "")
	require.Error(t, err)
	require.Nil(t, asset)
	require.EqualError(t, err, "mock error fetching asset")
}

func TestUpdateAsset_Success(t *testing.T) {
	tests.SetupLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockClaimRepositoryInterface(ctrl)
	service := NewClaimService(mockRepo)

	asset := &domain.Asset{ID: "assetID", Insured: "testuser", CoverageValue: 0, CoverageDuration: 0, CoverageType: 0, Premium: 0, Partner: "Partner"}
	body := &domain.InvokeRequest{
		ChannelID:   constants.ChannelID,
		ChaincodeID: constants.ChaincodeID,
		Function:    "UpdateAsset",
		Args:        []string{asset.ID, asset.Insured, fmt.Sprintf("%d", asset.CoverageDuration), fmt.Sprintf("%d", asset.CoverageValue), fmt.Sprintf("%d", asset.CoverageType), asset.Partner, fmt.Sprintf("%d", asset.Premium), "Pending"},
	}

	mockRepo.EXPECT().
		UpdateAsset(body, "").
		Return(nil).
		Times(1)

	err := service.UpdateAssetClaimStatus(asset, "Pending", "")
	require.NoError(t, err)
}

func TestUpdateAsset_Error(t *testing.T) {
	tests.SetupLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockClaimRepositoryInterface(ctrl)
	service := NewClaimService(mockRepo)

	asset := &domain.Asset{ID: "assetID", Insured: "testuser", CoverageValue: 0, CoverageDuration: 0, CoverageType: 0, Premium: 0, Partner: "Partner"}
	body := &domain.InvokeRequest{
		ChannelID:   constants.ChannelID,
		ChaincodeID: constants.ChaincodeID,
		Function:    "UpdateAsset",
		Args:        []string{asset.ID, asset.Insured, fmt.Sprintf("%d", asset.CoverageDuration), fmt.Sprintf("%d", asset.CoverageValue), fmt.Sprintf("%d", asset.CoverageType), asset.Partner, fmt.Sprintf("%d", asset.Premium), "Pending"},
	}

	mockRepo.EXPECT().
		UpdateAsset(body, "").
		Return(fmt.Errorf("mock update asset error")).
		Times(1)

	err := service.UpdateAssetClaimStatus(asset, "Pending", "")
	require.Error(t, err)
	require.EqualError(t, err, "mock update asset error")
}
