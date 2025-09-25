package service

import (
	"context"

	"github.com/untea/bottom_babruysk/internal/domain"
	"github.com/untea/bottom_babruysk/internal/repository"
)

type Users interface {
	CreateUser(ctx context.Context, request domain.CreateUserRequest) (*domain.CreateUserResponse, error)
	GetUser(ctx context.Context, request domain.GetUserRequest) (*domain.GetUserResponse, error)
	GetListUser(ctx context.Context, request domain.GetListUserRequest) (*domain.GetListUserResponse, error)
	UpdateUser(ctx context.Context, request domain.UpdateUserRequest) error
	DeleteUser(ctx context.Context, request domain.DeleteUserRequest) error
}

type usersService struct {
	repo repository.Users
}

func NewUsersSrv(repo repository.Users) Users {
	return &usersService{repo: repo}
}

func (s *usersService) CreateUser(ctx context.Context, request domain.CreateUserRequest) (*domain.CreateUserResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	return s.repo.CreateUser(ctx, request)
}

func (s *usersService) GetUser(ctx context.Context, request domain.GetUserRequest) (*domain.GetUserResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	return s.repo.GetUser(ctx, request)
}

func (s *usersService) GetListUser(ctx context.Context, request domain.GetListUserRequest) (*domain.GetListUserResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	return s.repo.GetListUser(ctx, request)
}

func (s *usersService) UpdateUser(ctx context.Context, request domain.UpdateUserRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}

	return s.repo.UpdateUser(ctx, request)
}

func (s *usersService) DeleteUser(ctx context.Context, request domain.DeleteUserRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}

	return s.repo.DeleteUser(ctx, request)
}
