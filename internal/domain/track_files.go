package domain

import (
	"time"

	"github.com/google/uuid"
)

type Codec string

const (
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
	FormatMP3  Format = "mp3"
	FormatMP4  Format = "mp4"
	FormatM4A  Format = "m4a"
	FormatOGG  Format = "ogg"
	FormatFLAC Format = "flac"
	FormatWAV  Format = "wav"
	FormatWEBM Format = "webm"
	FormatAAC  Format = "aac"
)

type TrackFile struct {
	ID         *uuid.UUID     `db:"id"`
	TrackID    *uuid.UUID     `db:"track_id"`
	Filename   *string        `db:"filename"`
	S3Key      *string        `db:"s3_key,omitempty"`
	Mime       *string        `db:"mime"`
	Format     *Format        `db:"format"`
	Codec      *Codec         `db:"codec"`
	Bitrate    *int           `db:"bitrate"`
	SampleRate *int           `db:"sample_rate"`
	Channels   *int           `db:"channels"`
	Size       *int64         `db:"size"`
	Duration   *time.Duration `db:"duration"`
	Checksum   *string        `db:"checksum"`
	CreatedAt  *time.Time     `db:"created_at"`
	UpdatedAt  *time.Time     `db:"updated_at"`
	UploadedAt *time.Time     `db:"uploaded_at"`
}
