package index

import (
	"okapi/lib/task"
)

// Name task name for trigger
var Name task.Name = "index"

// Task to index all the records in database
func Task(ctx *task.Context) (task.Pool, task.Worker, task.Finish, error) {
	return Pool(ctx), Worker, nil, nil
}
