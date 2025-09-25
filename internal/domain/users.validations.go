package domain

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

var sortableUserFields = map[string]struct{}{
	"email":        {},
	"display_name": {},
	"role":         {},
	"created_at":   {},
	"updated_at":   {},
}

func (r *CreateUserRequest) Validate() error {
	return validation.ValidateStruct(
		r,
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.PasswordHash, validation.Required),
		validation.Field(&r.DisplayName, validation.Required),
		validation.Field(&r.Role,
			validation.When(r.Role != nil,
				validation.In(UserRoleUser, UserRoleAdmin).
					Error("must be one of: user, admin"),
			),
		),
	)
}

func (r *GetUserRequest) Validate() error {
	return validation.ValidateStruct(
		r,
		validation.Field(&r.ID, validation.Required),
	)
}

func (r *GetListUserRequest) Normalize() {
	if r.Limit == nil || *r.Limit <= 0 {
		v := 50
		r.Limit = &v
	}

	if r.Offset == nil || *r.Offset <= 0 {
		v := 1
		r.Offset = &v
	}

	if r.SortField == nil || *r.SortField == "" {
		v := "created_at"
		r.SortField = &v
	} else {
		f := strings.ToLower(strings.TrimSpace(*r.SortField))
		r.SortField = &f
	}

	if r.SortOrder == nil || *r.SortOrder == "" {
		v := "desc"
		r.SortOrder = &v
	} else {
		o := strings.ToLower(strings.TrimSpace(*r.SortOrder))
		r.SortOrder = &o
	}
}

func (r *GetListUserRequest) Validate() error {
	r.Normalize()

	return validation.ValidateStruct(
		r,
		validation.Field(&r.Limit, validation.Required, validation.Min(1)),
		validation.Field(&r.Offset, validation.Required, validation.Min(1)),
		validation.Field(&r.SortField, validation.Required, validation.By(func(v any) error {
			f := *v.(*string)
			if _, ok := sortableUserFields[f]; !ok {
				return validation.NewError("validation", "invalid sort_field")
			}

			return nil
		})),
		validation.Field(&r.SortOrder, validation.Required, validation.By(func(v any) error {
			o := *v.(*string)
			if o != "asc" && o != "desc" {
				return validation.NewError("validation", "sort_order must be asc or desc")
			}

			return nil
		})),
		validation.Field(&r.Role,
			validation.When(r.Role != nil,
				validation.In(UserRoleUser, UserRoleAdmin).
					Error("role must be user or admin"),
			),
		),
	)
}

func (r *UpdateUserRequest) Validate() error {
	if r.DisplayName == nil && r.Role == nil {
		return validation.Errors{
			"body": validation.NewError("validation", "at least one of display_name or role must be provided"),
		}
	}

	return validation.ValidateStruct(
		r,
		validation.Field(&r.ID, validation.Required),
		validation.Field(&r.DisplayName,
			validation.When(r.DisplayName != nil, validation.Required),
		),
		validation.Field(&r.Role,
			validation.When(r.Role != nil,
				validation.In(UserRoleUser, UserRoleAdmin).
					Error("must be one of: user, admin"),
			),
		),
	)
}

func (r *DeleteUserRequest) Validate() error {
	return validation.ValidateStruct(
		r,
		validation.Field(&r.ID, validation.Required),
	)
}
