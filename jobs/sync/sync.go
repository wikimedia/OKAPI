package sync

import (
	"okapi/lib/task"
)

// Task for getting html for pages
func Task(ctx *task.Context) (task.Pool, task.Worker, task.Finish) {
	options := Options{}
	options.Init(ctx)
	worker := func(id int, payload task.Payload) (string, map[string]interface{}, error) {
		message, info, err := Worker(id, payload)
		options.Position++
		ctx.State.Set("offset", options.Position)
		return message, info, err
	}

	return Pool(ctx, &options), worker, nil
}
