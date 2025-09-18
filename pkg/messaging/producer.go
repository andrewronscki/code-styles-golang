package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/andre/code-styles-golang/pkg/messaging/options"
	"github.com/rabbitmq/amqp091-go"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type Producer[T any] interface {
	Produce(ctx context.Context, msg Message[T]) error
}

type RabbitMQProducer[T any] struct {
	uri             string
	destination     string
	destinationKind string
	routingKey      string
	mandatory       bool
	immediate       bool
	datadogEnabled  bool
}

func CreateProducer[T any](opts ...*options.ProducerOptions) (Producer[T], error) {
	opt := options.MergeProducerOptions(opts...)

	p := &RabbitMQProducer[T]{
		uri:             opt.URI,
		destination:     opt.Destination,
		destinationKind: opt.DestinationKind,
		routingKey:      opt.RoutingKey,
		mandatory:       opt.Mandatory,
		immediate:       opt.Immediate,
		datadogEnabled:  opt.DatadogIntegrationEnabled,
	}

	return p, nil
}

func (p *RabbitMQProducer[T]) Produce(ctx context.Context, msg Message[T]) error {
	connection, err := connect(p.uri)

	if err != nil {
		return fmt.Errorf("connection to broker failed: %s", err)
	}

	channel, err := connection.Channel()

	if err != nil {
		return fmt.Errorf("channel creation failed: %s", err)
	}

	defer connection.Close()
	defer channel.Close()

	body, err := json.Marshal(msg.Content)

	if err != nil {
		return err
	}

	publishing := amqp091.Publishing{
		Headers:     msg.Headers,
		ContentType: msg.ContentType,
		Body:        body,
	}

	if p.datadogEnabled {
		carrier := tracer.TextMapCarrier{}

		parent, ok := tracer.SpanFromContext(ctx)

		var span tracer.Span

		if !ok {
			span, ctx = tracer.StartSpanFromContext(ctx, "rabbitmq.publish")
		} else {
			span, ctx = tracer.StartSpanFromContext(ctx, "rabbitmq.publish", tracer.ChildOf(parent.Context()))
		}

		tracer.Inject(span.Context(), carrier)

		if publishing.Headers == nil {
			publishing.Headers = amqp091.Table{}
		}

		carrier.ForeachKey(func(key, val string) error {
			publishing.Headers[key] = val
			return nil
		})

		span.SetTag("span.kind", "producer")
		span.SetTag("messaging.system", "rabbitmq")
		span.SetTag("messaging.destination", p.destination)
		span.SetTag("messaging.destination_kind", p.destinationKind)

		if p.routingKey != "" {
			span.SetTag("messaging.routing_key", p.routingKey)
		}

		span.SetTag("messaging.protocol", "amqp")
		span.SetTag("messaging.protocol_version", "0.9.1")

		if uri, err := replacePassword(p.uri); err == nil {
			span.SetTag("messaging.url", uri)
		}

		span.SetTag("messaging.message_payload_size", len(publishing.Body))
		span.SetTag("messaging.operation", "send")

		if publishing.MessageId != "" {
			span.SetTag("messaging.message_id", publishing.MessageId)
		}

		if publishing.CorrelationId != "" {
			span.SetTag("messaging.conversation_id", publishing.CorrelationId)
		}

		err = channel.PublishWithContext(
			ctx,
			p.destination,
			p.routingKey,
			p.mandatory,
			p.immediate,
			publishing,
		)

		span.Finish()

		return err
	}

	err = channel.PublishWithContext(
		ctx,
		p.destination,
		p.routingKey,
		p.mandatory,
		p.immediate,
		publishing,
	)

	return err
}
