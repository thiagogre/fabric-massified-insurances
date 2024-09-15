package domain

type User struct {
	Id    string
	Token string
}

type AuthServiceInterface interface {
	AuthenticateUser(username, password string) (*User, error)
}

type UserRepositoryInterface interface {
	GetUserById(id string) (*User, error)
}
