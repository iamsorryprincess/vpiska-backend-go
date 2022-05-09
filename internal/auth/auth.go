package auth

type CreateTokenInput struct {
	ID      string
	Name    string
	ImageID string
}

type TokenManager interface {
	GetAccessToken(input CreateTokenInput) string
}
