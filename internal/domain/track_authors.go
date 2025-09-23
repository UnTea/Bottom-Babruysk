package domain

import (
	"time"

	"github.com/google/uuid"
)

type TrackAuthors struct {
	ID        *uuid.UUID `db:"id"`
	TrackID   *uuid.UUID `db:"track_id"`
	ArtistID  *uuid.UUID `db:"artist_id"`
	Ord       *int       `db:"ord"`
	CreatedAt *time.Time `db:"createdAt"`
}
