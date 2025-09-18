package savesnapshot

import (
	"context"
	"time"

	balance "github.com/andre/code-styles-golang/internal/balance/domain"
	"github.com/andre/code-styles-golang/pkg/cqrs"
	"github.com/andre/code-styles-golang/pkg/datadog/logger"
)

type Command struct {
	UserID int64
}

type CommandHandler struct {
	repository Repository
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

	logger.Info(ctx).Msgf("[%d] created a new snapshot with id: %s", balance.UserID, balance.ID)

	return balance, nil
}

func NewCommandHandler(repository Repository) cqrs.ICommandHandler[*Command, any] {
	return &CommandHandler{
		repository: repository,
	}
}
