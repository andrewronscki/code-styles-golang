package options

// ProducerOptions holds configuration options for a message producer.
type ProducerOptions struct {
	// URI specifies the AMQP server endpoint.
	URI string

	// Destination exchange or queue for the message to be sent.
	Destination string

	// DestinationKind must be either queue or exchange
	DestinationKind string

	// Routing key for the message.
	RoutingKey string

	// Mandatory flag for message publishing.
	Mandatory bool

	// Immediate flag for message publishing.
	Immediate bool

	// DatadogIntegrationEnabled enables automatic injection of TraceId and SpanId into AMQP message headers
	// for Datadog tracing integration when publishing messages. When set to true, this feature:
	//  1. Extracts the current TraceId and SpanId from the active Datadog span within the context.Context.
	//     This is done by creating a Datadog's TextMapCarrier and using the tracer.Inject method to populate it
	//     with the span's context.
	//  2. Appends the keys and values from the TextMapCarrier to the message headers before publishing.
	//     This ensures that the tracing information is propagated along with the message, allowing for
	//     distributed tracing across services that consume these messages.
	DatadogIntegrationEnabled bool
}

// Producer creates a new ProducerOptions instance with default values.
func Producer() *ProducerOptions {
	return &ProducerOptions{}
}

// SetURI sets the URI for the AMQP server.
func (p *ProducerOptions) SetURI(uri string) *ProducerOptions {
	p.URI = uri
	return p
}

// SetDestination sets the destination exchange or queue for the message to be sent.
func (p *ProducerOptions) SetDestination(destination string) *ProducerOptions {
	p.Destination = destination
	return p
}

// SetDestinationKind sets the destination kind to be either exchange or queue.
func (p *ProducerOptions) SetDestinationKind(kind string) *ProducerOptions {
	p.DestinationKind = kind
	return p
}

// SetRoutingKey sets the routing key for the message.
// The routing key is used to route messages to the correct queue.
func (p *ProducerOptions) SetRoutingKey(routingKey string) *ProducerOptions {
	p.RoutingKey = routingKey
	return p
}

// SetMandatory sets the mandatory flag for message publishing.
// If set to true, the server will return an undeliverable message with a Return method. If false, the server silently drops the message.
func (p *ProducerOptions) SetMandatory(mandatory bool) *ProducerOptions {
	p.Mandatory = mandatory
	return p
}

// SetImmediate sets the immediate flag for message publishing.
// If set to true, the message will be returned if it cannot be routed to a queue consumer immediately.
func (p *ProducerOptions) SetImmediate(immediate bool) *ProducerOptions {
	p.Immediate = immediate
	return p
}

func (p *ProducerOptions) EnableDatadogIntegration(enable bool) *ProducerOptions {
	p.DatadogIntegrationEnabled = enable
	return p
}

// MergeProducerOptions combines multiple ProducerOptions instances into one.
// For string fields, it uses the first non-empty value.
// For boolean fields, it returns true if any option is true.
func MergeProducerOptions(opts ...*ProducerOptions) *ProducerOptions {
	merged := &ProducerOptions{}

	for _, opt := range opts {
		if opt == nil {
			continue
		}

		if merged.URI == "" {
			merged.URI = opt.URI
		}

		if merged.Destination == "" {
			merged.Destination = opt.Destination
		}

		if merged.DestinationKind == "" {
			merged.DestinationKind = opt.DestinationKind
		}

		if merged.RoutingKey == "" {
			merged.RoutingKey = opt.RoutingKey
		}

		merged.Mandatory = merged.Mandatory || opt.Mandatory
		merged.Immediate = merged.Immediate || opt.Immediate
		merged.DatadogIntegrationEnabled = merged.DatadogIntegrationEnabled || opt.DatadogIntegrationEnabled
	}

	return merged
}
