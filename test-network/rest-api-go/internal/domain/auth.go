package domain

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

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
