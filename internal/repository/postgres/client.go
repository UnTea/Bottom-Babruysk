package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	driver        Driver
	configuration Configuration
}

type Configuration struct {
	ConnectionString string
	Timeout          time.Duration
}

func New(ctx context.Context, configuration Configuration) (*Client, error) {
	pollConf, err := pgxpool.ParseConfig(configuration.ConnectionString)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, pollConf)
	if err != nil {
		return nil, err
	}

	databaseClient := &Client{
		driver: &PgxPool{
			pool: pool,
		},
		configuration: configuration,
	}

	return databaseClient, nil
}

func (c *Client) Close() {
	if c.driver != nil {
		c.driver.Close()
	}
}

func (c *Client) QueryTimeout() time.Duration {
	return c.configuration.Timeout
}

func (c *Client) Driver() Driver {
	return c.driver
}

func (c *Client) Query(ctx context.Context, sqlQuery string, arguments ...any) (pgx.Rows, error) {
	return c.driver.Query(ctx, sqlQuery, arguments...)
}

func (c *Client) Exec(ctx context.Context, sqlQuery string, arguments ...any) (pgconn.CommandTag, error) {
	return c.driver.Exec(ctx, sqlQuery, arguments...)
}
