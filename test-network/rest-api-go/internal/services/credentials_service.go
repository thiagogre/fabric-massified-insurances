package services

import (
	"crypto/rand"
	"encoding/base64"
)

// Credentials holds the username and password
type Credentials struct {
	Username string
	Password string
}

// NewCredentials creates a new Credentials object
func NewCredentials() *Credentials {
	return &Credentials{}
}

// generateRandomBytes generates random bytes of a given length
func (c *Credentials) generateRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	return bytes, nil
}

// generateRandomString generates a random string of a given length
func (c *Credentials) generateRandomString(length int) (string, error) {
	bytes, err := c.generateRandomBytes(length)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

// Generate generates random credentials (username and password)
func (c *Credentials) Generate(usernameLength, passwordLength int) error {
	username, err := c.generateRandomString(usernameLength)
	if err != nil {
		return err
	}
	c.Username = username

	password, err := c.generateRandomString(passwordLength)
	if err != nil {
		return err
	}
	c.Password = password

	return nil
}
