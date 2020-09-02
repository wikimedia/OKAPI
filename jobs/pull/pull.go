package pull

import (
	"okapi/lib/task"
)

// Name task name for trigger
var Name task.Name = "pull"

// Task for getting html for pages
func Task(ctx *task.Context) (task.Pool, task.Worker, task.Finish, error) {
	options := Options{}
	err := options.Init(ctx)
	worker := func(id int, payload task.Payload) (string, map[string]interface{}, error) {
		message, info, err := Worker(id, payload)
		options.Position++
		ctx.State.Set("offset", options.Position)
		return message, info, err
	}

	return Pool(ctx, &options), worker, nil, err
}
