package domain

import (
	"time"

	"github.com/google/uuid"
)

type AlbumTrack struct {
	AlbumID    *uuid.UUID `db:"album_id"`
	DiscNumber *int       `db:"disc_number"`
	Position   *int       `db:"position"`
	TrackID    *uuid.UUID `db:"track_id"`
	CreatedAt  *time.Time `db:"created_at"`
}
