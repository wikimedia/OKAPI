package index

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	page_index "okapi/indexes/page"
	"okapi/lib/elastic"
	"okapi/lib/task"
	"okapi/models"
	"runtime"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v7/esutil"
	"gopkg.in/gookit/color.v1"
)

// Worker processing one page from the queue, getting html into s3
func Worker(id int, payload task.Payload) (string, map[string]interface{}, error) {
	start := time.Now().UTC()
	info := map[string]interface{}{}
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:      page_index.Name,
		Client:     elastic.Client(),
		NumWorkers: runtime.NumCPU(),
	})

	if err != nil {
		return "", info, err
	}

	for _, page := range payload.([]*models.Page) {
		body, err := json.Marshal(page.Index())

		if err != nil {
			color.Error.Println(err)
			continue
		}

		err = bi.Add(context.Background(), getIndex(page, body))

		if err != nil {
			return "", info, err
		}
	}

	if err := bi.Close(context.Background()); err != nil {
		return "", info, err
	}

	stats := bi.Stats()
	dur := time.Since(start)
	info["_num_failed"] = stats.NumFailed
	info["_num_flushed"] = stats.NumFlushed

	message := fmt.Sprintf(
		"indexed [%d] documents with [%d] errors in %s (%d docs/sec)",
		stats.NumFlushed,
		stats.NumFailed,
		dur.Truncate(time.Millisecond),
		int64(1000.0/float64(dur/time.Millisecond)*float64(stats.NumFlushed)),
	)

	return message, info, nil
}

func getIndex(page *models.Page, body []byte) esutil.BulkIndexerItem {
	return esutil.BulkIndexerItem{
		Action:     "index",
		DocumentID: strconv.Itoa(page.ID),
		Body:       bytes.NewReader(body),
		OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
			if err != nil {
				color.Error.Println(err)
			} else {
				color.Error.Println(res.Error.Type, res.Error.Reason)
			}
		},
	}
}
