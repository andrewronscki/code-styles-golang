package savesnapshot

import (
	"context"
	"time"

	balance "github.com/andre/code-styles-golang/internal/balance/domain"
	"github.com/andre/code-styles-golang/internal/shared/events"
	"github.com/andre/code-styles-golang/pkg/cqrs"
)

type Command struct {
	UserID int64
}

type CommandHandler struct {
	repository Repository
	store      events.EventStore
}

func (h *CommandHandler) Handle(ctx context.Context, command *Command) (any, error) {
	snapshotDate := time.Now().UTC()

	balance := balance.NewBalance(
		command.UserID,
		snapshotDate,
	)

	balance.SetBalance(10.0)

	id, err := h.repository.Insert(ctx, balance)

	if err != nil {
		return nil, err
	}

	balance.SetID(id)

	h.store.AddEvent(balance.RaiseSnapshotCreatedDomainEvent())

	return nil, nil
}

func NewCommandHandler(repository Repository, store events.EventStore) cqrs.ICommandHandler[*Command, any] {
	return &CommandHandler{
		repository: repository,
		store:      store,
	}
}
