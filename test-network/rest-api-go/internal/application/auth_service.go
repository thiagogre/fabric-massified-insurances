package application

import (
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/utils"
)

type AuthService struct {
	UserRepository domain.UserRepositoryInterface
}

func NewAuthService(repo domain.UserRepositoryInterface) *AuthService {
	return &AuthService{UserRepository: repo}
}

func (s *AuthService) AuthenticateUser(username, password string) (*domain.User, error) {
	user, err := s.UserRepository.GetUserById(username)
	if err != nil {
		return nil, err
	}

	if !utils.CheckPasswordHash(password, user.Token) {
		return nil, nil
	}

	return user, nil
}
