package cmd

import (
	"os/exec"

	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/constants"
)

type CommandExecutorInterface interface {
	ExecuteCommand(name string, args ...string) ([]byte, error)
}

type CommandExecutor struct{}

func (c *CommandExecutor) ExecuteCommand(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = constants.TestNetworkPath
	return cmd.CombinedOutput()
}

var _ CommandExecutorInterface = (*CommandExecutor)(nil)
