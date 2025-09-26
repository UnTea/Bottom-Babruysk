package service

import (
	"context"

	"github.com/untea/bottom_babruysk/internal/domain"
	"github.com/untea/bottom_babruysk/internal/repository"
)

type UsersService struct {
	repository repository.Users
}

func NewUsersService(repository repository.Users) *UsersService {
	return &UsersService{repository: repository}
}

func (s *UsersService) CreateUser(ctx context.Context, request domain.CreateUserRequest) (*domain.CreateUserResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, err
	}

	return s.repository.CreateUser(ctx, request)
}

func (s *UsersService) GetUser(ctx context.Context, request domain.GetUserRequest) (*domain.GetUserResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, err
	}

	return s.repository.GetUser(ctx, request)
}

func (s *UsersService) ListUsers(ctx context.Context, request domain.GetListUserRequest) (*domain.GetListUserResponse, error) {
	err := request.Validate()
	if err != nil {
		return nil, err
	}

	return s.repository.ListUsers(ctx, request)
}

func (s *UsersService) UpdateUser(ctx context.Context, request domain.UpdateUserRequest) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	return s.repository.UpdateUser(ctx, request)
}

func (s *UsersService) DeleteUser(ctx context.Context, request domain.DeleteUserRequest) error {
	err := request.Validate()
	if err != nil {
		return err
	}

	return s.repository.DeleteUser(ctx, request)
}
