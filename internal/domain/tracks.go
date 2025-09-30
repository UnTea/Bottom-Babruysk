package domain

import (
	"time"

	"github.com/google/uuid"
)

type Visibility string

const (
	VisibilityUnspecified Visibility = "unspecified"
	VisibilityPrivate     Visibility = "private"
	VisibilityUnlisted    Visibility = "unlisted"
	VisibilityPublic      Visibility = "public"
)

type Track struct {
	ID          *uuid.UUID     `db:"id"          json:"id,omitempty"`
	UploaderID  *uuid.UUID     `db:"uploader_id" json:"uploaderID,omitempty"`
	Title       *string        `db:"title"       json:"title,omitempty"`
	Subtitle    *string        `db:"subtitle"    json:"subtitle,omitempty"`
	Description *string        `db:"description" json:"description,omitempty"`
	Duration    *time.Duration `db:"duration"    json:"duration,omitempty"`
	Visibility  *Visibility    `db:"visibility"  json:"visibility,omitempty"`
	CreatedAt   *time.Time     `db:"created_at"  json:"createdAt,omitempty"`
	UpdatedAt   *time.Time     `db:"updated_at"  json:"updatedAt,omitempty"`
	UploadedAt  *time.Time     `db:"uploaded_at" json:"uploadedAt,omitempty"`
}

type (
	CreateTrackRequest struct {
		UploaderID  *uuid.UUID     `db:"uploader_id" json:"uploader_id,omitempty"`
		Title       *string        `db:"title"       json:"title,omitempty"`
		Subtitle    *string        `db:"subtitle"    json:"subtitle,omitempty"`
		Description *string        `db:"description" json:"description,omitempty"`
		Duration    *time.Duration `db:"duration"    json:"duration,omitempty"`
		Visibility  *Visibility    `db:"visibility"  json:"visibility,omitempty"`
		UploadedAt  *time.Time     `db:"uploaded_at" json:"uploaded_at,omitempty"`
	}

	CreateTrackResponse struct {
		ID *uuid.UUID `json:"track"`
	}
)

type (
	GetTrackRequest struct {
		ID *uuid.UUID `db:"id" json:"id,omitempty" path:"id"`
	}

	GetTrackResponse struct {
		Track *Track `json:"track"`
	}
)

type (
	ListTracksRequest struct {
		Limit      *int        `db:"limit"       query:"limit"`
		Offset     *int        `db:"offset"      query:"offset"`
		UploaderID *uuid.UUID  `db:"uploader_id" query:"uploader_id"`
		Visibility *Visibility `db:"visibility"  query:"visibility"`
		SortField  *string     `db:"sort_field"  query:"sort_field"`
		SortOrder  *string     `db:"sort_order"  query:"sort_order"`
	}

	ListTracksResponse struct {
		Tracks []*Track `json:"tracks,omitempty"`
	}
)

type UpdateTrackRequest struct {
	ID          *uuid.UUID  `db:"id"          json:"-" path:"id"`
	Title       *string     `db:"title"       json:"title,omitempty"`
	Subtitle    *string     `db:"subtitle"    json:"subtitle,omitempty"`
	Description *string     `db:"description" json:"description,omitempty"`
	Visibility  *Visibility `db:"visibility"  json:"visibility,omitempty"`
}

type DeleteTrackRequest struct {
	ID *uuid.UUID `db:"id" json:"-" path:"id"`
}
