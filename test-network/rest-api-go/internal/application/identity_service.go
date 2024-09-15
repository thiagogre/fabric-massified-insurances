package application

import (
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/constants"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/logger"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/utils"
)

type IdentityService struct {
	CommandExecutor domain.CommandExecutorInterface
}

func NewIdentityService(commandExecutor domain.CommandExecutorInterface) *IdentityService {
	return &IdentityService{CommandExecutor: commandExecutor}
}

func generateCredentials(usernameLength, passwordLength int) (*domain.Credentials, error) {
	credentials := domain.Credentials{}

	username, err := utils.GenerateRandomString(usernameLength)
	if err != nil {
		return nil, err
	}
	credentials.Username = username

	password, err := utils.GenerateRandomString(passwordLength)
	if err != nil {
		return nil, err
	}
	credentials.Password = password

	return &credentials, nil
}

func (s *IdentityService) Create() (*domain.Credentials, error) {
	credentials, err := generateCredentials(constants.DefaultUsernameLength, constants.DefaultPasswordLength)
	if err != nil {
		return nil, err
	}

	output, err := s.CommandExecutor.ExecuteCommand("/bin/bash", "./registerEnrollIdentity.sh", credentials.Username, credentials.Password)
	if err != nil {
		return nil, err
	}

	logger.Info("Script output: " + string(output))

	return credentials, nil
}
