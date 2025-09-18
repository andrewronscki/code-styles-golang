package ioc

import (
	"github.com/andre/code-styles-golang/internal/app/behaviors"
	balanceioc "github.com/andre/code-styles-golang/internal/balance/ioc"
	"github.com/andre/code-styles-golang/internal/shared/events"
	"github.com/andre/code-styles-golang/pkg/datadog/env"
	"go.uber.org/dig"
	"go.uber.org/multierr"
)

func Configure(ddEnvs *env.DatadogEnvironment) (*dig.Container, error) {
	container := dig.New()

	return container, multierr.Combine(
		events.Configure(container),
		behaviors.Configure(container),
		balanceioc.Configure(container, ddEnvs),
	)
}
