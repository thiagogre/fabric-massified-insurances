package services

import (
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/models"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/repositories"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/utils"
)

type AuthService struct {
	UserRepository repositories.UserRepository
}

func (s *AuthService) AuthenticateUser(username, password string) (*models.User, error) {
	user, err := s.UserRepository.GetUserById(username)
	if err != nil {
		return nil, err
	}

	if !utils.CheckPasswordHash(password, user.Token) {
		return nil, nil
	}

	return user, nil
}
