package messaging

import "github.com/rabbitmq/amqp091-go"

type Exchange struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp091.Table
}

type ExchangeBind struct {
	Destination string
	Key         string
	Source      string
	NoWait      bool
	Args        amqp091.Table
}

type Queue struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp091.Table
}

type QueueBind struct {
	Name     string
	Key      string
	Exchange string
	NoWait   bool
	Args     amqp091.Table
}

type RabbitMQConfig struct {
	URI           string
	Exchanges     []Exchange
	Queues        []Queue
	ExchangeBinds []ExchangeBind
	QueueBinds    []QueueBind
}

func (cfg *RabbitMQConfig) Apply() error {
	connection, err := connect(cfg.URI)

	if err != nil {
		return err
	}

	channel, err := connection.Channel()

	if err != nil {
		return err
	}

	defer connection.Close()
	defer channel.Close()

	for _, e := range cfg.Exchanges {
		channel.ExchangeDeclare(
			e.Name,
			e.Kind,
			e.Durable,
			e.AutoDelete,
			e.Internal,
			e.NoWait,
			e.Args,
		)
	}

	for _, q := range cfg.Queues {
		channel.QueueDeclare(
			q.Name,
			q.Durable,
			q.AutoDelete,
			q.Exclusive,
			q.NoWait,
			q.Args,
		)
	}

	for _, eb := range cfg.ExchangeBinds {
		channel.ExchangeBind(
			eb.Destination,
			eb.Key,
			eb.Source,
			eb.NoWait,
			eb.Args,
		)
	}

	for _, qb := range cfg.QueueBinds {
		channel.QueueBind(
			qb.Name,
			qb.Key,
			qb.Exchange,
			qb.NoWait,
			qb.Args,
		)
	}

	return nil
}
