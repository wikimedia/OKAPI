package pull

import (
	"okapi/lib/task"
)

// Name task name for trigger
var Name task.Name = "pull"

// Task for getting html for pages
func Task(ctx *task.Context) (task.Pool, task.Worker, task.Finish, error) {
	return Pool(ctx), Worker, nil, nil
}
