package behaviors

import (
	"context"
	"sync"

	// "sync"

	"github.com/andre/code-styles-golang/internal/shared/events"
	"github.com/andre/code-styles-golang/pkg/cqrs"
	"github.com/andre/code-styles-golang/pkg/datadog/logger"
)

type EventDispatcherBehavior struct {
	store events.EventStore
}

func (b *EventDispatcherBehavior) Handle(ctx context.Context, request any, next cqrs.NextFunc) (any, error) {
	res, err := next()

	if err != nil {
		return res, err
	}

	events := b.store.GetEventsAndClear()

	len := len(events)

	if len <= 0 {
		return res, err
	}

	logger.Info(ctx).Msgf("dispatching %d domain event(s)", len)
	wg := &sync.WaitGroup{}
	for _, event := range events {
		event := event
		wg.Add(1)
		go func(e any) {
			defer wg.Done()
			cqrs.PublishEvent(ctx, e)
		}(event)
	}

	wg.Wait()

	return res, err
}

func NewEventDispatcherBehavior(store events.EventStore) *EventDispatcherBehavior {
	return &EventDispatcherBehavior{
		store: store,
	}
}
