// Package postgres provides a Postgres connection.
package postgres

import (
	"context"

	"github.com/charmingruby/pipo-lib/logger"
	"github.com/jmoiron/sqlx"

	// Postgres driver.
	_ "github.com/lib/pq"
)

const driver = "postgres"

// Client represents a Postgres connection.
type Client struct {
	// Connection to the Postgres database.
	Conn *sqlx.DB
	// Logger.
	logger *logger.Logger
}

// New constructs a new Postgres client.
//
// logger: The logger instance.
// url: The connection URL. e.g. "postgresql://user:password@host:port/database?sslmode=disable"
//
// Returns a new Postgres client and an error if the connection fails.
func New(logger *logger.Logger, url string) (*Client, error) {
	db, err := sqlx.Connect(driver, url)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Client{Conn: db, logger: logger}, nil
}

// Ping pings the Postgres database.
//
// ctx: The context of the request.
//
// Returns an error if the ping fails.
func (c *Client) Ping(ctx context.Context) error {
	return c.Conn.PingContext(ctx)
}

// Close closes the Postgres connection.
//
// ctx: The context of the request.
//
// Returns an error if the close fails.
func (c *Client) Close(ctx context.Context) error {
	if err := c.Conn.Close(); err != nil {
		c.logger.Error("failed to close postgres connection", "error", err)
		return err
	}

	select {
	case <-ctx.Done():
		c.logger.Error("failed to close postgres connection", "error", ctx.Err())
		return ctx.Err()
	default:
		c.logger.Debug("postgres connection closed")
		return nil
	}
}
