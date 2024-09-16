package application

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
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
		SaveFile(fileHeader, "test.pdf").
		Return(nil).
		Times(1)

	err := service.StoreClaim(fileHeader)
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
		SaveFile(fileHeader, "test.pdf").
		Return(fmt.Errorf("mock save file error")).
		Times(1)

	err := service.StoreClaim(fileHeader)
	require.Error(t, err)
	require.EqualError(t, err, "mock save file error")
}
