package domain

import (
	"time"

	"github.com/google/uuid"
)

type Albums struct {
	ID          *uuid.UUID `db:"id"`
	OwnerID     *uuid.UUID `db:"owner_id"`
	Title       *string    `db:"title"`
	Description *string    `db:"description"`
	ReleaseDate *time.Time `db:"release_date"`
	CreatedAt   *time.Time `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}
