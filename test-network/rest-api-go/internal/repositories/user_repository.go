package repositories

import "github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/models"

type UserRepository interface {
	GetUserById(id string) (*models.User, error)
}
