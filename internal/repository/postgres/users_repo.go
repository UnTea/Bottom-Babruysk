package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"

	"bottom_babruysk/internal/domain"
	"bottom_babruysk/internal/repository"
)

type UsersRepository struct {
	DB *repository.Client
}

func NewUsersRepo(db *repository.Client) *UsersRepository {
	return &UsersRepository{DB: db}
}

func (r *UsersRepository) Create(ctx context.Context, request domain.CreateUserRequest) (*domain.User, error) {
	const createUserSQL = `
		insert into users (id, email, password_hash, display_name, role, created_at)
		values ($1, $2, $3, $4, coalesce($5::user_role, 'user'::user_role), $6)
		returning id, email, password_hash, display_name, role, created_at;
	`

	arguments := []any{
		uuid.New(),
		request.Email,
		request.PasswordHash,
		request.DisplayName,
		request.Role,
		time.Now().UTC(),
	}

	return repository.FetchOne[domain.User](ctx, r.DB.Driver(), createUserSQL, arguments...)
}

func (r *UsersRepository) Get(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	const getUserSQL = `
		select id, email, password_hash, display_name, role, created_at 
		from users 
		where id = $1
		limit 1;
	`

	return repository.FetchOne[domain.User](ctx, r.DB.Driver(), getUserSQL, id)
}

func (r *UsersRepository) List(ctx context.Context, page domain.Page, role, search *string) ([]domain.User, error) {
	const getUserListSQL = `
		select id, email, password_hash, display_name, role, created_at, updated_at
		from users
		where
		  ($1::user_role is null or role = $1::user_role)
		  and ($2 is null or (email ilike '%' || $2 || '%' or display_name ilike '%' || $2 || '%'))
		order by created_at desc
		limit $3 offset $4;
	`

	if page.Limit <= 0 {
		page.Limit = 50
	}

	if page.Offset < 0 {
		page.Offset = 0
	}

	var roleArgument, searchArgument any

	if role != nil && *role != "" {
		roleArgument = *role
	} else {
		roleArgument = nil
	}

	if search != nil && *search != "" {
		searchArgument = *search
	} else {
		searchArgument = nil
	}

	arguments := []any{
		roleArgument,
		searchArgument,
		page.Limit,
		page.Offset,
	}

	users, err := repository.FetchMany[domain.User](ctx, r.DB.Driver(), getUserListSQL, arguments...)
	if err != nil {
		return nil, err
	}

	result := make([]domain.User, 0, len(users))

	for _, user := range users {
		if user != nil {
			result = append(result, *user)
		}
	}

	return result, nil
}

func (r *UsersRepository) Update(ctx context.Context, id uuid.UUID, upd domain.UpdateUserRequest) (*domain.User, error) {
	const updateUserSQL = `
		update users
		set
			display_name = coalesce($2, display_name),
			role         = coalesce($3::user_role, role),
			updated_at   = now()
		where id = $1
		returning id, email, password_hash, display_name, role, created_at, updated_at;
	`

	arguments := []any{
		id,
		upd.DisplayName,
		upd.Role,
	}

	return repository.FetchOne[domain.User](ctx, r.DB.Driver(), updateUserSQL, arguments...)
}

func (r *UsersRepository) Delete(ctx context.Context, id uuid.UUID) error {
	const deleteUserSQL = `
		delete from users where id = $1;
	`

	affected, err := repository.ExecAffected(ctx, r.DB.Driver(), deleteUserSQL, id)
	if err != nil {
		return err
	}

	if affected == 0 {
		return repository.ErrNotFound
	}

	return nil
}
