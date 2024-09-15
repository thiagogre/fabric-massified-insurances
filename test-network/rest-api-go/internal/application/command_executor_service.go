package application

import (
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain"
)

type CommandExecutorService struct {
	CommandExecutor domain.CommandExecutorInterface
}

func NewCommandExecutorService(commandExecutor domain.CommandExecutorInterface) *CommandExecutorService {
	return &CommandExecutorService{CommandExecutor: commandExecutor}
}

func (s *CommandExecutorService) Execute(name string, args ...string) ([]byte, error) {
	return s.CommandExecutor.ExecuteCommand(name, args...)
}
