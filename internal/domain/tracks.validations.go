package domain

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"

	validatron "github.com/untea/bottom_babruysk/internal/domain/validation"
)

var (
	visibilitySet = validatron.NewSet(VisibilityPrivate, VisibilityUnlisted, VisibilityPublic)

	SortableTrackFields = validatron.NewSet(
		"title",
		"subtitle",
		"visibility",
		"uploaded_at",
		"created_at",
		"updated_at",
	)
)

func (r *CreateTrackRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.UploaderID, validation.Required),
		validation.Field(&r.Title, validation.Required),
		validation.Field(&r.Subtitle, validation.Required),
		validation.Field(&r.Description, validation.Required),
		validation.Field(&r.Duration, validation.Required),
		validation.Field(&r.Visibility, validatron.InSetPtr(visibilitySet)),
		validation.Field(&r.UploadedAt, validation.Required),
	)
}

func (r *GetTrackRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ID, validation.Required),
	)
}

func (r *ListTracksRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Limit, validation.Required, validation.Min(1)),
		validation.Field(&r.Offset, validation.Required, validation.Min(0)),
		validation.Field(&r.SortField, validation.Required, validatron.InStringsPtr(SortableTrackFields, "sort_field")),
		validation.Field(&r.SortOrder, validation.Required, validatron.InStringsPtr(validatron.SortOrders, "sort_order")),
		validation.Field(&r.Visibility, validatron.InSetPtr(visibilitySet)),
	)
}

func (r *UpdateTrackRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ID, validation.Required),
		validation.Field(&r.Title, validation.When(r.Title != nil, validation.Required)),
		validation.Field(&r.Subtitle, validation.When(r.Subtitle != nil, validation.Required)),
		validation.Field(&r.Description, validation.When(r.Description != nil, validation.Required)),
		validation.Field(&r.Visibility, validatron.InSetPtr(visibilitySet)),
	)
}

func (r *DeleteTrackRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ID, validation.Required),
	)
}
