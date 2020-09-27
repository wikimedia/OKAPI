package elastic

import (
	"bytes"
	"context"
	"fmt"
	"okapi/helpers/logger"
	"sync"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

// Index this is a index
type Index struct {
	ID   string
	Name string
	Body []byte
}

var wg = sync.WaitGroup{}

// Sync create index for document
func Sync(id string, name string, body []byte) {
	wg.Add(1)
	indexes <- &esapi.IndexRequest{
		Index:      name,
		DocumentID: id,
		Body:       bytes.NewReader(body),
	}
}

func indexer() {
	for req := range indexes {
		err := index(req)
		wg.Done()

		if err != nil {
			logger.Search.Error("failed to update index", err.Error())
		}
	}
}

func index(req *esapi.IndexRequest) error {
	res, err := req.Do(context.Background(), client)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("req failed status code '%d'", res.StatusCode)
	}

	return nil
}
