package auth

import "errors"

var ErrInvalidToken = errors.New("unauthorized")

type TokenData struct {
	ID      string
	Name    string
	ImageID string
}

type TokenManager interface {
	GetAccessToken(input TokenData) (string, error)
	ParseToken(token string) (TokenData, error)
}
