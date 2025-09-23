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
	ID           *uuid.UUID `db:"id"`
	Email        *string    `db:"email"`
	PasswordHash *string    `db:"password_hash"`
	DisplayName  *string    `db:"display_name"`
	Role         *UserRole  `db:"role"`
	CreatedAt    *time.Time `db:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at"`
}

type CreateUserRequest struct {
	Email        *string   `db:"email" json:"email,omitempty"`
	PasswordHash *string   `db:"password_hash" json:"password_hash,omitempty"`
	DisplayName  *string   `db:"display_name" json:"display_name,omitempty"`
	Role         *UserRole `db:"role" json:"role,omitempty"`
}

type UpdateUserRequest struct {
	DisplayName *string   `db:"display_name" json:"display_name,omitempty"`
	Role        *UserRole `db:"role" json:"role,omitempty"`
}
