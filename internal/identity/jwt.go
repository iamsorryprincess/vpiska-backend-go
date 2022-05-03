package identity

import "github.com/iamsorryprincess/vpiska-backend-go/internal/domain/user"

type JwtTokenProvider struct {
}

func NewJwtTokenProvider() user.IdentityProvider {
	return &JwtTokenProvider{}
}

func (p *JwtTokenProvider) GetAccessToken(user *user.User) string {
	return "token"
}
