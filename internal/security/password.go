package security

type passwordManager struct {
}

func NewPasswordManager() PasswordManager {
	return &passwordManager{}
}

func (m *passwordManager) HashPassword(password string) string {
	return password
}

func (m *passwordManager) VerifyHashedPassword(hashedPassword string, providedPassword string) bool {
	return hashedPassword == providedPassword
}
