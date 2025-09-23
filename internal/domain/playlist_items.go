package domain

import (
	"time"

	"github.com/google/uuid"
)

type PlaylistItems struct {
	PlaylistID *uuid.UUID `db:"playlist_id"`
	TrackID    *uuid.UUID `db:"track_id"`
	Position   *int       `db:"position"`
	AddedAt    *time.Time `db:"added_at"`
}
