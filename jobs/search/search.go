package search

import (
	"bytes"
	"encoding/json"
	"okapi/lib/storage"
	"okapi/lib/task"
)

// Name task name for trigger
var Name task.Name = "search"

const lang string = "en"

// Task for getting html for pages
func Task(ctx *task.Context) (task.Pool, task.Worker, task.Finish, error) {
	fields := map[string]interface{}{}
	getters := map[string]func(lang string) (interface{}, error){
		"ns_id":     getNamespaces,
		"site_code": getSiteCodes,
		"lang":      getLangs,
		"site_name": getSiteNames,
	}

	for fieldName, getValues := range getters {
		options, err := getValues(lang)

		if err == nil {
			fields[fieldName] = options
		} else {
			return nil, nil, nil, err
		}
	}

	content, err := json.Marshal(fields)
	if err != nil {
		return nil, nil, nil, err
	}

	err = storage.Local.Client().Put("options/"+lang+"/options.json", bytes.NewReader(content))
	return nil, nil, nil, err
}
