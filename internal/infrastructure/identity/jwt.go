package identity

import (
	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain/user/interfaces"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain/user/models"
)

type JwtProvider struct {
}

func InitJwtIdentityProvider() interfaces.IdentityProvider {
	return &JwtProvider{}
}

func (jwtProvider *JwtProvider) GetAccessToken(user *models.User) string {
	return "token"
}
