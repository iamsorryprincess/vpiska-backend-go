package auth

type jwtManager struct {
}

func NewJwtManager() TokenManager {
	return &jwtManager{}
}

func (m *jwtManager) GetAccessToken(input CreateTokenInput) string {
	return "token"
}
