package savesnapshot

import (
	"context"
	"sync"

	"github.com/andre/code-styles-golang/internal/shared/events"
	"github.com/andre/code-styles-golang/pkg/cqrs"
	"github.com/andre/code-styles-golang/pkg/datadog/logger"
	"github.com/andre/code-styles-golang/pkg/hosting"
	"github.com/andre/code-styles-golang/pkg/messaging"
	"github.com/andre/code-styles-golang/pkg/messaging/options"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type Worker struct {
	Consumer       messaging.Consumer[events.UserProcessIntegrationEvent]
	DatadogEnabled bool
}

func (w *Worker) Run(ctx context.Context, exit func()) {
	defer exit()

	wg := &sync.WaitGroup{}

	msgs, err := w.Consumer.Consume(ctx)

	if err != nil {
		logger.Err(ctx, err).Msg("save-snapshot worker exiting, consume messages failed")
		return
	}

	for {
		select {
		case <-ctx.Done():
			wg.Wait()
			return

		case msg, ok := <-msgs:

			if !ok {
				continue
			}

			wg.Add(1)
			go func() {
				defer func() {
					if err := recover(); err != nil {
						logger.Err(ctx, err.(error)).Msg("worker panic recovered, command cannot be handled")
						wg.Done()
					}
				}()

				var err error

				if w.DatadogEnabled {
					span, spanCtx := tracer.StartSpanFromContext(msg.Context, "worker.process")
					span.SetTag("worker.name", "save-snapshot")
					defer span.Finish()

					_, err = cqrs.Send[*Command, any](spanCtx, &Command{
						UserID: msg.Content.UserID,
					})
				} else {
					_, err = cqrs.Send[*Command, any](msg.Context, &Command{
						UserID: msg.Content.UserID,
					})
				}

				if err != nil {
					logger.Err(ctx, err)
				}

				wg.Done()
			}()
		}
	}
}

func CreateWorker(
	uri, consumerID, consumerQueue string, ddEnabled bool) hosting.Worker {
	opts := options.Consumer().
		SetURI(uri).
		SetConsumerID(consumerID).
		SetQueue(consumerQueue).
		SetAutoAck(true).
		EnableDatadogIntegration(ddEnabled)

	consumer, err := messaging.CreateConsumer[events.UserProcessIntegrationEvent](opts)

	if err != nil {
		logger.Fatal(context.Background()).Err(err)
	}

	return &Worker{
		Consumer:       consumer,
		DatadogEnabled: ddEnabled,
	}
}
