package repository

import (
	"context"

	"github.com/google/uuid"

	"bottom_babruysk/internal/domain"
)

type Users interface {
	Create(ctx context.Context, cu domain.CreateUserRequest) (*domain.User, error)
	Get(ctx context.Context, id uuid.UUID) (*domain.User, error)
	List(ctx context.Context, page domain.Page, role, search *string) ([]domain.User, error)
	Update(ctx context.Context, id uuid.UUID, upd domain.UpdateUserRequest) (*domain.User, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
