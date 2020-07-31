package bundle

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"okapi/helpers/writer"
	"okapi/lib/task"
	"okapi/models"
)

// Worker adding the html to write queue
func Worker(ctx *task.Context, options *Options) func(id int, payload task.Payload) (string, map[string]interface{}, error) {
	return func(id int, payload task.Payload) (string, map[string]interface{}, error) {
		message := "page title: '%s', page id: #%d"
		page := payload.(*models.Page)
		info := map[string]interface{}{
			"_title": page.Title,
			"_id":    page.ID,
		}

		body, err := options.Storage.Get(page.Path)
		if err != nil {
			return "", info, fmt.Errorf(message+", %s", page.Title, page.ID, err)
		}

		html, err := ioutil.ReadAll(body)
		if err != nil {
			return "", info, fmt.Errorf(message+", %s", page.Title, page.ID, err)
		}

		options.Writer.Add(writer.Payload{
			ReadCloser: ioutil.NopCloser(bytes.NewBuffer(html)),
			Name:       page.Title + ".html",
			Size:       int64(len(html)),
		})

		return fmt.Sprintf(message, page.Title, page.ID), info, nil
	}
}
