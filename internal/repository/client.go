package repository

import (
	"context"
	"time"

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

func New(ctx context.Context, cfg Configuration) (*Client, error) {
	pollConf, err := pgxpool.ParseConfig(cfg.ConnectionString)
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
		configuration: cfg,
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
