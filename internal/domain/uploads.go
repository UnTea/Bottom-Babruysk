package domain

import (
	"time"

	"github.com/google/uuid"
)

type UploadStatus string

const (
	UploadStatusPending    UploadStatus = "pending"
	UploadStatusProcessing UploadStatus = "processing"
	UploadStatusDone       UploadStatus = "done"
	UploadStatusFailed     UploadStatus = "failed"
)

type Uploads struct {
	ID        *uuid.UUID    `db:"id"`
	OwnerID   *uuid.UUID    `db:"owner_id"`
	Filename  *string       `db:"filename"`
	S3Key     *string       `db:"s3_key"`
	Mime      *string       `db:"mime"`
	Size      *int64        `db:"size"`
	Status    *UploadStatus `db:"status"`
	CreatedAt *time.Time    `db:"created_at"`
	UpdatedAt *time.Time    `db:"updated_at"`
}
