package savesnapshot

import (
	"github.com/andre/code-styles-golang/pkg/config"
	cqrsdig "github.com/andre/code-styles-golang/pkg/cqrs-dig"
	"go.uber.org/dig"
)

func Configure(container *dig.Container) {
	container.Provide(func() *BalanceRepositoryParams {
		database := config.Env.GetString("MONGO_DATABASE")
		collection := config.Env.GetString("MONGO_COLLECTION")

		return &BalanceRepositoryParams{
			Database:   database,
			Collection: collection,
		}
	})

	container.Provide(NewBalanceRepository)

	cqrsdig.ProvideCommandHandler[*Command, any](
		container,
		NewCommandHandler,
	)
}
