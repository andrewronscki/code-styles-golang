package datadog

import (
	"context"
	"fmt"

	"github.com/andre/code-styles-golang/pkg/datadog/env"
	"github.com/andre/code-styles-golang/pkg/datadog/logger"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"gopkg.in/DataDog/dd-trace-go.v1/profiler"
)

func Start(e *env.DatadogEnvironment) {
	if !e.DATADOG_ENABLED {
		return
	}

	agent := agentAddr(e)
	statsd := statsdAddr(e)

	if agent == "" || statsd == "" {
		logger.Warn(context.Background()).Msg("agent and/or statsd address must not be empty when datadog integration is enabled; tracer will not be started")
		return
	}

	tracer.Start(
		tracer.WithService(e.DD_SERVICE),
		tracer.WithEnv(e.DD_ENV),
		tracer.WithServiceVersion(e.DD_VERSION),
		tracer.WithAgentAddr(agent),
		tracer.WithDogstatsdAddress(statsd),
		tracer.WithSampler(tracer.NewAllSampler()),
		tracer.WithGlobalServiceName(true),
		tracer.WithTraceEnabled(true),
		tracer.WithRuntimeMetrics(),
	)

	profiler.Start(
		profiler.WithService(e.DD_SERVICE),
		profiler.WithEnv(e.DD_ENV),
		profiler.WithVersion(e.DD_VERSION),
		profiler.WithTags(fmt.Sprintf("version:%s", e.DD_VERSION)),
		profiler.WithAgentAddr(agent),
		profiler.WithProfileTypes(
			profiler.CPUProfile,
			profiler.HeapProfile,
			profiler.GoroutineProfile,
			profiler.MutexProfile,
			profiler.BlockProfile,
			profiler.MetricsProfile,
		),
	)

}

func Stop(e *env.DatadogEnvironment) {
	if !e.DATADOG_ENABLED {
		return
	}

	defer tracer.Stop()
	defer profiler.Stop()
}

func agentAddr(e *env.DatadogEnvironment) string {
	if e.DD_AGENT_HOST != "" {
		return fmt.Sprintf("%s:%d", e.DD_AGENT_HOST, env.DATADOG_AGENT_PORT)
	}

	return e.DATADOG_AGENT_ADDR
}

func statsdAddr(e *env.DatadogEnvironment) string {
	if e.DD_AGENT_HOST != "" {
		return fmt.Sprintf("%s:%d", e.DD_AGENT_HOST, env.DATADOG_DOGSTATSD_PORT)
	}

	return e.DATADOG_DOGSTATSD_ADDR
}
