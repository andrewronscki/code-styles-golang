package cqrs

import (
    "context"
    "errors"
    "fmt"
    "sort"
)

type NextFunc func() (any, error)

type IBehavior interface {
    Handle(ctx context.Context, request any, next NextFunc) (any, error)
}

var commandBehaviors map[int]any
var queryBehaviors map[int]any

func init() {
    commandBehaviors = make(map[int]any)
    queryBehaviors = make(map[int]any)
}

func RegisterCommandBehavior(order int, behavior IBehavior) error {
    _, found := commandBehaviors[order]

    if found {
        msg := fmt.Sprintf("position %d is taken by another command behavior.", order)
        return errors.New(msg)
    }

    commandBehaviors[order] = behavior

    return nil

}

func RegisterQueryBehavior(order int, behavior IBehavior) error {
    _, found := queryBehaviors[order]

    if found {
        msg := fmt.Sprintf("position %d is taken by another query behavior.", order)
        return errors.New(msg)
    }

    queryBehaviors[order] = behavior

    return nil
}

func sortBehaviors(behaviors map[int]any) []any {
    keys := make([]int, 0)

    for key := range behaviors {
        keys = append(keys, key)
    }

    sort.Sort(sort.Reverse(sort.IntSlice(keys)))

    sorted := make([]any, 0)

    for _, key := range keys {
        sorted = append(sorted, behaviors[key])
    }

    return sorted
}