package rabbitmq

import (
	"github.com/andre/code-styles-golang/pkg/config"
	"github.com/andre/code-styles-golang/pkg/messaging"
)

func ConfigureRabbitMQ() error {
	uri := config.Env.GetString("RABBITMQ_URI")

	cfg := messaging.RabbitMQConfig{
		URI: uri,
		Exchanges: []messaging.Exchange{
			{
				Name:    "balance-exchange",
				Kind:    "topic",
				Durable: true,
			},
		},
		Queues: []messaging.Queue{
			{
				Name:    "save-snapshot",
				Durable: true,
			},
			{
				Name:    "balance.snapshot-created",
				Durable: true,
			},
		},
		QueueBinds: []messaging.QueueBind{
			{
				Exchange: "balance-exchange",
				Name:     "save-snapshot",
				Key:      "snapshot-saved",
			},
			{
				Exchange: "balance-exchange",
				Name:     "balance.snapshot-created",
				Key:      "snapshot-created",
			},
		},
	}

	return cfg.Apply()
}
