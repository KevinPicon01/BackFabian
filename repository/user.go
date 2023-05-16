package repository

import (
	"context"
	"kevinPicon/go/rest-ws/models"
)

type UserRepository interface {
	InsertUser(ctx context.Context, user *models.User) error
	GetUserById(ctx context.Context, id string) (*models.UserPayload, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUsers(ctx context.Context) ([]*models.UserPayload, error)
	CreateUserService(ctx context.Context, service models.Service) error
	UpdateEcan(ctx context.Context, id string) error
	Close() error
}

var impl UserRepository

func SetRepository(repo UserRepository) {
	impl = repo
}

func InsertUser(ctx context.Context, user *models.User) error {
	return impl.InsertUser(ctx, user)
}

func GetUserById(ctx context.Context, id string) (*models.UserPayload, error) {
	return impl.GetUserById(ctx, id)
}

func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return impl.GetUserByEmail(ctx, email)
}

func GetUsers(ctx context.Context) ([]*models.UserPayload, error) {
	return impl.GetUsers(ctx)
}
func CreateUserService(ctx context.Context, service models.Service) error {
	return impl.CreateUserService(ctx, service)
}
func UpdateEcan(ctx context.Context, id string) error {
	return impl.UpdateEcan(ctx, id)
}
func Close() error {
	return impl.Close()
}
