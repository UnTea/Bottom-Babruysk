package repository

import (
	"context"

	"github.com/untea/bottom_babruysk/internal/domain"
)

type Users interface {
	CreateUser(context.Context, domain.CreateUserRequest) (*domain.CreateUserResponse, error)
	GetUser(context.Context, domain.GetUserRequest) (*domain.GetUserResponse, error)
	ListUsers(context.Context, domain.GetListUserRequest) (*domain.GetListUserResponse, error)
	UpdateUser(context.Context, domain.UpdateUserRequest) error
	DeleteUser(context.Context, domain.DeleteUserRequest) error
}
