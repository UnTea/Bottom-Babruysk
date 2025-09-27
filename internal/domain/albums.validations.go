package domain

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	validatron "github.com/untea/bottom_babruysk/internal/domain/validation"
)

var (
	sortableAlbumFields = validatron.NewSet(
		"title",
		"release_date",
		"created_at",
		"updated_at",
	)
)

func (r *CreateAlbumRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Title, validation.Required),
		validation.Field(&r.Description, validation.Required),
		validation.Field(&r.ReleaseDate, validation.Required),
	)
}

func (r *GetAlbumRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ID, validation.Required),
	)
}

func (r *ListAlbumsRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Limit, validation.Required, validation.Min(1)),
		validation.Field(&r.Offset, validation.Required, validation.Min(0)),
		validation.Field(&r.SortField, validation.Required, validatron.InStringsPtr(sortableAlbumFields, "sort_field")),
		validation.Field(&r.SortOrder, validation.Required, validatron.InStringsPtr(validatron.SortOrders, "sort_order")),
	)
}

func (r *UpdateAlbumRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ID, validation.Required),
	)
}

func (r *DeleteAlbumRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ID, validation.Required),
	)
}
