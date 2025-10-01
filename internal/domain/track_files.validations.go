package domain

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	validatron "github.com/untea/bottom_babruysk/internal/domain/validation"
)

var (
	formatSet = validatron.NewSet(
		"mp3",
		"mp4",
		"flac",
		"wav",
		"ogg",
		"ts",
		"m4a",
		"webm",
		"aac",
	)

	codecSet = validatron.NewSet(
		"wav",
		"wv",
		"wvc",
		"flac",
		"alac",
		"lpac",
		"ltac",
		"off",
		"ofr",
		"ofs",
		"thd",
		"ape",
		"shn",
		"opus",
		"vorbis",
		"pcm",
		"mp3",
		"aac",
		"wma",
		"ogg",
	)

	sortableTrackFileFields = validatron.NewSet(
		"filename",
		"format",
		"codec",
		"bitrate",
		"sample_rate",
		"channels",
		"size",
		"duration",
		"uploaded_at",
		"created_at",
		"updated_at",
	)
)

func (r *CreateTrackFileRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.TrackID, validation.Required),
		validation.Field(&r.Filename, validation.Required),
		validation.Field(&r.Mime, validation.Required),
		validation.Field(&r.Format, validation.Required, validatron.InStringsPtr(formatSet, "format")),
		validation.Field(&r.Codec, validation.Required, validatron.InStringsPtr(codecSet, "codec")),
		validation.Field(&r.Bitrate, validation.Required, validation.Min(1)),
		validation.Field(&r.SampleRate, validation.Required, validation.Min(1)),
		validation.Field(&r.Channels, validation.Required, validation.Min(1)),
		validation.Field(&r.Size, validation.Required, validation.Min(1)),
		validation.Field(&r.Duration, validation.Required),
		validation.Field(&r.Checksum, validation.Required),
		validation.Field(&r.UploadedAt, validation.Required),
	)
}

func (r *GetTrackFileRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ID, validation.Required),
		validation.Field(&r.TrackID, validation.Required),
	)
}

func (r *ListTrackFilesRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.TrackID, validation.Required),
		validation.Field(&r.Limit, validation.Min(1)),
		validation.Field(&r.Offset, validation.Min(0)),
		validation.Field(&r.SortField, validatron.InStringsPtr(sortableTrackFileFields, "sort_field")),
		validation.Field(&r.SortOrder, validatron.InStringsPtr(validatron.SortOrders, "sort_order")),
	)
}

func (r *UpdateTrackFileRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ID, validation.Required),
		validation.Field(&r.TrackID, validation.Required),
		validation.Field(&r.Format, validation.When(r.Format != nil, validatron.InStringsPtr(formatSet, "format"))),
		validation.Field(&r.Codec, validation.When(r.Codec != nil, validatron.InStringsPtr(codecSet, "codec"))),
		validation.Field(&r.Bitrate, validation.When(r.Bitrate != nil, validation.Min(1))),
		validation.Field(&r.SampleRate, validation.When(r.SampleRate != nil, validation.Min(1))),
		validation.Field(&r.Channels, validation.When(r.Channels != nil, validation.Min(1))),
		validation.Field(&r.Size, validation.When(r.Size != nil, validation.Min(int64(1)))),
	)
}

func (r *DeleteTrackFileRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ID, validation.Required),
		validation.Field(&r.TrackID, validation.Required),
	)
}
