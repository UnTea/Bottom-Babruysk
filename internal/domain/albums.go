package domain

import (
	"time"

	"github.com/google/uuid"
)

type Album struct {
	ID          *uuid.UUID `db:"id"           json:"id,omitempty"`
	OwnerID     *uuid.UUID `db:"owner_id"     json:"owner_id,omitempty"`
	Title       *string    `db:"title"        json:"title,omitempty"`
	Description *string    `db:"description"  json:"description,omitempty"`
	ReleaseDate *time.Time `db:"release_date" json:"release_date,omitempty"`
	CreatedAt   *time.Time `db:"created_at"   json:"created_at,omitempty"`
	UpdatedAt   *time.Time `db:"updated_at"   json:"updated_at,omitempty"`
}

type (
	CreateAlbumRequest struct {
		OwnerID     *uuid.UUID `db:"owner_id"     json:"owner_id,omitempty"`
		Title       *string    `db:"title"        json:"title,omitempty"`
		Description *string    `db:"description"  json:"description,omitempty"`
		ReleaseDate *time.Time `db:"release_date" json:"release_date,omitempty"`
	}

	CreateAlbumResponse struct {
		ID *uuid.UUID `json:"id,omitempty"`
	}
)

type (
	GetAlbumRequest struct {
		ID *uuid.UUID `db:"id" json:"id,omitempty" path:"id"`
	}

	GetAlbumResponse struct {
		Album *Album `json:"album,omitempty"`
	}
)

type (
	ListAlbumsRequest struct {
		Limit     *int    `db:"limit"        query:"limit"`
		Offset    *int    `db:"offset"       query:"offset"`
		SortField *string `db:"sort_field"   query:"sort_field"`
		SortOrder *string `db:"sort_order"   query:"sort_order"`
	}

	ListAlbumsResponse struct {
		Albums []*Album `json:"albums,omitempty"`
	}
)

type UpdateAlbumRequest struct {
	ID          *uuid.UUID `db:"id"           json:"-" path:"id"`
	Title       *string    `db:"title"        json:"title,omitempty"`
	Description *string    `db:"description"  json:"description,omitempty"`
	ReleaseDate *time.Time `db:"release_date" json:"release_date,omitempty"`
}

type DeleteAlbumRequest struct {
	ID *uuid.UUID `db:"id" json:"-" path:"id"`
}
