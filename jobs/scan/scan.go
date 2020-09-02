package scan

import (
	"okapi/lib/dump"
	"okapi/lib/task"
)

// Name task name for trigger
var Name task.Name = "scan"

// Task to get titles from dump and store them in database
func Task(ctx *task.Context) (task.Pool, task.Worker, task.Finish, error) {
	options := Options{}
	titles := []string{}
	worker := func(id int, payload task.Payload) (string, map[string]interface{}, error) {
		message, info, err := Worker(id, payload)
		options.Position++
		ctx.State.Set("offset", options.Position)
		return message, info, err
	}

	err := options.Init(ctx)

	if err == nil {
		titles, err = dump.Titles(ctx.Project.DBName, options.Folder)
	}

	return Pool(ctx, &options, titles), worker, nil, err
}
