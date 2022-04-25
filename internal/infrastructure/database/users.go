package database

import (
	"sync"

	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain/user/errors"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain/user/interfaces"
	"github.com/iamsorryprincess/vpiska-backend-go/internal/domain/user/models"
)

type UserRepository struct {
	storage map[string]*models.User
	mutex   sync.Mutex
}

func InitUserRepository() interfaces.Repository {
	return &UserRepository{
		storage: make(map[string]*models.User),
		mutex:   sync.Mutex{},
	}
}

func (r *UserRepository) Insert(user *models.User) {
	r.mutex.Lock()
	r.storage[user.ID] = user
	r.mutex.Unlock()
}

func (r *UserRepository) GetByID(id string) (*models.User, error) {
	user := r.storage[id]

	if user == nil {
		return nil, errors.UserNotFound
	}

	return user, nil
}

func (r *UserRepository) GetByPhone(phone string) (*models.User, error) {
	return nil, nil
}

func (r *UserRepository) ChangePassword(id string, password string) error {
	return nil
}

func (r *UserRepository) Update(id string, name string, phone string, imageID string) error {
	return nil
}

func (r *UserRepository) CheckNameAndPhone(name string, phone string) error {
	return nil
}
