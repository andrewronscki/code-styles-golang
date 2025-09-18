package behaviors

import (
	"context"
	"reflect"

	"github.com/andre/code-styles-golang/pkg/cqrs"
	"github.com/andre/code-styles-golang/pkg/datadog/logger"
)

type LoggingBehavior struct{}

func (b *LoggingBehavior) Handle(ctx context.Context, request any, next cqrs.NextFunc) (any, error) {
	logger.Info(ctx).Msgf("handling request of type %v", reflect.TypeOf(request))

	res, err := next()

	if err != nil {
		logger.Err(ctx, err).Msgf("failed to handle request of type %v", reflect.TypeOf(request))
		return res, err
	}

	logger.Info(ctx).Msgf("request of type %v handled successfully", reflect.TypeOf(request))

	return res, err
}

func NewLoggingBehavior() *LoggingBehavior {
	return &LoggingBehavior{}
}
