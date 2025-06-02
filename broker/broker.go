// Package broker provides a generic interface for message brokers implementation.
package broker

import "context"

// Broker interface that defines the methods for a message broker.
type Broker interface {
	// Publish publishes a message to a topic.
	//
	// ctx: The context of the request.
	// topic: The topic to publish the message to.
	// message: The message to publish in bytes.
	//
	// Returns an error if the message is not published.
	Publish(ctx context.Context, topic string, message []byte) error

	// Subscribe subscribes to a topic and calls a handler for each message.
	//
	// ctx: The context of the request.
	// topic: The topic to subscribe to.
	// handler: The handler to call for each message.
	//
	// Returns an error if the subscription fails.
	Subscribe(ctx context.Context, topic string, handler func(message []byte) error) error
}
