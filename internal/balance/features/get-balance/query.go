package getbalance

import (
	"context"
	"time"

	"github.com/andre/code-styles-golang/pkg/cqrs"
)

type Query struct {
	UserID       int64     `json:"user_id"`
	SnapshotDate time.Time `json:"snapshot_date"`
}

type QueryHandler struct {
	repository Repository
}

func (h *QueryHandler) Handle(ctx context.Context, query *Query) (*Model, error) {
	findUser, err := h.repository.FindBalance(ctx, (*Filter)(query))

	if err != nil {
		return nil, err
	}

	model := &Model{}
	findUser.Marshal(model)

	return model, nil
}

func NewQueryHandler(repository Repository) cqrs.IQueryHandler[*Query, *Model] {
	return &QueryHandler{
		repository: repository,
	}
}
