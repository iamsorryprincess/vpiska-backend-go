package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

type PasswordHashManager interface {
	HashPassword(password string) (string, error)
	VerifyPassword(password string, hashedPassword string) (bool, error)
}

type passwordHashManager struct {
	key []byte
}

func NewPasswordHashManager(key string) (PasswordHashManager, error) {
	bytes, err := hex.DecodeString(key)

	if err != nil {
		return nil, err
	}

	return &passwordHashManager{
		key: bytes,
	}, nil
}

func (m *passwordHashManager) HashPassword(password string) (string, error) {
	h := hmac.New(sha256.New, m.key)
	_, err := h.Write([]byte(password))

	if err != nil {
		return "", err
	}

	result := h.Sum(nil)
	return hex.EncodeToString(result), nil
}

func (m *passwordHashManager) VerifyPassword(password string, hashedPassword string) (bool, error) {
	hashedPasswordBytes, err := hex.DecodeString(hashedPassword)

	if err != nil {
		return false, err
	}

	h := hmac.New(sha256.New, m.key)
	_, err = h.Write([]byte(password))

	if err != nil {
		return false, err
	}

	result := h.Sum(nil)
	return hmac.Equal(result, hashedPasswordBytes), nil
}
