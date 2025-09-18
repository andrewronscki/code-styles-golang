package messaging

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func connect(uri string) (*amqp091.Connection, error) {
	var err error
	var connection *amqp091.Connection

	log.Printf("opening amqp connection")

	err = retry(func() error {
		connection, err = amqp091.Dial(uri)
		return err
	},
		"amqp dial",
		5,
	)
	return connection, err
}

func retry(f func() error, operation string, attempts int) error {
	if attempts < 1 {
		return fmt.Errorf("attempts must be greater than or equal to 1")
	}

	if f == nil {
		return fmt.Errorf("target function must not be nil")
	}

	var err error

	for attempt := 1; attempt < attempts; attempt++ {
		log.Printf("%s attempt %d of %d", operation, attempt, attempts)

		err = f()

		if err == nil {
			return nil
		}

		sleep := time.Duration(math.Pow(2, float64(attempt))) * time.Second

		log.Printf("%s failed, retrying in %v", operation, sleep.Seconds())

		time.Sleep(sleep)
	}

	return err
}
