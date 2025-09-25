package domain

import (
	"time"

	"github.com/google/uuid"
)

type TrackLike struct {
	UserID    *uuid.UUID `db:"user_id"`
	TrackID   *uuid.UUID `db:"track_id"`
	CreatedAt time.Time  `db:"created_at"`
}
