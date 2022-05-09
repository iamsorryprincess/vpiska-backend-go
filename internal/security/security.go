package security

type PasswordManager interface {
	HashPassword(password string) string
	VerifyHashedPassword(hashedPassword string, providedPassword string) bool
}
