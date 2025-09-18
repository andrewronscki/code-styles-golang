package cqrs

type INotifiable interface {
    AddEvent(event any)
    ClearEvents()
    GetEvents() []any
}