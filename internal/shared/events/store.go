package events

import (
	"sync"
)

type EventStore interface {
	AddEvent(event any)
	AddEvents(events []any)
	GetEventsAndClear() []any
}

type InMemoryEventStore struct {
	events []any
	mu     sync.Mutex
}

func (s *InMemoryEventStore) AddEvent(event any) {
	s.mu.Lock()
	s.events = append(s.events, event)
	s.mu.Unlock()
}

func (s *InMemoryEventStore) AddEvents(events []any) {
	s.mu.Lock()
	s.events = append(s.events, events...)
	s.mu.Unlock()
}

func (s *InMemoryEventStore) GetEventsAndClear() []any {
	s.mu.Lock()
	defer s.mu.Unlock()

	events := make([]any, 0, len(s.events))
	for _, e := range s.events {
		events = append(events, e)
	}
	s.events = []any{}

	return events
}

func NewInMemoryEventStore() EventStore {
	return &InMemoryEventStore{
		events: []any{},
	}
}
