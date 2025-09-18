package events

import "go.uber.org/dig"

func Configure(container *dig.Container) error {
	return container.Provide(NewInMemoryEventStore)
}
