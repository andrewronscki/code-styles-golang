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
		},
		QueueBinds: []messaging.QueueBind{
			{
				Exchange: "balance-exchange",
				Name:     "save-snapshot",
				Key:      "snapshot-saved",
			},
		},
	}

	return cfg.Apply()
}
