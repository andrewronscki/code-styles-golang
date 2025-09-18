package work

import (
	balanceworkers "github.com/andre/code-styles-golang/internal/balance/work"
	"github.com/andre/code-styles-golang/pkg/datadog/env"
	"github.com/andre/code-styles-golang/pkg/hosting"
)

func Workers(ddEnvs *env.DatadogEnvironment) []hosting.Worker {
	return balanceworkers.Workers(ddEnvs)
}
