package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	driver Driver
	config Config
}

type Config struct {
	ConnectionString string
	Timeout          time.Duration
}

func New(ctx context.Context, config Config) (*Client, error) {
	pollConf, err := pgxpool.ParseConfig(config.ConnectionString)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, pollConf)
	if err != nil {
		return nil, err
	}

	return &Client{
		driver: &PgxPool{
			pool: pool,
		},
		config: config,
	}, nil
}

func (c *Client) Close() {
	if c.driver != nil {
		c.driver.Close()
	}
}

func (c *Client) QueryTimeout() time.Duration {
	return c.config.Timeout
}

func (c *Client) Driver() Driver {
	return c.driver
}
