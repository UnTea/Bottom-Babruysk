package postgres

import (
	"context"

	"github.com/untea/bottom_babruysk/internal/domain"
	"github.com/untea/bottom_babruysk/internal/repository"
)

type UsersRepository struct {
	db *repository.Client // TODO реализовать интерфейс для fetch и прокидывать просто r.db
}

func NewUsersRepository(db *repository.Client) *UsersRepository {
	return &UsersRepository{db: db}
}

func (r *UsersRepository) CreateUser(ctx context.Context, request domain.CreateUserRequest) (*domain.CreateUserResponse, error) {
	const createUserSQL = `
		insert into users (email, password_hash, display_name, role)
		values ($1, $2, $3, coalesce($4::user_role, 'user'::user_role))
		returning id;
	`

	arguments := []any{
		request.Email,
		request.PasswordHash,
		request.DisplayName,
		request.Role,
	}

	user, err := repository.FetchOne[domain.User](ctx, r.db.Driver(), createUserSQL, arguments...) // TODO реализовать интерфейс для fetch и прокидывать просто r.db
	if err != nil {
		return nil, err
	}

	return &domain.CreateUserResponse{ID: user.ID}, nil
}

func (r *UsersRepository) GetUser(ctx context.Context, request domain.GetUserRequest) (*domain.GetUserResponse, error) {
	const getUserSQL = `
		select id, email, password_hash, display_name, role, created_at 
		from users 
		where id = $1;
	`

	arguments := []any{
		request.ID,
	}

	user, err := repository.FetchOne[domain.User](ctx, r.db.Driver(), getUserSQL, arguments...) // TODO реализовать интерфейс для fetch и прокидывать просто r.db
	if err != nil {
		return nil, err
	}

	return &domain.GetUserResponse{User: user}, nil
}

func (r *UsersRepository) ListUsers(ctx context.Context, request domain.ListUsersRequest) (*domain.ListUsersResponse, error) {
	const getListUsersSQL = `
		with params as (
			select
				$1::user_role                                 as role_filter,
				coalesce(nullif(lower($2), ''), 'created_at') as sort_field,
				coalesce(nullif(lower($3), ''), 'desc')       as sort_order,
				greatest(coalesce($4, 50), 1)                 as limit_val,
				greatest(coalesce($5, 1), 1)                  as offset_val
		)
		select
			u.id, u.email, u.password_hash, u.display_name, u.role, u.created_at, u.updated_at
		from users u, params p
		where
			(p.role_filter is null or u.role = p.role_filter)
		order by
			-- email
			case when p.sort_field = 'email'        and p.sort_order = 'asc'  then u.email        end asc  nulls last,
			case when p.sort_field = 'email'        and p.sort_order = 'desc' then u.email        end desc nulls last,
		
			-- display_name
			case when p.sort_field = 'display_name' and p.sort_order = 'asc'  then u.display_name end asc  nulls last,
			case when p.sort_field = 'display_name' and p.sort_order = 'desc' then u.display_name end desc nulls last,
		
			-- role
			case when p.sort_field = 'role'         and p.sort_order = 'asc'  then u.role::text   end asc  nulls last,
			case when p.sort_field = 'role'         and p.sort_order = 'desc' then u.role::text   end desc nulls last,
		
			-- created_at
			case when p.sort_field = 'created_at'   and p.sort_order = 'asc'  then u.created_at   end asc  nulls last,
			case when p.sort_field = 'created_at'   and p.sort_order = 'desc' then u.created_at   end desc nulls last,
		
			-- updated_at
			case when p.sort_field = 'updated_at'   and p.sort_order = 'asc'  then u.updated_at   end asc  nulls last,
			case when p.sort_field = 'updated_at'   and p.sort_order = 'desc' then u.updated_at   end desc nulls last,

			u.created_at desc
		limit (select limit_val from params)
		offset (select offset_val from params);
	`

	arguments := []any{
		request.Role,
		request.SortField,
		request.SortOrder,
		request.Limit,
		request.Offset,
	}

	users, err := repository.FetchMany[domain.User](ctx, r.db.Driver(), getListUsersSQL, arguments...) // TODO реализовать интерфейс для fetch и прокидывать просто r.db
	if err != nil {
		return nil, err
	}

	return &domain.ListUsersResponse{Users: users}, nil
}

func (r *UsersRepository) UpdateUser(ctx context.Context, request domain.UpdateUserRequest) error {
	const updateUserSQL = `
		update users
		set
			display_name = coalesce($2, display_name),
			role         = coalesce($3::user_role, role),
			updated_at   = now()
		where id = $1;
	`

	arguments := []any{
		request.ID,
		request.DisplayName,
		request.Role,
	}

	_, err := repository.FetchOne[domain.User](ctx, r.db.Driver(), updateUserSQL, arguments...) // TODO реализовать интерфейс для fetch и прокидывать просто r.db
	if err != nil {
		return err
	}

	return nil
}

func (r *UsersRepository) DeleteUser(ctx context.Context, request domain.DeleteUserRequest) error {
	const deleteUserSQL = `
		delete from users where id = $1;
	`

	affected, err := repository.ExecAffected(ctx, r.db.Driver(), deleteUserSQL, request.ID) // TODO реализовать интерфейс для fetch и прокидывать просто r.db
	if err != nil {
		return err
	}

	if affected == 0 {
		return repository.ErrNotFound
	}

	return nil
}
