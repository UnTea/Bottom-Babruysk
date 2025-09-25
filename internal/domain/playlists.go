package domain

import (
	"time"

	"github.com/google/uuid"
)

type Playlist struct {
	ID          *uuid.UUID  `db:"id"`
	OwnerID     *string     `db:"owner_id"`
	Title       *string     `db:"title"`
	Description *string     `db:"description"`
	Visibility  *Visibility `db:"visibility"`
	CreatedAt   *time.Time  `db:"created_at"`
	UpdatedAt   *time.Time  `db:"updated_at"`
}
