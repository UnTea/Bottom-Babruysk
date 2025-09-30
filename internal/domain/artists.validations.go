package domain

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	validatron "github.com/untea/bottom_babruysk/internal/domain/validation"
)

var (
	SortableArtistFields = validatron.NewSet(
		"name",
		"created_at",
		"updated_at",
	)
)

func (r *CreateArtistRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Bio, validation.Required),
	)
}

func (r *GetArtistRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ID, validation.Required),
	)
}

func (r *ListArtistsRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Limit, validation.Required, validation.Min(1)),
		validation.Field(&r.Offset, validation.Required, validation.Min(0)),
		validation.Field(&r.SortField, validation.Required, validatron.InStringsPtr(SortableArtistFields, "sort_field")),
		validation.Field(&r.SortOrder, validation.Required, validatron.InStringsPtr(validatron.SortOrders, "sort_order")),
	)
}

func (r *UpdateArtistRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ID, validation.Required),
		validation.Field(&r.Name, validation.When(r.Name != nil, validation.Required)),
		validation.Field(&r.Bio, validation.When(r.Bio != nil, validation.Required)),
	)
}

func (r *DeleteArtistRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ID, validation.Required),
	)
}
