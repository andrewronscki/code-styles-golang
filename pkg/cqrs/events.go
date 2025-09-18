package cqrs

import (
    "context"
    "errors"
    "fmt"
    "reflect"

    "go.uber.org/multierr"
)

type IEventHandler[TEvent any] interface {
    Handle(ctx context.Context, event TEvent) error
}

var eventHandlers map[reflect.Type][]any

func init() {
    eventHandlers = make(map[reflect.Type][]any)
}

func RegisterEventSubscriber[TEvent any](handler IEventHandler[TEvent]) error {
    var event TEvent
    eventType := reflect.TypeOf(event)
    handlers, found := eventHandlers[eventType]

    if !found {
        eventHandlers[eventType] = []any{
            handler,
        }
        return nil
    }

    eventHandlers[eventType] = append(handlers, handler)

    return nil
}

func RegisterEventSubscribers[TEvent any](handlers ...IEventHandler[TEvent]) error {
    if len(handlers) <= 0 {
        return errors.New("at least one handler must be provided")
    }

    for _, handler := range handlers {
        RegisterEventSubscriber(handler)
    }

    return nil
}

func PublishEvent[TEvent any](ctx context.Context, event TEvent) error {
    eventType := reflect.TypeOf(event)
    handlers, found := eventHandlers[eventType]

    if !found {
        msg := fmt.Sprintf("no event handler found event of type: %T", event)
        return errors.New(msg)
    }

    var err error = nil

    for _, h := range handlers {
        handler, ok := h.(IEventHandler[TEvent])

        if ok {
            handleErr := handler.Handle(ctx, event)

            if handleErr != nil {
                err = multierr.Append(err, handleErr)
            }

            continue
        }

        args := []reflect.Value{
            reflect.ValueOf(ctx),
            reflect.ValueOf(event),
        }

        r := reflect.ValueOf(h).MethodByName("Handle").Call(args)

        handleErr := r[0].Interface()

        if handleErr != nil {
            err = multierr.Append(err, handleErr.(error))
        }
    }

    return err
}