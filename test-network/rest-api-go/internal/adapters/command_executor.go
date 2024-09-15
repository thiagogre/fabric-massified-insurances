package adapters

import (
	"os/exec"

	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/constants"
)

type CommandExecutor struct{}

func NewCommandExecutor() *CommandExecutor {
	return &CommandExecutor{}
}

func (s *CommandExecutor) ExecuteCommand(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = constants.TestNetworkPath
	return cmd.CombinedOutput()
}
