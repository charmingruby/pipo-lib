// Package redis provides a Redis implementation of the Broker.
package redis

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis"
)

const (
	batchSize = 100
	blockTime = 5 * time.Second
)

var (
	// ErrInvalidMessageDataType is returned when the message data type is invalid.
	ErrInvalidMessageDataType = errors.New("invalid message data type")
)

// Client manages a Redis connection.
type Client struct {
	conn *redis.Client
}

// NewClient constructs a new Redis client.
//
// url: The URL of the Redis server. e.g. "redis://localhost:6379"
//
// Returns a new Redis client and an error if the connection fails.
func NewClient(url string) (*Client, error) {
	conn := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: "",
		DB:       0,
	})

	if _, err := conn.Ping().Result(); err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
	}, nil
}

// Ping pings the Redis server.
//
// ctx: The context of the request.
//
// Returns an error if the ping fails.
func (c *Client) Ping(_ context.Context) error {
	return c.conn.Ping().Err()
}

// Close closes the Redis connection.
//
// Returns an error if the close fails.
func (c *Client) Close() error {
	return c.conn.Close()
}

// Stream manages a Redis stream.
type Stream struct {
	client *Client
}

// NewStream constructs a new Redis stream.
//
// client: The Redis client.
//
// Returns a new Redis stream.
func NewStream(client *Client) *Stream {
	return &Stream{
		client: client,
	}
}

// Publish publishes a message to a topic.
//
// ctx: The context of the request.
// topic: The topic to publish the message to.
// message: The message to publish in bytes.
//
// Returns an error if the message is not published.
func (s *Stream) Publish(_ context.Context, topic string, message []byte) error {
	if _, err := s.client.conn.XAdd(&redis.XAddArgs{
		Stream: topic,
		Values: map[string]interface{}{
			"data": message,
		},
	}).Result(); err != nil {
		return err
	}

	return nil
}

// Subscribe subscribes to a topic and calls a handler for each message.
//
// ctx: The context of the request.
// topic: The topic to subscribe to.
// handler: The handler to call for each message.
//
// Returns an error if the subscription fails.
func (s *Stream) Subscribe(ctx context.Context, topic string, handler func(message []byte) error) error {
	lastID := "$"

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			messages, shouldContinue, err := s.readStreamMessages(topic, lastID)
			if err != nil {
				return err
			}

			if shouldContinue {
				continue
			}

			for _, message := range messages {
				data, ok := message.Values["data"].(string)
				if !ok {
					return ErrInvalidMessageDataType
				}

				if err := handler([]byte(data)); err != nil {
					return err
				}

				lastID = message.ID
			}
		}
	}
}

// readStreamMessages reads messages from a stream.
//
// topic: The topic to read messages from.
// lastID: The last message ID to read from.
//
// Returns a list of messages, a boolean indicating if the stream should continue, and an error if the read fails.
func (s *Stream) readStreamMessages(topic string, lastID string) ([]redis.XMessage, bool, error) {
	streams, err := s.client.conn.XRead(&redis.XReadArgs{
		Streams: []string{topic, lastID},
		Count:   batchSize,
		Block:   blockTime,
	}).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, true, nil
		}

		return nil, false, err
	}

	if len(streams) == 0 || len(streams[0].Messages) == 0 {
		return nil, true, nil
	}

	return streams[0].Messages, false, nil
}
