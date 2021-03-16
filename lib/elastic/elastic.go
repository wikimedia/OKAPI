package elastic

import (
	"errors"
	"okapi-data-service/lib/env"

	"github.com/elastic/go-elasticsearch/v7"
)

// ErrDuplicateClient duplication of elastic client
var ErrDuplicateClient = errors.New("duplicate elastic client")

var client *elasticsearch.Client

// Client get elasticsearch
func Client() *elasticsearch.Client {
	return client
}

// Init initialize new elastic instance
func Init() error {
	var err error

	if client != nil {
		return ErrDuplicateClient
	}

	client, err = elasticsearch.NewClient(elasticsearch.Config{
		Username: env.ElasticUsername,
		Password: env.ElasticPassword,
		Addresses: []string{
			env.ElasticURL,
		},
	})

	return err
}
