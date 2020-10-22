package upload

import (
	"okapi/helpers/bundle"
	"okapi/lib/task"
)

// Name task name for trigger
var Name task.Name = "upload"

// Task for bundling the html
func Task(ctx *task.Context) (task.Pool, task.Worker, task.Finish, error) {
	return nil, nil, nil, bundle.Upload(ctx.Project)
}
