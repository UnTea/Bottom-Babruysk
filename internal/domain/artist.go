package domain

import (
	"time"

	"github.com/google/uuid"
)

type Artists struct {
	ID        *uuid.UUID `db:"id"`
	Name      *string    `db:"name"`
	Bio       *string    `db:"bio"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}
