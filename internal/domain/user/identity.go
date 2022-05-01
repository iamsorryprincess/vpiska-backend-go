package user

type IdentityProvider interface {
	GetAccessToken(user *User) string
}
