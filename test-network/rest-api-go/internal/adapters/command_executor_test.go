package adapters

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/tests"
)

func TestExecuteCommand_Success(t *testing.T) {
	tests.SetupLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	executor := NewCommandExecutor()
	output, err := executor.ExecuteCommand("echo", "hello")

	require.NoError(t, err)
	require.Equal(t, "hello\n", string(output))
}
