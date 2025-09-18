package ioc

import (
	getbalance "github.com/andre/code-styles-golang/internal/balance/features/get-balance"
	savesnapshot "github.com/andre/code-styles-golang/internal/balance/features/save-snapshot"
	"github.com/andre/code-styles-golang/pkg/datadog/env"
	"go.uber.org/dig"
)

func Configure(container *dig.Container, ddEnvs *env.DatadogEnvironment) error {
	getbalance.Configure(container)
	savesnapshot.Configure(container)

	return nil
}
