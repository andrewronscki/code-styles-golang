package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/andre/code-styles-golang/pkg/messaging/options"
	"github.com/rabbitmq/amqp091-go"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type Consumer[T any] interface {
	Consume(ctx context.Context) (<-chan Message[T], error)
}

type RabbitMQConsumer[T any] struct {
	uri            string
	queue          string
	consumerID     string
	autoAck        bool
	exclusive      bool
	noLocal        bool
	noWait         bool
	datadogEnabled bool
	args           amqp091.Table
}

func CreateConsumer[T any](opts ...*options.ConsumerOptions) (Consumer[T], error) {
	opt := options.MergeConsumerOptions(opts...)

	c := &RabbitMQConsumer[T]{
		uri:            opt.URI,
		queue:          opt.Queue,
		consumerID:     opt.ConsumerID,
		autoAck:        opt.AutoAck,
		exclusive:      opt.Exclusive,
		noLocal:        opt.NoLocal,
		noWait:         opt.NoWait,
		datadogEnabled: opt.DatadogIntegrationEnabled,
		args:           opt.Args,
	}

	return c, nil
}

func (c *RabbitMQConsumer[T]) Consume(ctx context.Context) (<-chan Message[T], error) {
	messagesCh := make(chan Message[T])

	go func(ch chan<- Message[T]) {
		defer close(ch)

		for {
			if ctx.Err() != nil {
				return
			}

			if err := c.consume(ctx, ch); err != nil {
				log.Printf("consumer %s stopped, reconnecting in %d seconds", c.consumerID, 5)
				time.Sleep(5 * time.Second)
			} else {
				return
			}
		}
	}(messagesCh)

	return messagesCh, nil
}

func (c *RabbitMQConsumer[T]) consume(ctx context.Context, messagesCh chan<- Message[T]) error {
	connection, err := connect(c.uri)

	if err != nil {
		return fmt.Errorf("connection to broker failed: %s", err)
	}

	channel, err := connection.Channel()
	closeCh := channel.NotifyClose(make(chan *amqp091.Error))

	if err != nil {
		return fmt.Errorf("channel creation failed: %s", err)
	}

	defer connection.Close()
	defer channel.Close()

	log.Printf("consumer %s is waiting messages from %s", c.consumerID, c.queue)

	deliveries, err := channel.Consume(
		c.queue,
		c.consumerID,
		c.autoAck,
		c.exclusive,
		c.noLocal,
		c.noWait,
		c.args,
	)

	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			channel.Cancel(c.consumerID, false)
			return nil
		case err := <-closeCh:
			return err
		case delivery, ok := <-deliveries:
			if !ok {
				return fmt.Errorf("messages channel closed")
			}

			if c.datadogEnabled {
				carrier := tracer.TextMapCarrier{}

				for k, v := range delivery.Headers {
					carrier.Set(k, fmt.Sprint(v))
				}

				spanCtx, err := tracer.Extract(carrier)

				var span tracer.Span
				if err == nil {
					span = tracer.StartSpan("rabbitmq.consume", tracer.ChildOf(spanCtx))
				} else {
					span = tracer.StartSpan("rabbitmq.consume")
				}

				span.SetTag("span.kind", "consumer")
				span.SetTag("messaging.system", "rabbitmq")
				span.SetTag("messaging.destination", c.queue)
				span.SetTag("messaging.destination_kind", "queue")
				span.SetTag("messaging.protocol", "amqp")
				span.SetTag("messaging.protocol_version", "0.9.1")

				if uri, err := replacePassword(c.uri); err == nil {
					span.SetTag("messaging.url", uri)
				}

				span.SetTag("messaging.message_payload_size", len(delivery.Body))
				span.SetTag("messaging.operation", "receive")
				span.SetTag("messaging.consumer_id", c.consumerID)

				if delivery.MessageId != "" {
					span.SetTag("messaging.message_id", delivery.MessageId)
				}

				if delivery.CorrelationId != "" {
					span.SetTag("messaging.conversation_id", delivery.CorrelationId)
				}

				msgCtx := tracer.ContextWithSpan(ctx, span)
				var content T

				json.Unmarshal(delivery.Body, &content)

				message := Message[T]{
					Headers:     delivery.Headers,
					Content:     content,
					ContentType: delivery.ContentType,
					Context:     msgCtx,
					RawDelivery: delivery,
				}

				span.Finish()
				messagesCh <- message

			} else {
				var content T

				json.Unmarshal(delivery.Body, &content)

				message := Message[T]{
					Headers:     delivery.Headers,
					Content:     content,
					ContentType: delivery.ContentType,
					Context:     ctx,
					RawDelivery: delivery,
				}

				messagesCh <- message
			}
		}
	}
}
