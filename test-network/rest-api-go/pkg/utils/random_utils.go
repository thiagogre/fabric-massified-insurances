package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateRandomBytes generates random bytes of a given length
func GenerateRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	return bytes, nil
}

// GenerateRandomString generates a random string of a given length
func GenerateRandomString(length int) (string, error) {
	bytes, err := GenerateRandomBytes(length)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

// GenerateRandomCredentials generates a random username and password
func GenerateRandomCredentials(usernameLength, passwordLength int) (string, string, error) {
	username, err := GenerateRandomString(usernameLength)
	if err != nil {
		return "", "", err
	}

	password, err := GenerateRandomString(passwordLength)
	if err != nil {
		return "", "", err
	}

	return username, password, nil
}
