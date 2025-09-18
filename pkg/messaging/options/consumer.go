package options

import "github.com/rabbitmq/amqp091-go"

// ConsumerOptions configures the behavior of an AMQP consumer.
type ConsumerOptions struct {
	// URI specifies the AMQP server endpoint.
	URI string

	// Queue specifies the name of the queue from which to consume messages.
	Queue string

	// ConsumerID uniquely identifies the consumer on the channel.
	// If empty, the library will generate a unique identity.
	ConsumerID string

	// AutoAck, when set to true, automatically acknowledges deliveries.
	// This should be set to false if the consumer intends to manually acknowledge messages.
	AutoAck bool

	// Exclusive, when set to true, ensures that this is the sole consumer from the queue.
	// If set to false, the server will distribute deliveries across multiple consumers.
	Exclusive bool

	// NoLocal is not supported by RabbitMQ and should generally be set to false.
	NoLocal bool

	// NoWait, when set to true, the consumer will not wait for server confirmations and will
	// immediately start receiving messages. If it is not possible to consume, a channel exception
	// will occur and the channel will be closed.
	NoWait bool

	// DatadogIntegrationEnabled enables automatic extraction of TraceId and SpanId from AMQP message headers
	// for Datadog tracing integration. When set to true, this feature:
	//  1. Attempts to extract tracing information (TraceId and SpanId) from the message headers found in
	//     amqp091.Delivery using the Datadog's TextMapCarrier format.
	//  2. Creates a new Datadog span for the message consumption process. This new span is either:
	//     a. A child span, if TraceId and SpanId are successfully extracted from the message headers, thus
	//     continuing the trace that the message is a part of.
	//     b. A root span, if the message does not contain TraceId and SpanId, indicating the start of a new trace.
	//  3. Attaches the created span to a context.Context, which can then be passed down to other operations or
	//     functions that need tracing. This allows for an end-to-end view of the message processing, including
	//     any subsequent operations triggered by the message consumption.
	DatadogIntegrationEnabled bool

	// Args are optional arguments that can be provided with specific semantics for the queue or server.
	Args amqp091.Table
}

// Consumer initializes and returns a new instance of ConsumerOptions.
func Consumer() *ConsumerOptions {
	return &ConsumerOptions{}
}

// SetURI sets the URI of the AMQP server.
func (c *ConsumerOptions) SetURI(uri string) *ConsumerOptions {
	c.URI = uri
	return c
}

// SetQueue sets the name of the queue to consume from.
func (c *ConsumerOptions) SetQueue(queue string) *ConsumerOptions {
	c.Queue = queue
	return c
}

// SetConsumerID sets a unique identifier for the consumer.
func (c *ConsumerOptions) SetConsumerID(consumerID string) *ConsumerOptions {
	c.ConsumerID = consumerID
	return c
}

// SetAutoAck controls whether messages should be automatically acknowledged.
func (c *ConsumerOptions) SetAutoAck(autoAck bool) *ConsumerOptions {
	c.AutoAck = autoAck
	return c
}

// SetExclusive controls whether this consumer is the sole consumer from the queue.
// When true, ensures exclusive access to the queue. When false, allows multiple consumers.
func (c *ConsumerOptions) SetExclusive(exclusive bool) *ConsumerOptions {
	c.Exclusive = exclusive
	return c
}

// SetNoLocal sets the NoLocal flag. This is not supported by RabbitMQ and is typically set to false.
func (c *ConsumerOptions) SetNoLocal(noLocal bool) *ConsumerOptions {
	c.NoLocal = noLocal
	return c
}

// SetNoWait controls the waiting behavior on server confirmations for consuming messages.
func (c *ConsumerOptions) SetNoWait(noWait bool) *ConsumerOptions {
	c.NoWait = noWait
	return c
}

// SetArgs sets optional arguments with specific semantics for the queue or server.
func (c *ConsumerOptions) SetArgs(args amqp091.Table) *ConsumerOptions {
	c.Args = args
	return c
}

func (c *ConsumerOptions) EnableDatadogIntegration(enable bool) *ConsumerOptions {
	c.DatadogIntegrationEnabled = enable
	return c
}

// MergeConsumerOptions combines multiple ConsumerOptions instances into one.
// For string fields, it uses the first non-empty value.
// For boolean fields, it returns true if any option is true.
// For amqp091.Table, it merges all key-value pairs.
func MergeConsumerOptions(opts ...*ConsumerOptions) *ConsumerOptions {
	merged := &ConsumerOptions{}

	for _, opt := range opts {
		if opt == nil {
			continue
		}

		if merged.URI == "" {
			merged.URI = opt.URI
		}
		if merged.Queue == "" {
			merged.Queue = opt.Queue
		}
		if merged.ConsumerID == "" {
			merged.ConsumerID = opt.ConsumerID
		}
		merged.AutoAck = merged.AutoAck || opt.AutoAck
		merged.Exclusive = merged.Exclusive || opt.Exclusive
		merged.NoLocal = merged.NoLocal || opt.NoLocal
		merged.NoWait = merged.NoWait || opt.NoWait
		merged.DatadogIntegrationEnabled = merged.DatadogIntegrationEnabled || opt.DatadogIntegrationEnabled

		if opt.Args == nil {
			continue
		}

		for key, value := range opt.Args {
			if merged.Args == nil {
				merged.Args = amqp091.Table{}
			}

			merged.Args[key] = value
		}
	}

	return merged
}
