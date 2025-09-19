package getbalance

import (
	"github.com/andre/code-styles-golang/pkg/config"
	cqrsdig "github.com/andre/code-styles-golang/pkg/cqrs-dig"
	"go.uber.org/dig"
	"go.uber.org/multierr"
)

func Configure(container *dig.Container) error {
	return multierr.Combine(
		container.Provide(func() *BalanceRepositoryParams {
			database := config.Env.GetString("MONGO_DATABASE")
			collection := config.Env.GetString("MONGO_COLLECTION")

			return &BalanceRepositoryParams{
				Database:   database,
				Collection: collection,
			}
		}),
		container.Provide(NewBalanceRepository),
		cqrsdig.ProvideQueryHandler[*Query, *Model](
			container,
			NewQueryHandler,
		))
}
