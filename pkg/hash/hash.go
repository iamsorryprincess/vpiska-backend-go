package hash

type PasswordHashManager interface {
	HashPassword(password string) string
	VerifyHashedPassword(hashedPassword string, providedPassword string) bool
}

type passwordHashManager struct {
}

func NewPasswordHashManager() PasswordHashManager {
	return &passwordHashManager{}
}

func (m *passwordHashManager) HashPassword(password string) string {
	return password
}

func (m *passwordHashManager) VerifyHashedPassword(hashedPassword string, providedPassword string) bool {
	return hashedPassword == providedPassword
}
