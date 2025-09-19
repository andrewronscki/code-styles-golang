package notifysnapshotcreated

import (
	"context"

	balance "github.com/andre/code-styles-golang/internal/balance/domain"
	"github.com/andre/code-styles-golang/internal/shared/events"
	"github.com/andre/code-styles-golang/pkg/cqrs"
	"github.com/andre/code-styles-golang/pkg/messaging"
)

type EventProducer struct {
	producer messaging.Producer[*events.SnapshotCreatedIntegrationEvent]
}

func (h *EventProducer) Handle(ctx context.Context, event *balance.SnapshotCreatedDomainEvent) error {
	msg := messaging.Message[*events.SnapshotCreatedIntegrationEvent]{
		Content: &events.SnapshotCreatedIntegrationEvent{
			UserID:       event.UserID,
			SnapshotDate: event.SnapshotDate,
		},
		ContentType: "application/json",
	}

	h.producer.Produce(ctx, msg)

	return nil
}

func NewEventProducer(producer messaging.Producer[*events.SnapshotCreatedIntegrationEvent]) cqrs.IEventHandler[*balance.SnapshotCreatedDomainEvent] {
	return &EventProducer{
		producer: producer,
	}
}
