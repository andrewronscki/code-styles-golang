package main

import (
	"context"

	"github.com/andre/code-styles-golang/internal/app/ioc"
	"github.com/andre/code-styles-golang/internal/app/rabbitmq"
	"github.com/andre/code-styles-golang/internal/app/routes"
	"github.com/andre/code-styles-golang/internal/app/work"
	"github.com/andre/code-styles-golang/pkg/config"
	"github.com/andre/code-styles-golang/pkg/datadog"
	"github.com/andre/code-styles-golang/pkg/datadog/env"
	httpclient "github.com/andre/code-styles-golang/pkg/datadog/http"
	"github.com/andre/code-styles-golang/pkg/datadog/logger"
	"github.com/andre/code-styles-golang/pkg/datadog/mongodb"
	"github.com/andre/code-styles-golang/pkg/hosting"
)

// @title Balance API
// @version 2.0
// @description Code styles Golang API Example
// @termsOfService http://swagger.io/terms/

// @contact.name André
// @contact.url https://www.devtoolshq.dev
// @contact.email contato@devtoolshq.dev

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
// @schemes http
const devEnvPath = "../config/dev.env"

func main() {
	config.LoadEnv("", devEnvPath)

	ddEnvs := &env.DatadogEnvironment{
		DATADOG_ENABLED:        config.Env.GetBool("DATADOG_ENABLED"),
		DATADOG_AGENT_ADDR:     config.Env.GetString("DATADOG_AGENT_ADDR"),
		DATADOG_DOGSTATSD_ADDR: config.Env.GetString("DATADOG_DOGSTATSD_ADDR"),
		DD_SERVICE:             config.Env.GetString("DD_SERVICE"),
		DD_ENV:                 config.Env.GetString("DD_ENV"),
		DD_VERSION:             config.Env.GetString("DD_VERSION"),
		DD_AGENT_HOST:          config.Env.GetString("DD_AGENT_HOST"),
	}

	logger.ConfigureLogger(ddEnvs)

	if err := rabbitmq.ConfigureRabbitMQ(); err != nil {
		logger.Fatal(context.Background()).AnErr("error", err).Send()
	}

	httpclient.ConfigureHttpClient(ddEnvs)

	if err := mongodb.ConfigureClient(config.Env.GetString("MONGO_URI"), ddEnvs); err != nil {
		logger.Fatal(context.Background()).AnErr("error", err).Send()
	}

	if _, err := ioc.Configure(ddEnvs); err != nil {
		logger.Fatal(context.Background()).AnErr("error", err).Send()
	}

	host := &hosting.Host{
		Addr:    config.Env.GetString("ADDR"),
		Router:  routes.Router(ddEnvs),
		Workers: work.Workers(ddEnvs),
	}

	datadog.Start(ddEnvs)
	defer datadog.Stop(ddEnvs)

	host.Start()
}
