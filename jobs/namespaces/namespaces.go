package namespaces

import (
	"okapi/lib/task"
)

// Name task name for trigger
var Name task.Name = "namespaces"

// Task getting all namespaces
func Task(ctx *task.Context) (task.Pool, task.Worker, task.Finish, error) {
	return Pool(ctx), Worker(ctx), nil, nil
}
