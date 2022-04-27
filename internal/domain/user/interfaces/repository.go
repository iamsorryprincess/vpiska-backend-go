package interfaces

import (
	"context"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain/user/models"
)

type Repository interface {
	Insert(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByPhone(ctx context.Context, phone string) (*models.User, error)
	ChangePassword(ctx context.Context, id string, password string) error
	Update(ctx context.Context, id string, name string, phone string, imageID string) error
	CheckNameAndPhone(ctx context.Context, name string, phone string) error
}
