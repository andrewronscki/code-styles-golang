package ioc

import (
	getbalance "github.com/andre/code-styles-golang/internal/balance/features/get-balance"
	notifysnapshotcreated "github.com/andre/code-styles-golang/internal/balance/features/notify-snapshot-created"
	savesnapshot "github.com/andre/code-styles-golang/internal/balance/features/save-snapshot"
	"github.com/andre/code-styles-golang/pkg/datadog/env"
	"go.uber.org/dig"
	"go.uber.org/multierr"
)

func Configure(container *dig.Container, ddEnvs *env.DatadogEnvironment) error {
	multierr.Combine(
		getbalance.Configure(container),
		savesnapshot.Configure(container),
		notifysnapshotcreated.Configure(container, ddEnvs),
	)

	return nil
}
