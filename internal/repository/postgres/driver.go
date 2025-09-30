package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Driver interface {
	Query(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Close()
}

// ExecAffected helper для ExecAffected с возвратом числа изменённых строк.
func ExecAffected(ctx context.Context, driver Driver, sqlQuery string, arguments ...any) (int64, error) {
	tag, err := driver.Exec(ctx, sqlQuery, arguments...)
	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), nil
}
