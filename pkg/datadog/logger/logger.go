package logger

import (
	"context"
	"os"

	"github.com/andre/code-styles-golang/pkg/datadog/env"
	"github.com/rs/zerolog"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var logger zerolog.Logger
var envs *env.DatadogEnvironment

func ConfigureLogger(e *env.DatadogEnvironment) {
	logger = zerolog.New(os.Stdout).With().Caller().Timestamp().Logger()
	envs = e
}

func Info(ctx context.Context) *zerolog.Event {
	event := logger.Info().Ctx(ctx)

	if !envs.DATADOG_ENABLED {
		return event
	}

	return eventWithTags(ctx, event)
}

func Warn(ctx context.Context) *zerolog.Event {
	event := logger.Warn().Ctx(ctx)

	if !envs.DATADOG_ENABLED {
		return event
	}

	return eventWithTags(ctx, event)
}

func Err(ctx context.Context, err error) *zerolog.Event {
	event := logger.Err(err).Ctx(ctx)

	if !envs.DATADOG_ENABLED {
		return event
	}

	return eventWithTags(ctx, event)
}

func Error(ctx context.Context) *zerolog.Event {
	event := logger.Error().Ctx(ctx)

	if !envs.DATADOG_ENABLED {
		return event
	}

	return eventWithTags(ctx, event)
}

func Debug(ctx context.Context) *zerolog.Event {
	event := logger.Debug().Ctx(ctx)

	if !envs.DATADOG_ENABLED {
		return event
	}

	return eventWithTags(ctx, event)
}

func Fatal(ctx context.Context) *zerolog.Event {
	event := logger.Fatal().Ctx(ctx)

	if !envs.DATADOG_ENABLED {
		return event
	}

	return eventWithTags(ctx, event)
}

func eventWithTags(ctx context.Context, e *zerolog.Event) *zerolog.Event {
	e = e.
		Str("dd.service", envs.DD_SERVICE).
		Str("dd.env", envs.DD_ENV).
		Str("dd.version", envs.DD_VERSION)

	span, found := tracer.SpanFromContext(ctx)

	if !found {
		return e
	}

	return e.
		Uint64("dd.trace_id", span.Context().TraceID()).
		Uint64("dd.span_id", span.Context().SpanID())
}
