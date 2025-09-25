package repository

import (
	"context"

	"github.com/untea/bottom_babruysk/internal/domain"
)

type Users interface {
	CreateUser(ctx context.Context, request domain.CreateUserRequest) (*domain.CreateUserResponse, error)
	GetUser(ctx context.Context, request domain.GetUserRequest) (*domain.GetUserResponse, error)
	GetListUser(ctx context.Context, request domain.GetListUserRequest) (*domain.GetListUserResponse, error)
	UpdateUser(ctx context.Context, request domain.UpdateUserRequest) error
	DeleteUser(ctx context.Context, id domain.DeleteUserRequest) error
}
