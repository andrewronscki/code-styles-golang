package notifysnapshotcreated

import (
	balance "github.com/andre/code-styles-golang/internal/balance/domain"
	e "github.com/andre/code-styles-golang/internal/shared/events"
	"github.com/andre/code-styles-golang/pkg/config"
	cqrsdig "github.com/andre/code-styles-golang/pkg/cqrs-dig"
	"github.com/andre/code-styles-golang/pkg/datadog/env"
	m "github.com/andre/code-styles-golang/pkg/messaging"
	"github.com/andre/code-styles-golang/pkg/messaging/options"
	"go.uber.org/dig"
	"go.uber.org/multierr"
)

func Configure(container *dig.Container, ddEnvs *env.DatadogEnvironment) error {
	return multierr.Combine(
		container.Provide(func() (m.Producer[e.SnapshotCreatedIntegrationEvent], error) {
			opts := options.Producer().
				SetURI(config.Env.GetString("RABBITMQ_URI")).
				SetDestination("balance-exchange").
				SetDestinationKind("exchange").
				SetRoutingKey("snapshot-created").
				EnableDatadogIntegration(ddEnvs.DATADOG_ENABLED)

			return m.CreateProducer[e.SnapshotCreatedIntegrationEvent](opts)
		}),
		cqrsdig.ProvideCommandHandler[*Command, interface{}](
			container,
			NewCommandHandler,
		),
		cqrsdig.ProvideEventSubscriber[*balance.SnapshotCreatedDomainEvent](
			container,
			NewDomainEventHandler,
		),
	)
}
