package identity

import "github.com/iamsorryprincess/vpiska-backend-go/internal/domain/user/interfaces"

type PasswordHashProvider struct {
}

func InitPasswordHashProvider() interfaces.PasswordHashProvider {
	return &PasswordHashProvider{}
}

func (p *PasswordHashProvider) HashPassword(password string) string {
	return password
}

func (p *PasswordHashProvider) VerifyHashPassword(hashedPassword string, providedPassword string) bool {
	return true
}
