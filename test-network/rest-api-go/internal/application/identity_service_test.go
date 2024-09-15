package application

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain/mocks"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/tests"
)

func TestIdentityService_Create_Success(t *testing.T) {
	tests.SetupLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommandExecutor := mocks.NewMockCommandExecutorInterface(ctrl)
	service := NewIdentityService(mockCommandExecutor)

	mockCommandExecutor.EXPECT().ExecuteCommand("/bin/bash", "./registerEnrollIdentity.sh").Return([]byte("success"), nil)

	credentials, err := service.Create()
	require.NoError(t, err)
	credentialsJSONAsBytes, err := json.Marshal(credentials)
	require.NoError(t, err)

	require.Contains(t, string(credentialsJSONAsBytes), "Username")
	require.Contains(t, string(credentialsJSONAsBytes), "Password")
}

func TestIdentityService_Create_Failure_ExecuteCommand(t *testing.T) {
	tests.SetupLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCommandExecutor := mocks.NewMockCommandExecutorInterface(ctrl)
	service := NewIdentityService(mockCommandExecutor)

	mockCommandExecutor.EXPECT().ExecuteCommand("/bin/bash", "./registerEnrollIdentity.sh").Return(nil, errors.New("command execution failed"))

	credentials, err := service.Create()
	require.Error(t, err)
	require.Nil(t, credentials)
}
