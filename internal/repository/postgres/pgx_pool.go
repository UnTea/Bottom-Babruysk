package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxPool struct {
	pool *pgxpool.Pool
}

func (d *PgxPool) Close() {
	if d.pool != nil {
		d.pool.Close()
	}
}

func (d *PgxPool) Query(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error) {
	return d.pool.Query(ctx, sql, arguments...)
}

func (d *PgxPool) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return d.pool.Exec(ctx, sql, arguments...)
}
