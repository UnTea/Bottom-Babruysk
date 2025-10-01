package domain

import (
	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	validatron "github.com/untea/bottom_babruysk/internal/domain/validation"
)

var (
	roleSet = validatron.NewSet(RoleUser, RoleAdmin)

	sortableUserFields = validatron.NewSet(
		"email",
		"display_name",
		"role",
		"created_at",
		"updated_at",
	)
)

func (r *CreateUserRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.PasswordHash, validation.Required),
		validation.Field(&r.DisplayName, validation.Required),
		validation.Field(&r.Role, validatron.InSetPtr(roleSet)),
	)
}

func (r *GetUserRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ID, validation.Required),
	)
}

func (r *ListUsersRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Limit, validation.Min(1)),
		validation.Field(&r.Offset, validation.Min(0)),
		validation.Field(&r.SortField, validation.Required, validatron.InStringsPtr(sortableUserFields, "sort_field")),
		validation.Field(&r.SortOrder, validation.Required, validatron.InStringsPtr(validatron.SortOrders, "sort_order")),
		validation.Field(&r.Role, validation.When(r.Role != nil, validatron.InSetPtr(roleSet))),
	)
}

func (r *UpdateUserRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ID, validation.Required),
		validation.Field(&r.DisplayName, validation.When(r.DisplayName != nil, validation.Required)),
		validation.Field(&r.Role, validation.When(r.Role != nil, validatron.InSetPtr(roleSet))),
	)
}

func (r *DeleteUserRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ID, validation.Required),
	)
}
