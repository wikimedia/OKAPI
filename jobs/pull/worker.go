package pull

import (
	"fmt"
	"sync"

	"okapi/lib/task"
	"okapi/models"
)

const message = "page title: '%s', page id: #%d"

type getter func(page *models.Page) error

type errors struct {
	sync.RWMutex
	items []error
}

// Worker processing one page from the queue, getting html into s3
func Worker(id int, payload task.Payload) (string, map[string]interface{}, error) {
	page := payload.(*models.Page)
	page.WikitextPath = getWikitextPath(page)
	page.Path = getHTMLPath(page)
	info := getInfo(page)
	fields := []string{}
	wg := sync.WaitGroup{}
	errs := new(errors)
	getters := map[string]getter{
		"path": getHTML,
	}

	for field, get := range getters {
		wg.Add(1)
		go processor(page, get, errs, &wg)
		fields = append(fields, field)
	}

	wg.Wait()

	for _, err := range errs.items {
		if err != nil {
			return "", info, err
		}
	}

	err := updatePage(page, fields...)

	if err != nil {
		return "", info, getError(page, err)
	}

	return fmt.Sprintf(message, page.Title, page.ID), info, nil
}
