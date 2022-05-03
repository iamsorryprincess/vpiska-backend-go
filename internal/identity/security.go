package identity

import "github.com/iamsorryprincess/vpiska-backend-go/internal/domain/user"

type PasswordHashProvider struct {
}

func NewPasswordHashProvider() user.SecurityProvider {
	return &PasswordHashProvider{}
}

func (p *PasswordHashProvider) HashPassword(password string) string {
	return password
}

func (p *PasswordHashProvider) VerifyHashedPassword(hashedPassword string, providedPassword string) bool {
	return hashedPassword == providedPassword
}
