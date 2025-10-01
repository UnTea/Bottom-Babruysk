package connect

import (
	"errors"

	"connectrpc.com/connect"
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/untea/bottom_babruysk/internal/repository/postgres"
)

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
