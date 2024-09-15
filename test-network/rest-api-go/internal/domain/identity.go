package domain

type Credentials struct {
	Username string
	Password string
}

type IdentityInterface interface {
	Create() (*Credentials, error)
}
