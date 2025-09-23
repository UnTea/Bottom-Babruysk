package domain

import (
	"time"

	"github.com/google/uuid"
)

type Visibility string

const (
	VisibilityPrivate  Visibility = "private"
	VisibilityUnlisted Visibility = "unlisted"
	VisibilityPublic   Visibility = "public"
)

type Track struct {
	ID          *uuid.UUID     `db:"id"`
	UploaderID  *uuid.UUID     `db:"uploader_id"`
	Title       *string        `db:"title"`
	Subtitle    *string        `db:"subtitle"`
	Description *string        `db:"description"`
	Duration    *time.Duration `db:"duration"`
	Visibility  *Visibility    `db:"visibility"`
	CreatedAt   *time.Time     `db:"created_at"`
	UpdatedAt   *time.Time     `db:"updated_at"`
	UploadedAt  *time.Time     `db:"uploaded_at"`
}
