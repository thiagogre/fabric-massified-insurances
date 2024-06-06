package repositories

import "rest-api-go/internal/models"

type UserRepository interface {
	GetUserById(id string) (*models.User, error)
}
