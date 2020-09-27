package elastic

import (
	"context"
	"fmt"
	"okapi/helpers/logger"

	"github.com/elastic/go-elasticsearch/v7/esapi"
)

// Delete doc from the index
func Delete(id string, name string) {
	wg.Add(1)
	deletes <- &esapi.DeleteRequest{
		DocumentID: id,
		Index:      name,
	}
}

func deleter() {
	for req := range deletes {
		err := delete(req)
		wg.Done()

		if err != nil {
			logger.Search.Error("failed to delete index", err.Error())
		}
	}
}

func delete(req *esapi.DeleteRequest) error {
	res, err := req.Do(context.Background(), client)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("delete failed status code '%d'", res.StatusCode)
	}

	return nil
}
