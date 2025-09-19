package notifysnapshotcreated

import (
	"context"

	balance "github.com/andre/code-styles-golang/internal/balance/domain"

	"github.com/andre/code-styles-golang/pkg/cqrs"
)

type DomainEventHandler struct{}

func (h *DomainEventHandler) Handle(ctx context.Context, event *balance.SnapshotCreatedDomainEvent) error {
	_, err := cqrs.Send[*Command, any](ctx, &Command{
		UserID:       event.UserID,
		SnapshotDate: event.SnapshotDate,
	})

	return err
}

func NewDomainEventHandler() cqrs.IEventHandler[*balance.SnapshotCreatedDomainEvent] {
	return &DomainEventHandler{}
}
