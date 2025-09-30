package connect

import (
	"errors"

	"connectrpc.com/connect"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/uuid"

	"github.com/untea/bottom_babruysk/internal/repository/postgres"
)

// Ptr если нужно получить указатель на значение.
func Ptr[T any](v T) *T {
	return &v
}

// PtrIfNonZero если нужно отделить "nil" value от "zero" type value, ноль трактуется как значение не задано.
func PtrIfNonZero[T comparable](v T) *T {
	var z T

	if v == z {
		return nil
	}

	return &v
}

// ValueOrZero если нужно получить само значение, а nil трактовать как ноль типа.
func ValueOrZero[T any](p *T) (z T) {
	if p != nil {
		z = *p
	}

	return
}

// ValueOK если нужно знать, пришло ли значение.
func ValueOK[T any](p *T) (T, bool) {
	if p == nil {
		var z T

		return z, false
	}

	return *p, true
}

// ValueOr если нужно значение с явным дефолтом, если nil.
func ValueOr[T any](p *T, def T) T {
	if p != nil {
		return *p
	}

	return def
}

// UUIDToString конвертация uuid.UUID к string.
//
//	В случае uuid.Nil вернёт пустую строку.
func UUIDToString(u uuid.UUID) string {
	if u == uuid.Nil {
		return ""
	}

	return u.String()
}

// UUIDPtrToString конвертация *uuid.UUID к string.
//
//	В случае nil или uuid.Nil возвращается пустая строка.
func UUIDPtrToString(u *uuid.UUID) string {
	if u == nil || *u == uuid.Nil {
		return ""
	}

	return u.String()
}

// StringToUUID конвертация строки к uuid.UUID.
//
//	В случае пустой или невалидной строки вернёт uuid.Nil.
func StringToUUID(s string) uuid.UUID {
	if s == "" {
		return uuid.Nil
	}

	u, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil
	}

	return u
}

// StringToUUIDPtr конвертация строки к *uuid.UUID.
//
//	В случае пустой или невалидная строка вернёт nil.
func StringToUUIDPtr(s string) *uuid.UUID {
	u, err := uuid.Parse(s)
	if err != nil {
		return nil
	}

	return &u
}

func toConnectErr(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, postgres.ErrNotFound) {
		return connect.NewError(connect.CodeNotFound, err)
	}

	var verr validation.Errors

	if errors.As(err, &verr) {
		return connect.NewError(connect.CodeInvalidArgument, err)
	}

	return connect.NewError(connect.CodeInternal, err)
}
