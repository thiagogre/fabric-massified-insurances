package repositories

import (
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/models"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/db"
)

type SQLUserRepository struct {
	DB db.Database
}

func (repo *SQLUserRepository) GetUserById(id string) (*models.User, error) {
	query := "SELECT id, token FROM users WHERE id = ?"
	row := repo.DB.QueryRow(query, id)

	var user models.User
	err := row.Scan(&user.Id, &user.Token)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
