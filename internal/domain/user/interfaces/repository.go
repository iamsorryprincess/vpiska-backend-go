package interfaces

import "github.com/iamsorryprincess/vpiska-backend-go/internal/domain/user/models"

type Repository interface {
	Insert(user *models.User)
	GetByID(id string) (*models.User, error)
	GetByPhone(phone string) (*models.User, error)
	ChangePassword(id string, password string) error
	Update(id string, name string, phone string, imageID string) error
	CheckNameAndPhone(name string, phone string) error
}
