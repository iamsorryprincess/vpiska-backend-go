package interfaces

import "github.com/iamsorryprincess/vpiska-backend-go/internal/domain/user/models"

type IdentityProvider interface {
	GetAccessToken(user *models.User) string
}
