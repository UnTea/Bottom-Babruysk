package domain

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	validatron "github.com/untea/bottom_babruysk/internal/domain/validation"
)

var (
	playlistSortableFields = validatron.NewSet(
		"title",
		"created_at",
		"updated_at",
	)
)

func (r *CreatePlaylistRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.OwnerID, validation.Required),
		validation.Field(&r.Title, validation.Required),
		validation.Field(&r.Description, validation.Required),
		validation.Field(&r.Visibility, validation.When(r.Visibility != nil, validatron.InSetPtr(visibilitySet))),
	)
}

func (r *GetPlaylistRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ID, validation.Required),
	)
}

func (r *ListPlaylistsRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Limit, validation.Required, validation.Min(1)),
		validation.Field(&r.Offset, validation.Required, validation.Min(0)),
		validation.Field(&r.SortField, validation.Required, validatron.InStringsPtr(playlistSortableFields, "sort_field")),
		validation.Field(&r.SortOrder, validation.Required, validatron.InStringsPtr(validatron.SortOrders, "sort_order")),
		validation.Field(&r.Visibility, validation.When(r.Visibility != nil, validatron.InSetPtr(visibilitySet))),
	)
}

func (r *UpdatePlaylistRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ID, validation.Required),
		validation.Field(&r.Title, validation.When(r.Title != nil, validation.Required)),
		validation.Field(&r.Description, validation.When(r.Description != nil, validation.Required)),
		validation.Field(&r.Visibility, validation.When(r.Visibility != nil, validatron.InSetPtr(visibilitySet))),
	)
}

func (r *DeletePlaylistRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ID, validation.Required),
	)
}
