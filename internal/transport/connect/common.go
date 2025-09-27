package connect

import (
	"errors"
	"time"

	"connectrpc.com/connect"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/untea/bottom_babruysk/internal/domain"
	"github.com/untea/bottom_babruysk/internal/repository"
)

func tsPtr(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}

	return timestamppb.New(t.UTC())
}

func timePtrFromTS(ts *timestamppb.Timestamp) *time.Time {
	if ts == nil {
		return nil
	}

	t := ts.AsTime().UTC()

	return &t
}

func uuidStr(u *uuid.UUID) string {
	if u == nil {
		return ""
	}

	return u.String()
}

func uuidPtrFromStr(s string) *uuid.UUID {
	if s == "" {
		return nil
	}

	u, err := uuid.Parse(s)
	if err != nil {
		return nil
	}

	return &u
}

func uuidFromStr(s string) uuid.UUID {
	u, _ := uuid.Parse(s)

	return u
}

func strOrEmpty(p *string) string {
	if p == nil {
		return ""
	}

	return *p
}

func strPtrOrNil(s string) *string {
	if s == "" {
		return nil
	}

	return &s
}

func toConnectErr(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, repository.ErrNotFound) {
		return connect.NewError(connect.CodeNotFound, err)
	}

	var verr validation.Errors

	if errors.As(err, &verr) {
		return connect.NewError(connect.CodeInvalidArgument, err)
	}

	return connect.NewError(connect.CodeInternal, err)
}

func intPtr(v int) *int {
	return &v
}

func derefUserRole(r *domain.UserRole) domain.UserRole {
	if r == nil {
		return ""
	}

	return *r
}
