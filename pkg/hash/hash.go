package hash

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type PasswordHashManager interface {
	HashPassword(password string) (string, error)
	VerifyPassword(password string, hashedPassword string) (bool, error)
}

type passwordHashManager struct {
}

func NewPasswordHashManager() PasswordHashManager {
	return &passwordHashManager{}
}

func (m *passwordHashManager) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)

	if err != nil {
		return "", err
	}

	return string(hash), err
}

func (m *passwordHashManager) VerifyPassword(password string, hashedPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err != nil {
		if errors.Is(bcrypt.ErrMismatchedHashAndPassword, err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
