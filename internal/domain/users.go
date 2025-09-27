package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

type User struct {
	ID           *uuid.UUID `db:"id"            json:"id,omitempty"`
	Email        *string    `db:"email"         json:"email,omitempty"`
	PasswordHash *string    `db:"password_hash" json:"password_hash,omitempty"`
	DisplayName  *string    `db:"display_name"  json:"display_name,omitempty"`
	Role         *UserRole  `db:"role"          json:"role,omitempty"`
	CreatedAt    *time.Time `db:"created_at"    json:"created_at,omitempty"`
	UpdatedAt    *time.Time `db:"updated_at"    json:"updated_at,omitempty"`
}

type (
	CreateUserRequest struct {
		Email        *string   `db:"email"         json:"email,omitempty"`
		PasswordHash *string   `db:"password_hash" json:"password_hash,omitempty"`
		DisplayName  *string   `db:"display_name"  json:"display_name,omitempty"`
		Role         *UserRole `db:"role"          json:"role,omitempty"`
	}

	CreateUserResponse struct {
		ID *uuid.UUID `json:"id,omitempty"`
	}
)

type (
	GetUserRequest struct {
		ID uuid.UUID `db:"id" json:"id,omitempty" path:"id"`
	}

	GetUserResponse struct {
		User *User `json:"user,omitempty"`
	}
)

type (
	ListUsersRequest struct {
		Limit     *int      `db:"limit"        json:"limit,omitempty"      query:"limit"`
		Offset    *int      `db:"offset"       json:"offset,omitempty"     query:"offset"`
		Role      *UserRole `db:"role"         json:"role,omitempty"       query:"role"`
		SortField *string   `db:"sort_field"   json:"sort_field,omitempty" query:"sort_field"`
		SortOrder *string   `db:"sort_order"   json:"sort_order,omitempty" query:"sort_order"`
	}

	ListUsersResponse struct {
		Users []*User `json:"users,omitempty"`
	}
)

type (
	UpdateUserRequest struct {
		ID          *uuid.UUID `db:"id"           json:"-" path:"id"`
		DisplayName *string    `db:"display_name" json:"-" query:"display_name"`
		Role        *UserRole  `db:"role"         json:"-" query:"role"`
	}
)

type (
	DeleteUserRequest struct {
		ID uuid.UUID `db:"id" json:"-" path:"id"`
	}
)
