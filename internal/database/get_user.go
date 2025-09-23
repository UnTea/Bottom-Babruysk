package database

import (
	"context"
	"time"

	"bottom_babruysk/internal/repository"
)

type (
	User struct {
		Id          *string    `db:"id"`
		Email       *string    `db:"email"`
		DisplayName *string    `db:"display_name"`
		CreatedAt   *time.Time `db:"created_at"`
		UpdatedAt   *time.Time `db:"updated_at"`
	}

	GetUserRequest struct {
		UserID string
	}

	GetUserResponse struct {
		Users *User
	}
)

func (r *GetUserRequest) GetUser(ctx context.Context, client *repository.Client) (*GetUserResponse, error) {
	const getUserSQL = `
	select id, 
	       email, 
	       display_name, 
	       created_at, 
	       updated_at 
	from users 
	where id=$1 
	limit 1;
	`

	user, err := repository.FetchOne[User](ctx, client.Driver(), getUserSQL, r.UserID)
	if err != nil {
		return nil, err
	}

	return &GetUserResponse{
		Users: user,
	}, nil
}
