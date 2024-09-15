package adapters

import (
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/internal/domain"
	"github.com/thiagogre/fabric-massified-insurances/test-network/rest-api-go/pkg/db"
)

type AuthRepository struct {
	DB db.Database
}

func NewAuthRepository(db db.Database) *AuthRepository {
	return &AuthRepository{DB: db}
}

func (repo *AuthRepository) GetUserById(id string) (*domain.User, error) {
	query := "SELECT id, token FROM users WHERE id = ?"
	row := repo.DB.QueryRow(query, id)

	var user domain.User
	err := row.Scan(&user.Id, &user.Token)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
