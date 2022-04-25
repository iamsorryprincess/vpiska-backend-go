package interfaces

type PasswordHashProvider interface {
	HashPassword(password string) string
	VerifyHashPassword(hashedPassword string, providedPassword string) bool
}
