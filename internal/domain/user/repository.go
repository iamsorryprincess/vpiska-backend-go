package user

import (
	"context"
)

type Repository interface {
	CheckNameAndPhone(ctx context.Context, name string, phone string) error
	CreateUser(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, id string) (*User, error)
	GetUserByPhone(ctx context.Context, phone string) (*User, error)
	ChangePassword(ctx context.Context, id string, password string) error
}
