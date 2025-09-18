package work

import (
	savesnapshot "github.com/andre/code-styles-golang/internal/balance/features/save-snapshot"
	"github.com/andre/code-styles-golang/pkg/config"
	"github.com/andre/code-styles-golang/pkg/datadog/env"
	"github.com/andre/code-styles-golang/pkg/hosting"
)

func Workers(ddEnvs *env.DatadogEnvironment) []hosting.Worker {
	uri := config.Env.GetString("RABBITMQ_URI")
	return []hosting.Worker{
		savesnapshot.CreateWorker(uri,
			"save-snapshot-consumer",
			"save-snapshot",
			ddEnvs.DATADOG_ENABLED,
		),
	}
}
