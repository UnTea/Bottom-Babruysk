package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

var (
	ErrNotFound = errors.New("not found")
)

// FetchOne выполняет запрос и возвращает одну модель. Если строка не найдена - возвращает ErrNotFound.
func FetchOne[Model any](ctx context.Context, driver Driver, sqlQuery string, arguments ...any) (*Model, error) {
	rows, err := driver.Query(ctx, sqlQuery, arguments...)
	if err != nil {
		return nil, err
	}

	model, err := pgx.CollectOneRow[*Model](rows, pgx.RowToAddrOfStructByNameLax[Model])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return model, nil
}

// FetchMany выполняет запрос и возвращает слайс моделей.
func FetchMany[Model any](ctx context.Context, driver Driver, sqlQuery string, arguments ...any) ([]*Model, error) {
	rows, err := driver.Query(ctx, sqlQuery, arguments...)
	if err != nil {
		return nil, err
	}

	models, err := pgx.CollectRows[*Model](rows, pgx.RowToAddrOfStructByNameLax[Model])
	if err != nil {
		return nil, err
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return models, nil
}
