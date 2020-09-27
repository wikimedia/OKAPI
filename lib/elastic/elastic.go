package elastic

import (
	"okapi/lib/env"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

var client *elasticsearch.Client
var indexes = make(chan *esapi.IndexRequest)
var deletes = make(chan *esapi.DeleteRequest)

// Init function to initialize elastic client
func Init() (err error) {
	client, err = elasticsearch.NewClient(elasticsearch.Config{
		Username: env.Context.ElasticUsername,
		Password: env.Context.ElasticPassword,
		Addresses: []string{
			env.Context.ElasticURL,
		},
	})

	if err == nil {
		go indexer()
		go deleter()
	}

	return err
}

// Client getter for elastic seacrch client
func Client() *elasticsearch.Client {
	return client
}

// Close function to close index channel
func Close() {
	close(indexes)
	wg.Wait()
}
