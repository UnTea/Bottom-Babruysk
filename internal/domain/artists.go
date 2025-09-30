package domain

import (
	"time"

	"github.com/google/uuid"
)

type Artist struct {
	ID        *uuid.UUID `db:"id"         json:"id,omitempty"`
	Name      *string    `db:"name"       json:"name,omitempty"`
	Bio       *string    `db:"bio"        json:"bio,omitempty"`
	CreatedAt *time.Time `db:"created_at" json:"createdAt,omitempty"`
	UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
}

type (
	CreateArtistRequest struct {
		Name *string `db:"name" json:"name,omitempty"`
		Bio  *string `db:"bio"  json:"bio,omitempty"`
	}

	CreateArtistResponse struct {
		ID *uuid.UUID `json:"id,omitempty"`
	}
)

type (
	GetArtistRequest struct {
		ID *uuid.UUID `db:"id" json:"id,omitempty" path:"id"`
	}

	GetArtistResponse struct {
		Artist *Artist `json:"artist,omitempty"`
	}
)

type (
	ListArtistsRequest struct {
		Name      *string `db:"name"       query:"name"`
		Bio       *string `db:"bio"        query:"bio"`
		Limit     *int    `db:"limit"      query:"limit"`
		Offset    *int    `db:"offset"     query:"offset"`
		SortField *string `db:"sort_field" query:"sort_field"`
		SortOrder *string `db:"sort_order" query:"sort_order"`
	}

	ListArtistsResponse struct {
		Artists []*Artist `json:"artists,omitempty"`
	}
)

type UpdateArtistRequest struct {
	ID   *uuid.UUID `db:"id"   json:"-" path:"id"`
	Name *string    `db:"name" json:"name,omitempty"`
	Bio  *string    `db:"bio"  json:"bio,omitempty"`
}

type DeleteArtistRequest struct {
	ID *uuid.UUID `db:"id" json:"-" path:"id"`
}
