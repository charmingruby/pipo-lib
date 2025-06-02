// Package postgres provides a Postgres connection.
package postgres

import (
	"context"
	"fmt"
	"net"

	"github.com/charmingruby/pipo/lib/logger"
	"github.com/jmoiron/sqlx"

	// Postgres driver.
	_ "github.com/lib/pq"
)

// Client represents a Postgres connection.
type Client struct {
	// Connection to the Postgres database.
	Conn *sqlx.DB
	// Logger.
	logger *logger.Logger
}

// ConnectionInput represents the input for a Postgres connection.
type ConnectionInput struct {
	User         string
	Password     string
	Host         string
	Port         string
	DatabaseName string
	SSL          string
}

// New constructs a new Postgres client.
//
// logger: The logger.
// in: The connection input.
//
// Returns a new Postgres client and an error if the connection fails.
func New(logger *logger.Logger, in ConnectionInput) (*Client, error) {
	connectionString := fmt.Sprintf(
		"postgresql://%s:%s@%s/%s?sslmode=%s",
		in.User,
		in.Password,
		net.JoinHostPort(in.Host, in.Port),
		in.DatabaseName,
		in.SSL,
	)

	dbDriver := "postgres"

	db, err := sqlx.Connect(dbDriver, connectionString)
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
