package domain

type IdentityResponse struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Credentials struct {
	Username string
	Password string
}

type IdentityInterface interface {
	Create() (*Credentials, error)
}
