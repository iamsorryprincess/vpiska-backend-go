package user

type SecurityProvider interface {
	HashPassword(password string) string
	VerifyHashedPassword(hashedPassword string, providedPassword string) bool
}
