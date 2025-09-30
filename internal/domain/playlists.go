package domain

import (
	"time"

	"github.com/google/uuid"
)

type Playlist struct {
	ID          *uuid.UUID  `db:"id"`
	OwnerID     *uuid.UUID  `db:"owner_id"`
	Title       *string     `db:"title"`
	Description *string     `db:"description"`
	Visibility  *Visibility `db:"visibility"`
	CreatedAt   *time.Time  `db:"created_at"`
	UpdatedAt   *time.Time  `db:"updated_at"`
}

type (
	CreatePlaylistRequest struct {
		OwnerID     *uuid.UUID  `db:"owner_id"    json:"owner_id,omitempty"`
		Title       *string     `db:"title"       json:"title,omitempty"`
		Description *string     `db:"description" json:"description,omitempty"`
		Visibility  *Visibility `db:"visibility"  json:"visibility,omitempty"`
	}

	CreatePlaylistResponse struct {
		ID *uuid.UUID `json:"id,omitempty"`
	}
)

type (
	GetPlaylistRequest struct {
		ID *uuid.UUID `db:"id" json:"id,omitempty" path:"id"`
	}

	GetPlaylistResponse struct {
		Playlist *Playlist `json:"playlist"`
	}
)

type (
	ListPlaylistsRequest struct {
		Limit      *int        `db:"limit"       query:"limit"`
		Offset     *int        `db:"offset"      query:"offset"`
		OwnerID    *uuid.UUID  `db:"owner_id"    query:"owner_id"`
		Visibility *Visibility `db:"visibility"  query:"visibility"`
		SortField  *string     `db:"sort_field"  query:"sort_field"`
		SortOrder  *string     `db:"sort_order"  query:"sort_order"`
	}

	ListPlaylistsResponse struct {
		Playlists []*Playlist `json:"playlists,omitempty"`
	}
)

type UpdatePlaylistRequest struct {
	ID          *uuid.UUID  `db:"id"          json:"-" path:"id"`
	Title       *string     `db:"title"       json:"title,omitempty"`
	Description *string     `db:"description" json:"description,omitempty"`
	Visibility  *Visibility `db:"visibility"  json:"visibility,omitempty"`
}

type DeletePlaylistRequest struct {
	ID *uuid.UUID `db:"id" json:"-" path:"id"`
}
