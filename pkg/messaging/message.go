package messaging

import (
	"context"

	"github.com/rabbitmq/amqp091-go"
)

type Message[T any] struct {
	Headers     amqp091.Table
	Content     T
	ContentType string
	Context     context.Context
	RawDelivery amqp091.Delivery
}
