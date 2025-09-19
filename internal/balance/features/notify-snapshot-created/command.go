package notifysnapshotcreated

import (
	"context"
	"time"

	"github.com/andre/code-styles-golang/internal/shared/events"
	"github.com/andre/code-styles-golang/pkg/cqrs"
	"github.com/andre/code-styles-golang/pkg/datadog/logger"
	m "github.com/andre/code-styles-golang/pkg/messaging"
)

type Command struct {
	UserID       int64
	SnapshotDate time.Time
}

type CommandHandler struct {
	store    events.EventStore
	producer m.Producer[events.SnapshotCreatedIntegrationEvent]
}

func (h *CommandHandler) Handle(ctx context.Context, command *Command) (any, error) {
	logger.Info(ctx).Msgf("O snapshot para userID %d do dia %s foi gerado com sucesso", command.UserID, command.SnapshotDate)

	msg := m.Message[events.SnapshotCreatedIntegrationEvent]{
		Content: events.SnapshotCreatedIntegrationEvent{
			UserID:       command.UserID,
			SnapshotDate: command.SnapshotDate,
		},
		ContentType: "application/json",
	}

	err := h.producer.Produce(ctx, msg)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func NewCommandHandler(store events.EventStore, producer m.Producer[events.SnapshotCreatedIntegrationEvent]) cqrs.ICommandHandler[*Command, any] {
	return &CommandHandler{
		store:    store,
		producer: producer,
	}
}
