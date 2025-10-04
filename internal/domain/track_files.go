package domain

import (
	"time"

	"github.com/google/uuid"
)

type Codec string

const (
	CodecUnspecified Codec = "unspecified"

	// Lossless
	CodecWAV  Codec = "wav"
	CodecFLAC Codec = "flac"
	CodecALAC Codec = "alac"
	CodecAPE  Codec = "ape"
	CodecSHN  Codec = "shn"

	// Lossy
	CodecMP3    Codec = "mp3"
	CodecAAC    Codec = "aac"
	CodecOPUS   Codec = "opus"
	CodecVORBIS Codec = "vorbis"
	CodecWMA    Codec = "wma"
)

type Format string

const (
	FormatUnspecified Format = "unspecified"
	FormatMP3         Format = "mp3"
	FormatMP4         Format = "mp4"
	FormatM4A         Format = "m4a"
	FormatOGG         Format = "ogg"
	FormatFLAC        Format = "flac"
	FormatWAV         Format = "wav"
	FormatWEBM        Format = "webm"
	FormatAAC         Format = "aac"
)

type TrackFile struct {
	ID         *uuid.UUID     `db:"id"          json:"id,omitempty"`
	TrackID    *uuid.UUID     `db:"track_id"    json:"track_id,omitempty"`
	Filename   *string        `db:"filename"    json:"filename,omitempty"`
	S3Key      *string        `db:"s3_key"      json:"s3_key,omitempty"`
	Mime       *string        `db:"mime"        json:"mime,omitempty"`
	Format     *Format        `db:"format"      json:"format,omitempty"`
	Codec      *Codec         `db:"codec"       json:"codec,omitempty"`
	Bitrate    *int           `db:"bitrate"     json:"bitrate,omitempty"`
	SampleRate *int           `db:"sample_rate" json:"sample_rate,omitempty"`
	Channels   *int           `db:"channels"    json:"channels,omitempty"`
	Size       *int64         `db:"size"        json:"size,omitempty"`
	Duration   *time.Duration `db:"duration"    json:"duration,omitempty"`
	Checksum   []byte         `db:"checksum"    json:"checksum,omitempty"`
	CreatedAt  *time.Time     `db:"created_at"  json:"created_at,omitempty"`
	UpdatedAt  *time.Time     `db:"updated_at"  json:"updated_at,omitempty"`
	UploadedAt *time.Time     `db:"uploaded_at" json:"uploaded_at,omitempty"`
}

type (
	CreateTrackFileRequest struct {
		TrackID    *uuid.UUID     `db:"track_id"    json:"track_id,omitempty"`
		Filename   *string        `db:"filename"    json:"filename,omitempty"`
		S3Key      *string        `db:"s3_key"      json:"s3_key,omitempty"`
		Mime       *string        `db:"mime"        json:"mime,omitempty"`
		Format     *Format        `db:"format"      json:"format,omitempty"`
		Codec      *Codec         `db:"codec"       json:"codec,omitempty"`
		Bitrate    *int           `db:"bitrate"     json:"bitrate,omitempty"`
		SampleRate *int           `db:"sample_rate" json:"sample_rate,omitempty"`
		Channels   *int           `db:"channels"    json:"channels,omitempty"`
		Size       *int64         `db:"size"        json:"size,omitempty"`
		Duration   *time.Duration `db:"duration"    json:"duration,omitempty"`
		Checksum   []byte         `db:"checksum"    json:"checksum,omitempty"`
		UploadedAt *time.Time     `db:"uploaded_at" json:"uploaded_at,omitempty"`
	}

	CreateTrackFileResponse struct {
		ID *uuid.UUID `json:"id"`
	}
)

type (
	GetTrackFileRequest struct {
		ID      *uuid.UUID `db:"id"       json:"id,omitempty"       path:"id"`
		TrackID *uuid.UUID `db:"track_id" json:"track_id,omitempty" path:"track_id"`
	}

	GetTrackFileResponse struct {
		TrackFile *TrackFile `json:"trackfile,omitempty"`
	}
)

type (
	ListTrackFilesRequest struct {
		TrackID   *uuid.UUID `db:"track_id"   query:"track_id"`
		Limit     *int       `db:"limit"      query:"limit"`
		Offset    *int       `db:"offset"     query:"offset"`
		SortField *string    `db:"sort_field" query:"sort_field"`
		SortOrder *string    `db:"sort_order" query:"sort_order"`
	}

	ListTrackFilesResponse struct {
		TrackFiles []*TrackFile `json:"trackfiles,omitempty"`
	}
)

type UpdateTrackFileRequest struct {
	ID         *uuid.UUID     `db:"id"       json:"-" path:"id"`
	TrackID    *uuid.UUID     `db:"track_id" json:"-" path:"track_id"`
	Filename   *string        `db:"filename"    json:"filename,omitempty"`
	S3Key      *string        `db:"s3_key"      json:"s3_key,omitempty"`
	Mime       *string        `db:"mime"        json:"mime,omitempty"`
	Format     *Format        `db:"format"      json:"format,omitempty"`
	Codec      *Codec         `db:"codec"       json:"codec,omitempty"`
	Bitrate    *int           `db:"bitrate"     json:"bitrate,omitempty"`
	SampleRate *int           `db:"sample_rate" json:"sample_rate,omitempty"`
	Channels   *int           `db:"channels"    json:"channels,omitempty"`
	Size       *int64         `db:"size"        json:"size,omitempty"`
	Duration   *time.Duration `db:"duration" json:"duration,omitempty"`
	Checksum   *string        `db:"checksum"    json:"checksum,omitempty"`
}

type DeleteTrackFileRequest struct {
	ID      *uuid.UUID `db:"id"       json:"-" path:"id"`
	TrackID *uuid.UUID `db:"track_id" json:"-" path:"track_id"`
}
