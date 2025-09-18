package mongodb

import (
	"context"

	"github.com/andre/code-styles-golang/pkg/datadog/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	mongotrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go.mongodb.org/mongo-driver/mongo"
)

var client *mongo.Client

func ConfigureClient(uri string, e *env.DatadogEnvironment) error {
	opts := options.Client()
	opts.ApplyURI(uri)

	if e.DATADOG_ENABLED {
		opts.Monitor = mongotrace.NewMonitor(mongotrace.WithServiceName(e.DD_SERVICE))
	}

	var err error
	client, err = mongo.Connect(context.Background(), opts)

	return err
}

func Client() *mongo.Client {
	return client
}
